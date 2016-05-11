package hybrid

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
	εl   float64
	εt   float64

	k uint

	active   *internal.Active
	hash     *internal.Hash
	unique   *internal.Unique
	position map[string]uint

	offset []uint
	global []float64
	local  []float64
}

// Guide is a grid-refinement tool of a basic strategy.
type Guide interface {
	grid.Indexer
	grid.RefinerToward
}

// NewStrategy creates a basic strategy.
func NewStrategy(inputs, outputs uint, guide Guide, minLevel, maxLevel uint,
	localError, totalError float64) *Strategy {

	return &Strategy{
		ni: inputs,
		no: outputs,

		guide: guide,

		lmin: minLevel,
		lmax: maxLevel,
		εl:   localError,
		εt:   totalError,

		active:   internal.NewActive(inputs),
		hash:     internal.NewHash(inputs),
		unique:   internal.NewUnique(inputs),
		position: make(map[string]uint),
	}
}

func (self *Strategy) First() *algorithm.State {
	state := &algorithm.State{}
	state.Lindices = self.active.First()
	state.Indices, state.Counts = internal.Index(self.guide, state.Lindices, self.ni)
	return state
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
		state = &algorithm.State{}
		state.Lindices = self.active.Next(k)
		state.Indices, state.Counts = self.index(state.Lindices, surrogate)
		if len(state.Indices) > 0 {
			return state
		}
	}
}

func (self *Strategy) Score(element *algorithm.Element) float64 {
	return internal.MaxAbsolute(element.Surplus) * element.Volume
}

func (self *Strategy) check() bool {
	total := 0.0
	for i := range self.active.Positions {
		total += self.global[i]
		if total > self.εt {
			return false
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
	ni, ng, nl := self.ni, uint(len(self.global)), uint(len(self.local))
	nn := uint(len(state.Counts))

	levels := internal.Levelize(state.Lindices, ni)

	self.offset = append(self.offset, make([]uint, nn)...)
	offset := self.offset[ng:]

	self.global = append(self.global, make([]float64, nn)...)
	global := self.global[ng:]

	self.local = append(self.local, state.Scores...)
	local := self.local[nl:]

	for i, o := uint(0), uint(0); i < nn; i++ {
		count := state.Counts[i]
		offset[i] = nl + o
		if levels[i] < uint64(self.lmin) {
			global[i] = infinity
			for j := uint(0); j < count; j++ {
				local[o+j] = infinity
			}
		} else if levels[i] >= uint64(self.lmax) {
			global[i] = 0.0
			for j := uint(0); j < count; j++ {
				local[o+j] = 0.0
			}
		} else {
			global[i] = internal.Max(state.Scores[o:(o + count)])
		}
		self.position[self.hash.Key(state.Lindices[i*ni:(i+1)*ni])] = ng + i
		o += count
	}
}

func (self *Strategy) index(lindices []uint64, surrogate *algorithm.Surrogate) ([]uint64, []uint) {
	ni, nl := self.ni, uint(len(self.local))
	nn := uint(len(lindices)) / ni

	indices, counts := []uint64(nil), make([]uint, nn)
	for i, o := uint(0), uint(0); i < nn; i++ {
		lindex := lindices[i*ni : (i+1)*ni]
		for j := uint(0); j < ni; j++ {
			level := lindex[j]
			if level == 0 {
				continue
			}

			lindex[j] = level - 1
			k, ok := self.position[self.hash.Key(lindex)]
			lindex[j] = level
			if !ok {
				continue
			}

			from, till := self.offset[k], nl
			if uint(len(self.offset)) > k+1 {
				till = self.offset[k+1]
			}

			for k := from; k < till; k++ {
				if self.local[k] >= self.εl {
					indices = append(indices, self.unique.Distil(self.guide.RefineToward(
						surrogate.Indices[k*ni:(k+1)*ni], j))...)
				}
			}
		}
		counts[i] = uint(len(indices))/ni - o
		o += counts[i]
	}

	return indices, counts
}
