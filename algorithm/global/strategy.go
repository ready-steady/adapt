package global

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

var (
	infinity = math.Inf(1.0)
)

// Strategy is a basic strategy.
type Strategy struct {
	*internal.Active

	ni uint
	no uint

	guide Guide

	lmin uint
	lmax uint

	k uint

	global []float64
	local  []float64

	threshold *internal.Threshold
}

// Guide is a grid-refinement tool of a basic strategy.
type Guide interface {
	internal.GridIndexer
}

// NewStrategy creates a basic strategy.
func NewStrategy(inputs, outputs uint, guide Guide, minLevel, maxLevel uint,
	absoluteError, relativeError float64) *Strategy {

	return &Strategy{
		Active: internal.NewActive(inputs),

		ni: inputs,
		no: outputs,

		guide: guide,

		lmin: minLevel,
		lmax: maxLevel,

		k: ^uint(0),

		threshold: internal.NewThreshold(outputs, absoluteError, relativeError),
	}
}

func (self *Strategy) First() *external.State {
	state := &external.State{}
	state.Lindices = self.Active.First()
	state.Indices, state.Counts = internal.Index(self.guide, state.Lindices, self.ni)
	return state
}

func (self *Strategy) Done(_ *external.State, _ *external.Surrogate) bool {
	if self.k == ^uint(0) {
		return false
	}
	no := self.no
	nl := uint(len(self.local)) / no
	δ := self.threshold.Values
	for i := range self.Positions {
		if i < nl {
			for j := uint(0); j < no; j++ {
				if self.local[i*no+j] > δ[j] {
					return false
				}
			}
		}
	}
	return true
}

func (self *Strategy) Score(element *external.Element) float64 {
	return internal.SumAbsolute(element.Surplus)
}

func (self *Strategy) Next(state *external.State, _ *external.Surrogate) *external.State {
	for {
		self.consume(state)
		self.Active.Drop(self.k)
		if len(self.Positions) == 0 {
			return nil
		}
		self.k = internal.LocateMax(self.global, self.Positions)
		if self.global[self.k] <= 0.0 {
			return nil
		}
		state = &external.State{}
		state.Lindices = self.Active.Next(self.k)
		state.Indices, state.Counts = internal.Index(self.guide, state.Lindices, self.ni)
		if len(state.Indices) > 0 {
			return state
		}
	}
}

func (self *Strategy) consume(state *external.State) {
	no, ng, nl := self.no, uint(len(self.global)), uint(len(self.local))
	nn := uint(len(state.Counts))

	levels := internal.Levelize(state.Lindices, self.ni)

	self.global = append(self.global, make([]float64, nn)...)
	global := self.global[ng:]

	self.local = append(self.local, make([]float64, nn*no)...)
	local := self.local[nl:]

	for i, o := uint(0), uint(0); i < nn; i++ {
		count := state.Counts[i]
		if levels[i] < uint64(self.lmin) {
			global[i] = infinity
			for j := uint(0); j < no; j++ {
				local[i*no+j] = infinity
			}
		} else if levels[i] >= uint64(self.lmax) {
			global[i] = 0.0
			for j := uint(0); j < no; j++ {
				local[i*no+j] = 0.0
			}
		} else {
			global[i] = internal.Average(state.Scores[o:(o + count)])
			for j, m := uint(0), count*no; j < m; j++ {
				k := i*no + j%no
				local[k] = math.Max(local[k], math.Abs(state.Surpluses[o+j]))
			}
		}
		o += count
	}

	self.threshold.Update(state.Values)
}
