package global

import (
	"math"

	"github.com/ready-steady/adapt/algorithm"
	"github.com/ready-steady/adapt/algorithm/internal"
	"github.com/ready-steady/adapt/grid"
)

const (
	none = ^uint(0)
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

	accuracy []float64
	priority []float64

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
		if k == none {
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
	no := self.no
	for i := range self.active.Positions {
		if !self.threshold.Check(self.accuracy[i*no : (i+1)*no]) {
			return false
		}
	}
	return true
}

func (self *Strategy) choose() uint {
	if len(self.active.Positions) == 0 {
		return none
	}
	k, max := none, 0.0
	for i := range self.active.Positions {
		value := self.priority[i]
		if k == none || value > max || (value == max && k > i) {
			k, max = i, value
		}
	}
	if max <= 0.0 {
		return none
	}
	return k
}

func (self *Strategy) consume(state *algorithm.State) {
	ni, no := self.ni, self.no
	nop := uint(len(self.priority))
	noa := uint(len(self.accuracy))
	nnp := uint(len(state.Counts))
	nna := nnp * no

	levels := internal.Levelize(state.Lndices, ni)

	self.priority = append(self.priority, make([]float64, nnp)...)
	priority := self.priority[nop:]

	self.accuracy = append(self.accuracy, make([]float64, nna)...)
	accuracy := self.accuracy[noa:]

	for i, offset := uint(0), uint(0); i < nnp; i++ {
		count := state.Counts[i]
		if levels[i] < uint64(self.lmin) {
			priority[i] = infinity
			for j := uint(0); j < no; j++ {
				accuracy[i*no+j] = infinity
			}
		} else if levels[i] < uint64(self.lmax) {
			priority[i] = internal.Average(state.Scores[offset:(offset + count)])
			self.threshold.Compress(accuracy[i*no:(i+1)*no],
				state.Surpluses[offset*no:(offset+count)*no])
		}
		offset += count
	}

	self.threshold.Update(state.Values)
}

func (self *Strategy) initiate(lndices []uint64, _ *algorithm.Surrogate) (state *algorithm.State) {
	state = &algorithm.State{Lndices: lndices}
	state.Indices, state.Counts = internal.Index(self.guide, lndices, self.ni)
	return
}
