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

	grid  Grid
	basis Basis
}

// New creates an interpolator.
func New(inputs, outputs uint, grid Grid, basis Basis) *Interpolator {
	return &Interpolator{
		ni: inputs,
		no: outputs,

		grid:  grid,
		basis: basis,
	}
}

// Compute constructs an interpolant for a function.
func (self *Interpolator) Compute(target external.Target,
	strategy external.Strategy) *external.Surrogate {

	ni, no := self.ni, self.no
	surrogate := external.NewSurrogate(ni, no)
	for s := strategy.First(); !strategy.Done(s, surrogate); s = strategy.Next(s, surrogate) {
		s.Volumes = internal.Measure(self.basis, s.Indices, ni)
		s.Nodes = self.grid.Compute(s.Indices)
		s.Values = external.Invoke(target, s.Nodes, ni, no)
		s.Estimates = internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, s.Nodes, ni, no)
		s.Surpluses = internal.Subtract(s.Values, s.Estimates)
		s.Scores = score(strategy, s, ni, no)
		surrogate.Push(s.Indices, s.Surpluses, s.Volumes)
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
				Value:   state.Values[j*no : (j+1)*no],
				Surplus: state.Surpluses[j*no : (j+1)*no],
			}))
		}
	}
	return scores
}
