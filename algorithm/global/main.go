// Package global provides an algorithm for hierarchical interpolation with
// global adaptation.
package global

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

var (
	infinity = math.Inf(1.0)
)

// Basis is a functional basis.
type Basis interface {
	// Compute evaluates the value of a basis function.
	Compute([]uint64, []float64) float64

	// Integrate computes the integral of a basis function.
	Integrate([]uint64) float64
}

// Grid is a sparse grid.
type Grid interface {
	// Compute returns the nodes corresponding to a set of indices.
	Compute([]uint64) []float64

	// Index returns the indices of a set of levels.
	Index([]uint64) []uint64
}

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	grid   Grid
	basis  Basis
	config Config
}

// New creates an interpolator.
func New(grid Grid, basis Basis, config *Config) *Interpolator {
	return &Interpolator{
		grid:   grid,
		basis:  basis,
		config: *config,
	}
}

// Compute constructs an interpolant for a function.
func (self *Interpolator) Compute(target Target) *external.Surrogate {
	config := &self.config

	ni, no := target.Dimensions()
	nw := config.Workers

	active := external.NewActive(ni, config.MaxLevel, config.MaxIndices)
	progress := external.NewProgress()
	surrogate := external.NewSurrogate(ni, no)

	k := ^uint(0)
	indices, counts := index(self.grid, active.Begin(), ni)
	progress.Push(indices, ni)
	for target.Continue(active, progress) {
		nodes := self.grid.Compute(indices)
		values := internal.Invoke(target.Compute, nodes, ni, no, nw)
		surpluses := internal.Subtract(values, internal.Approximate(self.basis,
			surrogate.Indices, surrogate.Surpluses, nodes, ni, no, nw))

		surrogate.Push(self.basis, indices, surpluses)
		score(self.basis, target, counts, indices, values, surpluses, ni, no)

		active.Remove(k)
		k = target.Select(active)

		indices, counts = index(self.grid, active.Forward(k), ni)
		progress.Push(indices, ni)
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *external.Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
}

func index(grid Grid, lindices []uint64, ni uint) ([]uint64, []uint) {
	nn := uint(len(lindices)) / ni
	indices, counts := []uint64(nil), make([]uint, nn)
	for i := uint(0); i < nn; i++ {
		newIndices := grid.Index(lindices[:ni])
		indices = append(indices, newIndices...)
		counts[i] = uint(len(newIndices)) / ni
		lindices = lindices[ni:]
	}
	return indices, counts
}

func score(basis Basis, target Target, counts []uint, indices []uint64,
	values, surpluses []float64, ni, no uint) {

	for _, count := range counts {
		oi, oo := count*ni, count*no
		target.Score(&Location{
			Indices:   indices[:oi],
			Volumes:   internal.Measure(basis, indices[:oo], ni),
			Values:    values[:oo],
			Surpluses: surpluses[:oo],
		})
		indices, values, surpluses = indices[oi:], values[oo:], surpluses[oo:]
	}
}
