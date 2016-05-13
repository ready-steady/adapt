package global

import (
	"math"

	"github.com/ready-steady/adapt/algorithm"
	"github.com/ready-steady/adapt/algorithm/internal"
	"github.com/ready-steady/adapt/grid"
)

var (
	infinity = math.Inf(1.0)
)

// Strategy is a basic strategy.
type Strategy struct {
	ni uint
	no uint

	guide Guide

	lmin uint
	lmax uint

	priority []float64
	accuracy []float64

	active    *internal.Active
	threshold *internal.Threshold
}

// Guide is a grid-refinement tool of a basic strategy.
type Guide interface {
	grid.Indexer
}

// NewStrategy creates a basic strategy.
func NewStrategy(inputs, outputs uint, guide Guide, minLevel, maxLevel uint,
	absoluteError, relativeError float64) *Strategy {

	return &Strategy{
		ni: inputs,
		no: outputs,

		guide: guide,

		lmin: minLevel,
		lmax: maxLevel,

		active:    internal.NewActive(inputs),
		threshold: internal.NewThreshold(outputs, absoluteError, relativeError),
	}
}

func (self *Strategy) First(surrogate *algorithm.Surrogate) *algorithm.State {
	return self.initiate(self.active.First(), surrogate)
}

func (self *Strategy) Next(state *algorithm.State,
	surrogate *algorithm.Surrogate) *algorithm.State {

	for {
		self.consume(state)
		if self.check() {
			return nil
		}
		k := self.choose()
		if k == internal.None {
			return nil
		}
		state = self.initiate(self.active.Next(k), surrogate)
		if len(state.Indices) > 0 {
			return state
		}
	}
}

func (self *Strategy) Score(element *algorithm.Element) float64 {
	return internal.SumAbsolute(element.Surplus)
}

func (self *Strategy) check() bool {
	for i := range self.active.Positions {
		if !self.threshold.Check(self.accuracy[i*self.no : (i+1)*self.no]) {
			return false
		}
	}
	return true
}

func (self *Strategy) choose() uint {
	return internal.Choose(self.priority, self.active.Positions)
}

func (self *Strategy) consume(state *algorithm.State) {
	ni, no := self.ni, self.no
	np := uint(len(self.priority))
	na := uint(len(self.accuracy))
	nn := uint(len(state.Counts))

	self.priority = append(self.priority, make([]float64, nn)...)
	priority := self.priority[np:]

	self.accuracy = append(self.accuracy, make([]float64, nn*no)...)
	accuracy := self.accuracy[na:]

	levels := internal.Levelize(state.Lndices, ni)

	for i, o := uint(0), uint(0); i < nn; i++ {
		count := state.Counts[i]
		if levels[i] < uint64(self.lmin) {
			priority[i] = infinity
			internal.Set(accuracy[i*no:(i+1)*no], infinity)
		} else if levels[i] < uint64(self.lmax) {
			priority[i] = internal.Average(state.Scores[o:(o + count)])
			self.threshold.Compress(accuracy[i*no:(i+1)*no],
				state.Surpluses[o*no:(o+count)*no])
		}
		o += count
	}

	self.threshold.Update(state.Values)
}

func (self *Strategy) initiate(lndices []uint64, _ *algorithm.Surrogate) (state *algorithm.State) {
	state = &algorithm.State{Lndices: lndices}
	state.Indices, state.Counts = internal.Index(self.guide, lndices, self.ni)
	return
}
