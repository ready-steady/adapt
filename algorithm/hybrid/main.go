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

	// Index returns the nodal indices of a set of level indices.
	Index([]uint64) []uint64

	// ChildrenToward returns the child indices of a set of indices with respect
	// to a particular dimension.
	ChildrenToward([]uint64, uint) []uint64
}

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	grid   Grid
	basis  Basis
	config Config
}

type state struct {
	Lindices []uint64
	Indices  []uint64
	Counts   []uint

	Nodes        []float64
	Volumes      []float64
	Observations []float64
	Predictions  []float64
	Surpluses    []float64
	Scores       []float64
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
	strategy := newStrategy(ni, no, self.grid, config)

	state := strategy.Next(nil, nil)
	progress.Push(state.Indices, ni)
	for !target.Done(progress) && !strategy.Done() {
		state.Volumes = internal.Measure(self.basis, state.Indices, ni)
		state.Nodes = self.grid.Compute(state.Indices)
		state.Observations = internal.Invoke(target.Compute, state.Nodes, ni, no, nw)
		state.Predictions = internal.Approximate(self.basis, surrogate.Indices,
			surrogate.Surpluses, state.Nodes, ni, no, nw)
		state.Surpluses = internal.Subtract(state.Observations, state.Predictions)
		state.Scores = score(target, state.Indices, state.Volumes, state.Observations,
			state.Surpluses, state.Counts, ni, no)

		surrogate.Push(state.Indices, state.Surpluses, state.Volumes)
		state = strategy.Next(state, surrogate)
		progress.Push(state.Indices, ni)
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *external.Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
}
