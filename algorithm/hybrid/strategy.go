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

	active *internal.Active
	hash   *internal.Hash
	unique *internal.Unique

	lndices []lndex
	indices []index

	lcursor map[string]uint
	icursor map[string]uint
}

// Guide is a grid-refinement tool of a basic strategy.
type Guide interface {
	grid.Indexer
	grid.RefinerToward
}

type lndex struct {
	score float64
	from  uint
	till  uint
}

type index struct {
	score float64
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

		active: internal.NewActive(inputs),
		hash:   internal.NewHash(inputs),
		unique: internal.NewUnique(inputs),

		lcursor: make(map[string]uint),
		icursor: make(map[string]uint),
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
	return internal.MaxAbsolute(element.Surplus) * element.Volume
}

func (self *Strategy) check() bool {
	total := 0.0
	for i := range self.active.Positions {
		total += self.lndices[i].score
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
	k, max := none, 0.0
	for i := range self.active.Positions {
		if score := self.lndices[i].score; score > max {
			k, max = i, score
		}
	}
	if max <= 0.0 {
		return none
	}
	return k
}

func (self *Strategy) consume(state *algorithm.State) {
	ni := self.ni
	nol, noi := uint(len(self.lndices)), uint(len(self.indices))
	nnl, nni := uint(len(state.Counts)), uint(len(state.Scores))

	levels := internal.Levelize(state.Lndices, ni)

	self.lndices = append(self.lndices, make([]lndex, nnl)...)
	lndices := self.lndices[nol:]

	self.indices = append(self.indices, make([]index, nni)...)
	indices := self.indices[noi:]

	for i, offset := uint(0), uint(0); i < nnl; i++ {
		count := state.Counts[i]

		if levels[i] < uint64(self.lmin) {
			lndices[i].score = infinity
			for j := uint(0); j < count; j++ {
				indices[offset+j].score = infinity
			}
		} else if levels[i] < uint64(self.lmax) {
			for j := uint(0); j < count; j++ {
				lndices[i].score = math.Max(lndices[i].score, state.Scores[offset+j])
				indices[offset+j].score = state.Scores[offset+j]
			}
		}
		lndices[i].from = noi + offset
		lndices[i].till = noi + offset + count

		lndex := state.Lndices[i*ni : (i+1)*ni]
		self.lcursor[self.hash.Key(lndex)] = nol + i
		for j := uint(0); j < count; j++ {
			index := state.Indices[(offset+j)*ni : (offset+j+1)*ni]
			self.icursor[self.hash.Key(index)] = noi + offset + j
		}

		offset += count
	}
}

func (self *Strategy) index(lndices []uint64, surrogate *algorithm.Surrogate) [][]uint64 {
	ni := self.ni
	nn := uint(len(lndices)) / ni
	groups := make([][]uint64, nn)
	for i := uint(0); i < nn; i++ {
		root, lndex := true, lndices[i*ni:(i+1)*ni]
		for j := uint(0); j < ni; j++ {
			level := lndex[j]
			if level == 0 {
				continue
			}
			root = false

			lndex[j] = level - 1
			k, ok := self.lcursor[self.hash.Key(lndex)]
			lndex[j] = level
			if !ok {
				panic("the index set is not admissible")
			}

			for k, m := self.lndices[k].from, self.lndices[k].till; k < m; k++ {
				if self.indices[k].score >= self.εl {
					index := surrogate.Indices[k*ni : (k+1)*ni]
					groups[i] = append(groups[i], self.guide.RefineToward(index, j)...)
				}
			}
		}
		if root {
			groups[i] = append(groups[i], self.guide.Index(lndex)...)
		}
	}
	return groups
}

func (self *Strategy) initiate(lndices []uint64,
	surrogate *algorithm.Surrogate) (state *algorithm.State) {

	groups := self.index(lndices, surrogate)
	nn := uint(len(groups))
	state = &algorithm.State{
		Lndices: lndices,
		Counts:  make([]uint, nn),
		Data:    groups,
	}
	for i := uint(0); i < nn; i++ {
		indices := self.unique.Distil(groups[i])
		state.Indices = append(state.Indices, indices...)
		state.Counts[i] = uint(len(indices)) / self.ni
	}
	return
}
