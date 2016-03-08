package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy guides the interpolation process.
type strategy interface {
	// Start returns the initial level and nodal indices.
	Start() ([]uint64, []uint64, []uint)

	// Check decides if the interpolation process should go on.
	Check() bool

	// Push takes into account new level indices, nodal indices, function
	// values, hierarchical surpluses, and scores.
	Push([]uint64, []uint64, []float64, []float64, []float64, []uint)

	// Move selects an active level index, searches admissible level indices in
	// the forward neighborhood of the selected level index, searches admissible
	// nodal indices with respect to each admissible level index, and returns
	// all the identified level and nodal indices.
	Move() ([]uint64, []uint64, []uint)
}

type basicStrategy struct {
	internal.Active

	ni uint
	no uint

	grid Grid

	εt float64
	εl float64

	k uint

	hash *internal.Hash
	find map[string]uint

	offset []uint
	global []float64
	local  []float64
}

func newStrategy(ni, no uint, grid Grid, config *Config) *basicStrategy {
	return &basicStrategy{
		Active: *internal.NewActive(ni, config.MaxLevel, config.MaxIndices),

		ni: ni,
		no: no,

		grid: grid,

		εt: config.TotalError,
		εl: config.LocalError,

		k: ^uint(0),

		hash: internal.NewHash(ni),
		find: make(map[string]uint),
	}
}

func (self *basicStrategy) Start() ([]uint64, []uint64, []uint) {
	lindices := self.Active.Start()
	indices, counts := internal.Index(self.grid, lindices, self.ni)
	return lindices, indices, counts
}

func (self *basicStrategy) Check() bool {
	total := 0.0
	for i := range self.Positions {
		total += self.global[i]
	}
	return total > self.εt
}

func (self *basicStrategy) Push(lindices, _ []uint64, _, _, local []float64, counts []uint) {
	ni, ng, nl := self.ni, uint(len(self.global)), uint(len(self.local))
	nn := uint(len(counts))
	for i, offset := uint(0), uint(0); i < nn; i++ {
		global := 0.0
		for _, ε := range local[offset:(offset + counts[i])] {
			global += ε
		}
		self.find[self.hash.Key(lindices[i*ni:(i+1)*ni])] = ng + i
		self.offset = append(self.offset, nl+offset)
		self.global = append(self.global, global)
		offset += counts[i]
	}
	self.local = append(self.local, local...)
}

func (self *basicStrategy) Move() ([]uint64, []uint64, []uint) {
	self.Remove(self.k)
	self.k = internal.LocateMaxFloat64s(self.global, self.Positions)
	lindices := self.Active.Move(self.k)

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
			}
		}
	}

	return lindices, indices, counts
}
