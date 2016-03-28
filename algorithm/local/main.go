// Package local provides an algorithm for hierarchical interpolation with local
// adaptation.
package local

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

	// Children returns the child indices corresponding to a set of indices.
	Children([]uint64) []uint64
}

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
	unique := internal.NewUnique(ni)

	state := &State{}
	state.Indices = make([]uint64, 1*ni)
	progress.Push(state.Indices, ni)
	for !target.Done(progress) && progress.More > 0 {
		state.Volumes = internal.Measure(self.basis, state.Indices, ni)
		state.Nodes = self.grid.Compute(state.Indices)
		state.Observations = internal.Invoke(target.Compute, state.Nodes, ni, no, internal.Workers)
		state.Predictions = internal.Approximate(self.basis, surrogate.Indices,
			surrogate.Surpluses, state.Nodes, ni, no, internal.Workers)
		state.Surpluses = internal.Subtract(state.Observations, state.Predictions)
		state.Scores = score(target, state, ni, no)

		surrogate.Push(state.Indices, state.Surpluses, state.Volumes)
		state.Indices = unique.Distil(self.grid.Children(filter(state.Indices, state.Scores,
			self.config.MinLevel, self.config.MaxLevel, ni)))
		progress.Push(state.Indices, ni)
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *external.Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, internal.Workers)
}
