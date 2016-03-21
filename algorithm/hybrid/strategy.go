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
	Next(*state, *external.Surrogate) *state
}

type basicStrategy struct {
	internal.Active

	ni uint
	no uint

	grid Grid

	lmin uint
	lmax uint

	εt float64
	εl float64

	k uint

	hash     *internal.Hash
	unique   *internal.Unique
	position map[string]uint

	offset []uint
	global []float64
	local  []float64
}

func newStrategy(ni, no uint, grid Grid, config *Config) *basicStrategy {
	return &basicStrategy{
		Active: *internal.NewActive(ni),

		ni: ni,
		no: no,

		grid: grid,

		lmin: config.MinLevel,
		lmax: config.MaxLevel,

		εt: config.TotalError,
		εl: config.LocalError,

		k: ^uint(0),

		hash:     internal.NewHash(ni),
		unique:   internal.NewUnique(ni),
		position: make(map[string]uint),
	}
}

func (self *basicStrategy) Done() bool {
	ng := uint(len(self.global))
	if ng == 0 {
		return false
	}
	total := 0.0
	for i := range self.Positions {
		if i >= ng {
			continue
		}
		total += self.global[i]
	}
	return total <= self.εt
}

func (self *basicStrategy) Next(current *state, surrogate *external.Surrogate) *state {
	next := &state{}
	if current == nil {
		next.Lindices = self.Active.Next(self.k)
		next.Indices, next.Counts = internal.Index(self.grid, next.Lindices, self.ni)
	} else {
		self.consume(current)
		if !self.advance() {
			return nil
		}
		next.Lindices = self.Active.Next(self.k)
		next.Indices, next.Counts = self.indexSet(next.Lindices, surrogate)
	}
	return next
}

func (self *basicStrategy) advance() bool {
	for {
		self.Active.Drop(self.k)
		if len(self.Positions) == 0 {
			return false
		}
		self.k = internal.LocateMaxFloat64s(self.global, self.Positions)
		if self.Norms[self.k] < uint64(self.lmax) {
			return true
		}
	}
}

func (self *basicStrategy) consume(state *state) {
	ni, ng, nl := self.ni, uint(len(self.global)), uint(len(self.local))
	nn := uint(len(state.Counts))
	for i, offset := uint(0), uint(0); i < nn; i++ {
		global := 0.0
		for _, ε := range state.Scores[offset:(offset + state.Counts[i])] {
			global += ε
		}
		self.position[self.hash.Key(state.Lindices[i*ni:(i+1)*ni])] = ng + i
		self.offset = append(self.offset, nl+offset)
		self.global = append(self.global, global)
		offset += state.Counts[i]
	}
	self.local = append(self.local, state.Scores...)
}

func (self *basicStrategy) indexSet(lindices []uint64,
	surrogate *external.Surrogate) ([]uint64, []uint) {

	ni := self.ni
	nn := uint(len(lindices)) / ni
	indices, counts := []uint64(nil), make([]uint, nn)
	for i, offset := uint(0), uint(0); i < nn; i++ {
		indices = append(indices, self.indexOne(lindices[i*ni:(i+1)*ni], surrogate)...)
		counts[i] = uint(len(indices))/ni - offset
		offset += counts[i]
	}
	return indices, counts
}

func (self *basicStrategy) indexOne(lindex []uint64, surrogate *external.Surrogate) []uint64 {
	ni := self.ni

	indices := []uint64(nil)
	for i := uint(0); i < ni; i++ {
		level := lindex[i]
		if level == 0 {
			continue
		}

		lindex[i] = level - 1
		k, ok := self.position[self.hash.Key(lindex)]
		lindex[i] = level
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

		for j := from; j < till; j++ {
			if self.local[j] < self.εl {
				continue
			}
			indices = append(indices, self.unique.Distil(self.grid.ChildrenToward(
				surrogate.Indices[j*ni:(j+1)*ni], i))...)
		}
	}

	return indices
}
