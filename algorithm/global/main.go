// Package global provides an algorithm for hierarchical interpolation with
// global adaptation.
package global

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
}

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	grid   Grid
	basis  Basis
	config Config
}

type state struct {
	lindices []uint64
	indices  []uint64
	counts   []uint

	nodes        []float64
	volumes      []float64
	observations []float64
	predictions  []float64
	surpluses    []float64
	scores       []float64
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
	strategy := newStrategy(ni, no, self.grid, surrogate, config)

	state := &state{}
	strategy.Next(state)
	progress.Push(state.indices, ni)
	for !target.Done(progress) && !strategy.Done() {
		state.volumes = internal.Measure(self.basis, state.indices, ni)
		state.nodes = self.grid.Compute(state.indices)
		state.observations = internal.Invoke(target.Compute, state.nodes, ni, no, nw)
		state.predictions = internal.Approximate(self.basis, surrogate.Indices,
			surrogate.Surpluses, state.nodes, ni, no, nw)
		state.surpluses = internal.Subtract(state.observations, state.predictions)
		state.scores = score(target, state.indices, state.volumes, state.observations,
			state.surpluses, state.counts, ni, no)

		strategy.Push(state)
		strategy.Next(state)
		progress.Push(state.indices, ni)
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *external.Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
}
