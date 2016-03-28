// Package global provides an algorithm for hierarchical interpolation with
// global adaptation.
package global

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Basis is an interpolation basis.
type Basis interface {
	// Compute evaluates the value of a basis function.
	Compute([]uint64, []float64) float64

	// Integrate computes the integral of a basis function.
	Integrate([]uint64) float64
}

// Grid is an interpolation grid.
type Grid interface {
	// Compute returns the nodes corresponding to a set of indices.
	Compute([]uint64) []float64

	// Index returns the nodal indices of a set of level indices.
	Index([]uint64) []uint64
}

// Target is a function to be interpolated.
type Target func([]float64, []float64)

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	ni uint
	no uint

	grid   Grid
	basis  Basis
	config *Config
}

// Element contains information about an interpolation element.
type Element struct {
	Lindex  []uint64 // Level index
	Indices []uint64 // Nodal indices

	Volumes      []float64 // Basis-function volumes
	Observations []float64 // Target-function values
	Surpluses    []float64 // Hierarchical surpluses
}

// State contains information about an interpolation iteration.
type State struct {
	Lindices []uint64 // Level indices
	Indices  []uint64 // Nodal indices
	Counts   []uint   // Number of nodal indices for each level index

	Nodes        []float64 // Grid nodes
	Volumes      []float64 // Basis-function volumes
	Observations []float64 // Target-function values
	Predictions  []float64 // Approximated values
	Surpluses    []float64 // Hierarchical surpluses
	Scores       []float64 // Level-index scores
}

// New creates an interpolator.
func New(inputs, outputs uint, grid Grid, basis Basis, config *Config) *Interpolator {
	return &Interpolator{
		ni: inputs,
		no: outputs,

		grid:   grid,
		basis:  basis,
		config: config,
	}
}

// Compute constructs an interpolant for a function.
func (self *Interpolator) Compute(target Target) *external.Surrogate {
	ni, no := self.ni, self.no

	progress := external.NewProgress()
	surrogate := external.NewSurrogate(ni, no)
	strategy := NewStrategy(ni, no, self.grid, self.config)

	state := strategy.Next(nil, nil)
	progress.Push(state.Indices, ni)
	for !strategy.Done() {
		state.Volumes = internal.Measure(self.basis, state.Indices, ni)
		state.Nodes = self.grid.Compute(state.Indices)
		state.Observations = internal.Invoke(target, state.Nodes, ni, no, internal.Workers)
		state.Predictions = internal.Approximate(self.basis, surrogate.Indices,
			surrogate.Surpluses, state.Nodes, ni, no, internal.Workers)
		state.Surpluses = internal.Subtract(state.Observations, state.Predictions)
		state.Scores = score(strategy, state, ni, no)

		surrogate.Push(state.Indices, state.Surpluses, state.Volumes)
		state = strategy.Next(state, surrogate)
		progress.Push(state.Indices, ni)
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *external.Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, internal.Workers)
}
