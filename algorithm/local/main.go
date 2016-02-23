// Package local provides an algorithm for hierarchical interpolation with local
// adaptation.
package local

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

const (
	levelMask = 0x3F
	levelSize = 6
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

	// Children returns the child indices corresponding to a set of indices.
	Children([]uint64) []uint64
}

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	grid   Grid
	basis  Basis
	config Config
}

// New creates an interpolator.
func New(grid Grid, basis Basis, config *Config) *Interpolator {
	if config.Workers == 0 {
		panic("the number of workers should be positive")
	}
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
	unique := internal.NewUnique(ni)

	indices := make([]uint64, 1*ni)
	progress := &Progress{Active: 1}
	for target.Continue(progress) && progress.Active > 0 {
		nodes := self.grid.Compute(indices)
		values := internal.Invoke(target.Compute, nodes, ni, no, nw)
		surpluses := internal.Subtract(values, internal.Approximate(self.basis,
			surrogate.Indices, surrogate.Surpluses, nodes, ni, no, nw))

		surrogate.Push(self.basis, indices, surpluses)

		scores := assess(self.basis, target, indices, values, surpluses, ni, no)
		indices = filter(indices, scores, config.MinLevel, config.MaxLevel, ni)
		indices = unique.Distil(self.grid.Children(indices))

		progress.Passive += progress.Active
		progress.Active = uint(len(indices)) / ni
		progress.Level++
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *external.Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
}
