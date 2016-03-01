package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy guides the interpolation process.
type Strategy interface {
	// Continue decides if the interpolation process should go on.
	Continue(*internal.Active) bool

	// Push takes into account a new interpolation element and its score.
	Push(*Element, []float64)

	// Forward selects an active index for refinement and returns its admissible
	// neighbors from the forward neighborhood.
	Forward(*internal.Active) []uint64
}

type defaultStrategy struct {
	ni uint
	no uint

	εt float64
	εl float64

	k uint

	global []float64
	local  []float64
}

func newStrategy(ni, no uint, total, local float64) *defaultStrategy {
	return &defaultStrategy{
		ni: ni,
		no: no,

		εt: total,
		εl: local,

		k: ^uint(0),
	}
}

func (self *defaultStrategy) Continue(active *internal.Active) bool {
	total := 0.0
	for i := range active.Positions {
		total += self.global[i]
	}
	return total > self.εt
}

func (self *defaultStrategy) Push(element *Element, local []float64) {
	global := 0.0
	for i := range local {
		global += local[i]
	}
	self.global = append(self.global, global)
	self.local = append(self.local, local...)
}

func (self *defaultStrategy) Forward(active *internal.Active) []uint64 {
	active.Remove(self.k)
	self.k = internal.LocateMaxFloat64s(self.global, active.Positions)
	return active.Forward(self.k)
}
