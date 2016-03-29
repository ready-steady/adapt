package hybrid

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

var (
	infinity = math.Inf(1.0)
)

// Strategy controls the interpolation process.
type Strategy interface {
	// First returns the initial state of the first iteration.
	First() *State

	// Continue returns true if the interpolation process should continue.
	Continue(*State, *external.Surrogate) bool

	// Score assigns a score to an interpolation element.
	Score(*Element) []float64

	// Next consumes the result of the current iteration and returns the initial
	// state of the next one.
	Next(*State, *external.Surrogate) *State
}

type BasicStrategy struct {
	*internal.Active

	ni uint
	no uint

	lmin uint
	lmax uint
	εl   float64
	εt   float64

	grid Grid

	k uint

	hash     *internal.Hash
	unique   *internal.Unique
	position map[string]uint

	offset []uint
	global []float64
	local  []float64
}

func NewStrategy(inputs, outputs, minLevel, maxLevel uint,
	localError, totalError float64, grid Grid) *BasicStrategy {

	return &BasicStrategy{
		Active: internal.NewActive(inputs),

		ni: inputs,
		no: outputs,

		lmin: minLevel,
		lmax: maxLevel,
		εl:   localError,
		εt:   totalError,

		grid: grid,
	}
}

func (self *BasicStrategy) First() *State {
	self.k = ^uint(0)
	self.hash = internal.NewHash(self.ni)
	self.unique = internal.NewUnique(self.ni)
	self.position = make(map[string]uint)

	state := &State{}
	state.Lindices = self.Active.First()
	state.Indices, state.Counts = internal.Index(self.grid, state.Lindices, self.ni)

	return state
}

func (self *BasicStrategy) Continue(_ *State, _ *external.Surrogate) bool {
	if self.k == ^uint(0) {
		return true
	}
	ng := uint(len(self.global))
	total := 0.0
	for i := range self.Positions {
		if i >= ng {
			continue
		}
		total += self.global[i]
	}
	return total > self.εt
}

func (self *BasicStrategy) Score(element *Element) []float64 {
	no := self.no
	local := make([]float64, len(element.Volumes))
	for i, m := uint(0), uint(len(element.Surpluses)); i < m; i++ {
		j := i % no
		local[j] = math.Max(local[j], math.Abs(element.Volumes[j]*element.Surpluses[i]))
	}
	return local
}

func (self *BasicStrategy) Next(current *State, surrogate *external.Surrogate) *State {
	self.consume(current)

	self.Active.Drop(self.k)
	if len(self.Positions) == 0 {
		return nil
	}
	self.k = internal.LocateMax(self.global, self.Positions)
	if self.global[self.k] <= 0.0 {
		return nil
	}

	state := &State{}
	state.Lindices = self.Active.Next(self.k)
	state.Indices, state.Counts = self.index(state.Lindices, surrogate)

	return state
}

func (self *BasicStrategy) consume(state *State) {
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
			global[i] = -infinity
			for j := uint(0); j < count; j++ {
				local[o+j] = -infinity
			}
		} else {
			for _, ε := range state.Scores[o:(o + count)] {
				global[i] += ε
			}
		}
		self.position[self.hash.Key(state.Lindices[i*ni:(i+1)*ni])] = ng + i
		o += count
	}
}

func (self *BasicStrategy) index(lindices []uint64,
	surrogate *external.Surrogate) ([]uint64, []uint) {

	ni := self.ni
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

			var from, till uint
			from = self.offset[k]
			if uint(len(self.offset)) > k+1 {
				till = self.offset[k+1]
			} else {
				till = uint(len(self.local))
			}

			for k := from; k < till; k++ {
				if self.local[k] < self.εl {
					continue
				}
				indices = append(indices, self.unique.Distil(self.grid.RefineToward(
					surrogate.Indices[k*ni:(k+1)*ni], j))...)
			}
		}
		counts[i] = uint(len(indices))/ni - o
		o += counts[i]
	}

	return indices, counts
}
