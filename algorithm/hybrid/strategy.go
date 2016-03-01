package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy guides the interpolation process.
type Strategy interface {
	// Continue decides if the interpolation process should go on.
	Continue(*external.Active) bool

	// Push takes into account a new interpolation element and its score.
	Push(*Element, []float64)

	// Select chooses an active index for refinement.
	Select(*external.Active) uint
}

type defaultStrategy struct {
	ni uint
	no uint

	εt float64
	εl float64

	global []float64
	local  []float64
}

func newStrategy(ni, no uint, total, local float64) *defaultStrategy {
	return &defaultStrategy{
		ni: ni,
		no: no,

		εt: total,
		εl: local,
	}
}

func (self *defaultStrategy) Continue(active *external.Active) bool {
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

func (self *defaultStrategy) Select(active *external.Active) uint {
	return internal.LocateMaxFloat64s(self.global, active.Positions)
}
