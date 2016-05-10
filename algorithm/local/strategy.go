package local

import (
	"github.com/ready-steady/adapt/algorithm"
	"github.com/ready-steady/adapt/algorithm/internal"
	"github.com/ready-steady/adapt/grid"
)

// Strategy is a basic strategy.
type Strategy struct {
	ni uint
	no uint

	guide Guide

	lmin uint
	lmax uint
	εl   float64

	unique *internal.Unique
}

// Guide is a grid-refinement tool of a basic strategy.
type Guide interface {
	grid.Refiner
}

// NewStrategy creates a basic strategy.
func NewStrategy(inputs, outputs uint, guide Guide, minLevel, maxLevel uint,
	localError float64) *Strategy {

	return &Strategy{
		ni: inputs,
		no: outputs,

		guide: guide,

		lmin: minLevel,
		lmax: maxLevel,
		εl:   localError,

		unique: internal.NewUnique(inputs),
	}
}

func (self *Strategy) First() *algorithm.State {
	return &algorithm.State{
		Indices: make([]uint64, 1*self.ni),
	}
}

func (self *Strategy) Done(state *algorithm.State, _ *algorithm.Surrogate) bool {
	return state == nil || len(state.Indices) == 0
}

func (self *Strategy) Score(element *algorithm.Element) float64 {
	return internal.MaxAbsolute(element.Surplus)
}

func (self *Strategy) Next(state *algorithm.State, _ *algorithm.Surrogate) *algorithm.State {
	return &algorithm.State{
		Indices: self.unique.Distil(self.guide.Refine(filter(state.Indices,
			state.Scores, self.lmin, self.lmax, self.εl, self.ni))),
	}
}

func filter(indices []uint64, scores []float64, lmin, lmax uint, εl float64, ni uint) []uint64 {
	nn := uint(len(scores))
	levels := internal.Levelize(indices, ni)
	na, ne := uint(0), nn
	for i, j := uint(0), uint(0); i < nn; i++ {
		if levels[i] >= uint64(lmin) && (scores[i] <= εl || levels[i] >= uint64(lmax)) {
			j++
			continue
		}
		if j > na {
			copy(indices[na*ni:], indices[j*ni:ne*ni])
			ne -= j - na
			j = na
		}
		na++
		j++
	}
	return indices[:na*ni]
}
