package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy guides the interpolation process.
type strategy interface {
	// Begin returns the initial order indices.
	Begin() ([]uint64, []uint)

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

		εt: config.TotalError,
		εl: config.LocalError,

		k: ^uint(0),
	}
}

func (self *basicStrategy) Begin() ([]uint64, []uint) {
	return index(self.grid, self.Active.Initialize(), self.ni)
}

func (self *basicStrategy) Check() bool {
	total := 0.0
	for i := range self.Positions {
		total += self.global[i]
	}
	return total > self.εt
}

func (self *basicStrategy) Push(element *Element, local []float64) {
	global := 0.0
	for i := range local {
		global += local[i]
	}
	self.global = append(self.global, global)
	self.local = append(self.local, local...)
}

func (self *basicStrategy) Move() ([]uint64, []uint) {
	self.Remove(self.k)
	self.k = internal.LocateMaxFloat64s(self.global, self.Positions)
	return index(self.grid, self.Active.Forward(self.k), self.ni)
}
