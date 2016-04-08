// Package global provides an algorithm for hierarchical interpolation with
// global adaptation.
package global

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
	internal.GridIndexer
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
	surrogate := external.NewSurrogate(ni, no)
	state := self.strategy.First()
	for self.strategy.Continue(state, surrogate) {
		state.Volumes = internal.Measure(self.basis, state.Indices, ni)
		state.Nodes = self.grid.Compute(state.Indices)
		state.Observations = internal.Invoke(target, state.Nodes, ni, no, internal.Workers)
		state.Predictions = internal.Approximate(self.basis, surrogate.Indices,
			surrogate.Surpluses, state.Nodes, ni, no, internal.Workers)
		state.Surpluses = internal.Subtract(state.Observations, state.Predictions)
		state.Scores = score(self.strategy, state, ni, no)
		surrogate.Push(state.Indices, state.Surpluses, state.Volumes)
		state = self.strategy.Next(state, surrogate)
	}
	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *external.Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, internal.Workers)
}

func score(strategy Strategy, state *State, ni, no uint) []float64 {
	nn := uint(len(state.Counts))
	scores := make([]float64, nn)
	for i, o := uint(0), uint(0); i < nn; i++ {
		count := state.Counts[i]
		element := Element{
			Lindex:       state.Lindices[i*ni : (i+1)*ni],
			Indices:      state.Indices[o*ni : (o+count)*ni],
			Volumes:      state.Volumes[o:(o + count)],
			Observations: state.Observations[o*no : (o+count)*no],
			Surpluses:    state.Surpluses[o*no : (o+count)*no],
		}
		scores[i] = strategy.Score(&element)
		o += count
	}
	return scores
}
