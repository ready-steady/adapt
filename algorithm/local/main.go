// Package local provides an algorithm for hierarchical interpolation with local
// adaptation.
package local

import (
	"github.com/ready-steady/adapt/algorithm"
	"github.com/ready-steady/adapt/algorithm/internal"
	"github.com/ready-steady/adapt/basis"
	"github.com/ready-steady/adapt/grid"
)

// Algorithm is the interpolation algorithm.
type Algorithm struct {
	ni uint
	no uint

	grid  Grid
	basis Basis
}

// Basis is an interpolation basis.
type Basis interface {
	basis.Computer
	basis.Integrator
}

// Grid is an interpolation grid.
type Grid interface {
	grid.Computer
}

// New creates an interpolator.
func New(inputs, outputs uint, grid Grid, basis Basis) *Algorithm {
	return &Algorithm{
		ni: inputs,
		no: outputs,

		grid:  grid,
		basis: basis,
	}
}

// Compute constructs an interpolant for a function.
func (self *Algorithm) Compute(target algorithm.Target,
	strategy algorithm.Strategy) *algorithm.Surrogate {

	ni, no := self.ni, self.no
	surrogate := algorithm.NewSurrogate(ni, no)
	for s := strategy.First(surrogate); s != nil; s = strategy.Next(s, surrogate) {
		s.Volumes = internal.Measure(self.basis, s.Indices, ni)
		s.Nodes = self.grid.Compute(s.Indices)
		s.Values = algorithm.Invoke(target, s.Nodes, ni, no)
		s.Estimates = internal.Estimate(self.basis, surrogate.Indices,
			surrogate.Surpluses, s.Nodes, ni, no)
		s.Surpluses = internal.Subtract(s.Values, s.Estimates)
		s.Scores = score(strategy, s, ni, no)
		surrogate.Push(s.Indices, s.Surpluses, s.Volumes)
	}
	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Algorithm) Evaluate(surrogate *algorithm.Surrogate, points []float64) []float64 {
	return internal.Estimate(self.basis, surrogate.Indices, surrogate.Surpluses,
		points, surrogate.Inputs, surrogate.Outputs)
}

func score(strategy algorithm.Strategy, state *algorithm.State, ni, no uint) []float64 {
	nn := uint(len(state.Indices)) / ni
	scores := make([]float64, nn)
	for i := uint(0); i < nn; i++ {
		scores[i] = strategy.Score(&algorithm.Element{
			Index:   state.Indices[i*ni : (i+1)*ni],
			Node:    state.Nodes[i*ni : (i+1)*ni],
			Volume:  state.Volumes[i],
			Value:   state.Values[i*no : (i+1)*no],
			Surplus: state.Surpluses[i*no : (i+1)*no],
		})
	}
	return scores
}
