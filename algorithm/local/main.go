// Package local provides an algorithm for hierarchical interpolation with local
// adaptation.
package local

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Basis is an interpolation basis.
type Basis interface {
	internal.BasisComputer
	internal.BasisIntegrator
}

// Grid is an interpolation grid.
type Grid interface {
	internal.GridComputer
	internal.GridRefiner
}

// Target is a function to be interpolated.
type Target func([]float64, []float64)

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	ni uint
	no uint

	grid     Grid
	basis    Basis
	strategy Strategy
}

// Element contains information about an interpolation element.
type Element struct {
	Index []uint64 // Nodal index

	Volume      float64   // Basis-function volume
	Observation []float64 // Target-function value
	Surplus     []float64 // Hierarchical surplus
}

// State contains information about an interpolation iteration.
type State struct {
	Indices []uint64 // Nodal indices

	Nodes        []float64 // Grid nodes
	Volumes      []float64 // Basis-function volumes
	Observations []float64 // Target-function values
	Predictions  []float64 // Approximated values
	Surpluses    []float64 // Hierarchical surpluses
	Scores       []float64 // Nodal-index scores
}

// New creates an interpolator.
func New(inputs, outputs uint, grid Grid, basis Basis, strategy Strategy) *Interpolator {
	return &Interpolator{
		ni: inputs,
		no: outputs,

		grid:     grid,
		basis:    basis,
		strategy: strategy,
	}
}

// Compute constructs an interpolant for a function.
func (self *Interpolator) Compute(target Target) *external.Surrogate {
	ni, no := self.ni, self.no

	progress := external.NewProgress()
	surrogate := external.NewSurrogate(ni, no)

	state := self.strategy.First()
	progress.Push(state.Indices, ni)
	for !self.strategy.Check(progress) {
		state.Volumes = internal.Measure(self.basis, state.Indices, ni)
		state.Nodes = self.grid.Compute(state.Indices)
		state.Observations = internal.Invoke(target, state.Nodes, ni, no, internal.Workers)
		state.Predictions = internal.Approximate(self.basis, surrogate.Indices,
			surrogate.Surpluses, state.Nodes, ni, no, internal.Workers)
		state.Surpluses = internal.Subtract(state.Observations, state.Predictions)
		state.Scores = score(self.strategy, state, ni, no)

		surrogate.Push(state.Indices, state.Surpluses, state.Volumes)
		state = self.strategy.Next(state, surrogate)
		progress.Push(state.Indices, ni)
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *external.Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, internal.Workers)
}

func score(strategy Strategy, state *State, ni, no uint) []float64 {
	nn := uint(len(state.Indices)) / ni
	scores := make([]float64, nn)
	for i := uint(0); i < nn; i++ {
		scores[i] = strategy.Score(&Element{
			Index:       state.Indices[i*ni : (i+1)*ni],
			Volume:      state.Volumes[i],
			Observation: state.Observations[i*no : (i+1)*no],
			Surplus:     state.Surpluses[i*no : (i+1)*no],
		})
	}
	return scores
}
