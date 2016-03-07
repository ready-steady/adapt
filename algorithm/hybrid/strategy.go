package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy guides the interpolation process.
type strategy interface {
	// Start returns the initial order indices.
	Start() ([]uint64, []uint)

	// Check decides if the interpolation process should go on.
	Check() bool

	// Push takes into account a new interpolation element and its score.
	Push(*Element, []float64)

	// Move selects an active level index, searches admissible level indices in
	// the forward neighborhood of the selected level index, searches admissible
	// order indices with respect to each admissible level index, and returns
	// all the identified order indices.
	Move() ([]uint64, []uint)
}

type basicStrategy struct {
	internal.Active

	ni uint
	no uint

	grid Grid

	hash *internal.Hash
	find map[string]uint

	εt float64
	εl float64

	k uint

	global []float64
	local  []float64
}

func newStrategy(ni, no uint, grid Grid, config *Config) *basicStrategy {
	return &basicStrategy{
		Active: *internal.NewActive(ni, config.MaxLevel, config.MaxIndices),

		ni: ni,
		no: no,

		grid: grid,

		hash: internal.NewHash(ni),
		find: make(map[string]uint),

		εt: config.TotalError,
		εl: config.LocalError,

		k: ^uint(0),
	}
}

func (self *basicStrategy) Start() ([]uint64, []uint) {
	return internal.Index(self.grid, self.Active.Start(), self.ni)
}

func (self *basicStrategy) Check() bool {
	total := 0.0
	for i := range self.Positions {
		total += self.global[i]
	}
	return total > self.εt
}

func (self *basicStrategy) Push(element *Element, local []float64) {
	ni, ne := self.ni, uint(len(self.local))
	nn := uint(len(element.Indices)) / self.ni

	global := 0.0
	for i := uint(0); i < nn; i++ {
		global += local[i]
		self.find[self.hash.Key(element.Indices[i*ni:(i+1)*ni])] = ne + i
	}

	self.global = append(self.global, global)
	self.local = append(self.local, local...)
}

func (self *basicStrategy) Move() ([]uint64, []uint) {
	self.Remove(self.k)
	self.k = internal.LocateMaxFloat64s(self.global, self.Positions)
	lindices := self.Active.Move(self.k)

	ni := self.ni
	nn := uint(len(lindices)) / ni

	indices, counts := []uint64(nil), make([]uint, nn)
	for i := uint(0); i < nn; i++ {
		lindex := lindices[i*ni : (i+1)*ni]
		for j := uint(0); j < ni; j++ {
			l := lindex[j]
			if l == 0 {
				continue
			}
			lindex[j] = l - 1
			_, ok := self.find[self.hash.Key(lindex)]
			lindex[j] = l
			if !ok {
				continue
			}
		}
	}

	return indices, counts
}
