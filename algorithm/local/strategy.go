package local

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy controls the interpolation process.
type Strategy interface {
	// First returns the initial state of the first iteration.
	First() *State

	// Continue returns true if the interpolation process should continue.
	Continue(*State, *external.Surrogate) bool

	// Score assigns a score to an interpolation element.
	Score(*Element) float64

	// Next consumes the result of the current iteration and returns the initial
	// state of the next one.
	Next(*State, *external.Surrogate) *State
}

// BasicStrategy is a basic target satisfying the Target interface.
type BasicStrategy struct {
	ni uint
	no uint

	lmin uint
	lmax uint
	εa   float64

	grid   Grid
	unique *internal.Unique
}

// NewStrategy creates a basic strategy.
func NewStrategy(inputs, outputs, minLevel, maxLevel uint,
	absoluteError float64, grid Grid) *BasicStrategy {

	return &BasicStrategy{
		ni: inputs,
		no: outputs,

		lmin: minLevel,
		lmax: maxLevel,
		εa:   absoluteError,

		grid: grid,
	}
}

func (self *BasicStrategy) First() *State {
	self.unique = internal.NewUnique(self.ni)
	return &State{
		Indices: make([]uint64, 1*self.ni),
	}
}

func (self *BasicStrategy) Continue(state *State, _ *external.Surrogate) bool {
	return state != nil && len(state.Indices) > 0
}

func (self *BasicStrategy) Score(element *Element) float64 {
	for _, εa := range element.Surplus {
		if math.Abs(εa) > self.εa {
			return 1.0
		}
	}
	return 0.0
}

func (self *BasicStrategy) Next(state *State, _ *external.Surrogate) *State {
	return &State{
		Indices: self.unique.Distil(self.grid.Refine(filter(
			state.Indices, state.Scores, self.lmin, self.lmax, self.ni))),
	}
}

func filter(indices []uint64, scores []float64, lmin, lmax, ni uint) []uint64 {
	nn := uint(len(scores))
	levels := internal.Levelize(indices, ni)
	na, ne := uint(0), nn
	for i, j := uint(0), uint(0); i < nn; i++ {
		if levels[i] >= uint64(lmin) && (scores[i] <= 0.0 || levels[i] >= uint64(lmax)) {
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
