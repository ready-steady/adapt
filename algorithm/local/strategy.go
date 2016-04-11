package local

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy is a basic strategy.
type Strategy struct {
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
	localError float64, grid Grid) *Strategy {

	return &Strategy{
		ni: inputs,
		no: outputs,

		lmin: minLevel,
		lmax: maxLevel,
		εl:   localError,

		grid: grid,
	}
}

func (self *Strategy) First() *external.State {
	self.unique = internal.NewUnique(self.ni)
	return &external.State{
		Indices: make([]uint64, 1*self.ni),
	}
}

func (self *Strategy) Check(state *external.State, _ *external.Surrogate) bool {
	return state != nil && len(state.Indices) > 0
}

func (self *Strategy) Score(element *external.Element) float64 {
	return internal.MaxAbsolute(element.Surplus)
}

func (self *Strategy) Next(state *external.State, _ *external.Surrogate) *external.State {
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
