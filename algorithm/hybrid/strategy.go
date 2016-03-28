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
	// Done checks if the stopping criteria have been satisfied.
	Done() bool

	// Score assigns a score to an interpolation element.
	Score(*Element) []float64

	// Next consumes the result of the current iteration and configures the
	// level and nodal indices for the next iteration.
	Next(*State, *external.Surrogate) *State
}

type BasicStrategy struct {
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

func NewStrategy(inputs, outputs uint, grid Grid, config *Config) *BasicStrategy {
	return &BasicStrategy{
		Active: *internal.NewActive(inputs),

		ni: inputs,
		no: outputs,

		grid: grid,

		lmin: config.MinLevel,
		lmax: config.MaxLevel,

		εt: config.TotalError,
		εl: config.LocalError,

		k: ^uint(0),

		hash:     internal.NewHash(inputs),
		unique:   internal.NewUnique(inputs),
		position: make(map[string]uint),
	}
}

func (self *BasicStrategy) Done() bool {
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
	if current == nil {
		next := &State{}
		next.Lindices = self.Active.Next(self.k)
		next.Indices, next.Counts = internal.Index(self.grid, next.Lindices, self.ni)
		return next
	}

	self.consume(current)

	self.Active.Drop(self.k)
	if len(self.Positions) == 0 {
		return nil
	}
	self.k = internal.LocateMax(self.global, self.Positions)
	if self.global[self.k] <= 0.0 {
		return nil
	}

	next := &State{}
	next.Lindices = self.Active.Next(self.k)
	next.Indices, next.Counts = self.index(next.Lindices, surrogate)
	return next
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
				indices = append(indices, self.unique.Distil(self.grid.ChildrenToward(
					surrogate.Indices[k*ni:(k+1)*ni], j))...)
			}
		}
		counts[i] = uint(len(indices))/ni - o
		o += counts[i]
	}

	return indices, counts
}
