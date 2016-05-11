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

	global []float64
	local  []float64

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

func (self *Strategy) First() *algorithm.State {
	state := &algorithm.State{}
	state.Lndices = self.active.First()
	state.Indices, state.Counts = internal.Index(self.guide, state.Lndices, self.ni)
	return state
}

func (self *Strategy) Next(state *algorithm.State, _ *algorithm.Surrogate) *algorithm.State {
	for {
		self.consume(state)
		if self.check() {
			return nil
		}
		k := self.choose()
		if k == none {
			return nil
		}
		state = &algorithm.State{}
		state.Lndices = self.active.Next(k)
		state.Indices, state.Counts = internal.Index(self.guide, state.Lndices, self.ni)
		if len(state.Indices) > 0 {
			return state
		}
	}
}

func (self *Strategy) Score(element *algorithm.Element) float64 {
	return internal.SumAbsolute(element.Surplus)
}

func (self *Strategy) check() bool {
	no, δ := self.no, self.threshold.Values
	for i := range self.active.Positions {
		for j := uint(0); j < no; j++ {
			if self.local[i*no+j] > δ[j] {
				return false
			}
		}
	}
	return true
}

func (self *Strategy) choose() uint {
	if len(self.active.Positions) == 0 {
		return none
	}
	k := internal.LocateMax(self.global, self.active.Positions)
	if self.global[k] <= 0.0 {
		return none
	}
	return k
}

func (self *Strategy) consume(state *algorithm.State) {
	no, ng, nl := self.no, uint(len(self.global)), uint(len(self.local))
	nn := uint(len(state.Counts))

	levels := internal.Levelize(state.Lndices, self.ni)

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
