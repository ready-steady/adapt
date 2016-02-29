// Package hybrid provides an algorithm for hierarchical interpolation with
// hybrid adaptation.
package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
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

	progress := external.NewProgress()
	surrogate := external.NewSurrogate(ni, no)
	active := external.NewActive(ni, config.MaxLevel, config.MaxIndices)
	strategy := newStrategy(ni, no, config.AbsoluteError, config.RelativeError)

	k := ^uint(0)
	indices, counts := index(self.grid, active.Begin(), ni)
	progress.Push(indices, ni)
	for target.Continue(progress) && strategy.Continue(active) {
		nodes := self.grid.Compute(indices)
		values := internal.Invoke(target.Compute, nodes, ni, no, nw)
		surpluses := internal.Subtract(values, internal.Approximate(self.basis,
			surrogate.Indices, surrogate.Surpluses, nodes, ni, no, nw))

		surrogate.Push(self.basis, indices, surpluses)
		score(self.basis, strategy, target, counts, indices, values, surpluses, ni, no)

		active.Remove(k)
		k = strategy.Select(active)

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
