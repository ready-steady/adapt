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
	ε    float64

	unique *internal.Unique
}

// Guide is a grid-refinement tool of a basic strategy.
type Guide interface {
	grid.Indexer
	grid.Refiner
}

// NewStrategy creates a basic strategy.
func NewStrategy(inputs, outputs uint, guide Guide, minLevel, maxLevel uint,
	scoreError float64) *Strategy {

	return &Strategy{
		ni: inputs,
		no: outputs,

		guide: guide,

		lmin: minLevel,
		lmax: maxLevel,
		ε:    scoreError,

		unique: internal.NewUnique(inputs),
	}
}

func (self *Strategy) First(_ *algorithm.Surrogate) *algorithm.State {
	lndex := make([]uint64, self.ni)
	return &algorithm.State{Indices: self.guide.Index(lndex)}
}

func (self *Strategy) Next(state *algorithm.State, _ *algorithm.Surrogate) *algorithm.State {
	indices := self.unique.Distil(self.guide.Refine(filter(state.Indices,
		state.Scores, self.lmin, self.lmax, self.ε, self.ni)))
	if len(indices) == 0 {
		return nil
	}
	return &algorithm.State{
		Indices: indices,
	}
}

func (self *Strategy) Score(element *algorithm.Element) float64 {
	return internal.MaxAbsolute(element.Surplus)
}

func filter(indices []uint64, scores []float64, lmin, lmax uint, ε float64, ni uint) []uint64 {
	nn := uint(len(scores))
	levels := internal.Levelize(indices, ni)
	na, ne := uint(0), nn
	for i, j := uint(0), uint(0); i < nn; i++ {
		if levels[i] >= uint64(lmin) && (scores[i] <= ε || levels[i] >= uint64(lmax)) {
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
