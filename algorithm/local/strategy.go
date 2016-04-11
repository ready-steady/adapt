package local

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy controls the interpolation process.
type Strategy interface {
	// First returns the initial state of the first iteration.
	First() *external.State

	// Check returns true if the interpolation process should continue.
	Check(*external.State, *external.Surrogate) bool

	// Score assigns a score to an interpolation element.
	Score(*external.Element) float64

	// Next consumes the result of the current iteration and returns the initial
	// state of the next one.
	Next(*external.State, *external.Surrogate) *external.State
}

// BasicStrategy is a basic strategy satisfying the Strategy interface.
type BasicStrategy struct {
	ni uint
	no uint

	lmin uint
	lmax uint
	εl   float64

	grid   Grid
	unique *internal.Unique
}

// NewStrategy creates a basic strategy.
func NewStrategy(inputs, outputs, minLevel, maxLevel uint,
	localError float64, grid Grid) *BasicStrategy {

	return &BasicStrategy{
		ni: inputs,
		no: outputs,

		lmin: minLevel,
		lmax: maxLevel,
		εl:   localError,

		grid: grid,
	}
}

func (self *BasicStrategy) First() *external.State {
	self.unique = internal.NewUnique(self.ni)
	return &external.State{
		Indices: make([]uint64, 1*self.ni),
	}
}

func (self *BasicStrategy) Check(state *external.State, _ *external.Surrogate) bool {
	return state != nil && len(state.Indices) > 0
}

func (self *BasicStrategy) Score(element *external.Element) float64 {
	return internal.MaxAbsolute(element.Surplus)
}

func (self *BasicStrategy) Next(state *external.State, _ *external.Surrogate) *external.State {
	return &external.State{
		Indices: self.unique.Distil(self.grid.Refine(filter(state.Indices,
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
