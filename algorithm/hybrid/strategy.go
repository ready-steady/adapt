package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy guides the interpolation process.
type Strategy interface {
	// Continue decides if the interpolation process should go on.
	Continue(*external.Active) bool

	// Push takes into account a new location and its score.
	Push(*Location, float64, []float64)

	// Select chooses an active index for refinement.
	Select(*external.Active) uint
}

type defaultStrategy struct {
	ni uint
	no uint

	εg float64
	εl float64

	global []float64
	local  []float64
}

func newStrategy(ni, no uint, global, local float64) *defaultStrategy {
	return &defaultStrategy{
		ni: ni,
		no: no,

		εg: global,
		εl: local,
	}
}

func (self *defaultStrategy) Continue(active *external.Active) bool {
	Σ := 0.0
	for i := range active.Positions {
		Σ += self.global[i]
	}
	return Σ > self.εg
}

func (self *defaultStrategy) Push(location *Location, global float64, local []float64) {
	self.global = append(self.global, global)
	self.local = append(self.local, local...)
}

func (self *defaultStrategy) Select(active *external.Active) uint {
	return internal.LocateMaxFloat64s(self.global, active.Positions)
}
