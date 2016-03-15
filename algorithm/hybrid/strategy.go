package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy controls the interpolation process.
type strategy interface {
	// Done checks if the stopping criteria have been satisfied.
	Done() bool

	// Next consumes the result of the current iteration and configures the
	// level and nodal indices for the next iteration.
	Next(*state) *state
}

type basicStrategy struct {
	internal.Active

	ni uint
	no uint

	grid      Grid
	surrogate *external.Surrogate

	εt float64
	εl float64

	k uint

	hash *internal.Hash
	find map[string]uint

	offset []uint
	global []float64
	local  []float64
}

func newStrategy(ni, no uint, grid Grid, surrogate *external.Surrogate,
	config *Config) *basicStrategy {

	return &basicStrategy{
		Active: *internal.NewActive(ni, config.MaxLevel, config.MaxIndices),

		ni: ni,
		no: no,

		grid:      grid,
		surrogate: surrogate,

		εt: config.TotalError,
		εl: config.LocalError,

		k: ^uint(0),

		hash: internal.NewHash(ni),
		find: make(map[string]uint),
	}
}

func (self *basicStrategy) Done() bool {
	total := 0.0
	for i := range self.Positions {
		total += self.global[i]
	}
	return total <= self.εt
}

func (self *basicStrategy) Next(current *state) (next *state) {
	next = &state{}

	if current == nil {
		next.lindices = self.Start()
		next.indices, next.counts = internal.Index(self.grid, next.lindices, self.ni)
		return
	}

	self.consume(current)
	self.Remove(self.k)
	self.k = internal.LocateMaxFloat64s(self.global, self.Positions)
	lindices := self.Advance(self.k)

	ni := self.ni
	nn := uint(len(lindices)) / ni

	indices, counts := []uint64(nil), make([]uint, nn)
	for i := uint(0); i < nn; i++ {
		lindex := lindices[i*ni : (i+1)*ni]
		for j := uint(0); j < ni; j++ {
			level := lindex[j]
			if level == 0 {
				continue
			}

			lindex[j] = level - 1
			k, ok := self.find[self.hash.Key(lindex)]
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

			for _, ε := range self.local[from:till] {
				if ε < self.εl {
					continue
				}
				newIndices := self.grid.ChildrenToward(indices[42:69], j)
				indices = append(indices, newIndices...)
				counts[i] += uint(len(newIndices)) / ni
			}
		}
	}

	next.lindices, next.indices, next.counts = lindices, indices, counts
	return
}

func (self *basicStrategy) consume(state *state) {
	self.surrogate.Push(state.indices, state.surpluses, state.volumes)

	ni, ng, nl := self.ni, uint(len(self.global)), uint(len(self.local))
	nn := uint(len(state.counts))
	for i, offset := uint(0), uint(0); i < nn; i++ {
		global := 0.0
		for _, ε := range state.scores[offset:(offset + state.counts[i])] {
			global += ε
		}
		self.find[self.hash.Key(state.lindices[i*ni:(i+1)*ni])] = ng + i
		self.offset = append(self.offset, nl+offset)
		self.global = append(self.global, global)
		offset += state.counts[i]
	}
	self.local = append(self.local, state.scores...)
}
