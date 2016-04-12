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

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	ni uint
	no uint

	grid     Grid
	basis    Basis
	strategy external.Strategy
}

// New creates an interpolator.
func New(inputs, outputs uint, grid Grid, basis Basis, strategy external.Strategy) *Interpolator {
	return &Interpolator{
		ni: inputs,
		no: outputs,

		grid:     grid,
		basis:    basis,
		strategy: strategy,
	}
}

// Compute constructs an interpolant for a function.
func (self *Interpolator) Compute(target external.Target) *external.Surrogate {
	ni, no := self.ni, self.no
	surrogate := external.NewSurrogate(ni, no)
	state := self.strategy.First()
	for self.strategy.Check(state, surrogate) {
		state.Volumes = internal.Measure(self.basis, state.Indices, ni)
		state.Nodes = self.grid.Compute(state.Indices)
		state.Observations = external.Invoke(target, state.Nodes, ni, no)
		state.Predictions = internal.Approximate(self.basis, surrogate.Indices,
			surrogate.Surpluses, state.Nodes, ni, no)
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
		surrogate.Inputs, surrogate.Outputs)
}

func score(strategy external.Strategy, state *external.State, ni, no uint) []float64 {
	nn := uint(len(state.Counts))
	scores := make([]float64, 0, nn)
	for i, j := uint(0), uint(0); i < nn; i++ {
		for m := j + state.Counts[i]; j < m; j++ {
			scores = append(scores, strategy.Score(&external.Element{
				Index:   state.Indices[j*ni : (j+1)*ni],
				Volume:  state.Volumes[j],
				Value:   state.Observations[j*no : (j+1)*no],
				Surplus: state.Surpluses[j*no : (j+1)*no],
			}))
		}
	}
	return scores
}
