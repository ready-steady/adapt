// Package hybrid provides an algorithm for hierarchical interpolation with
// hybrid adaptation.
package hybrid

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

var (
	infinity = math.Inf(1)
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

	surrogate := external.NewSurrogate(ni, no)
	active := external.NewActive(ni, config.MaxLevel, config.MaxIndices)

	k, indices, counts := ^uint(0), make([]uint64, 1*ni), []uint{1}
	progress := &Progress{More: 1}
	for target.Continue(active, progress) {
		nodes := self.grid.Compute(indices)
		values := internal.Invoke(target.Compute, nodes, ni, no, nw)
		surpluses := internal.Subtract(values, internal.Approximate(self.basis,
			surrogate.Indices, surrogate.Surpluses, nodes, ni, no, nw))

		surrogate.Push(self.basis, indices, surpluses)
		score(self.basis, target, counts, indices, values, surpluses, ni, no)

		active.Forget(k)
		k = target.Select(active)
		_ = active.Advance(k)

		progress.Done += progress.More
		progress.More = uint(len(indices)) / ni
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *external.Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
}

func score(basis Basis, target Target, counts []uint, indices []uint64,
	values, surpluses []float64, ni, no uint) {

	for _, count := range counts {
		oi, oo := count*ni, count*no
		target.Score(&Location{
			Indices:   indices[:oi],
			Values:    values[:oo],
			Surpluses: surpluses[:oo],
			Volumes:   internal.Measure(basis, indices[:oo], ni),
		})
		indices, values, surpluses = indices[oi:], values[oo:], surpluses[oo:]
	}
}
