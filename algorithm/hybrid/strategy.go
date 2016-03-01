package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy guides the interpolation process.
type Strategy interface {
	// Continue decides if the interpolation process should go on.
	Continue() bool

	// Push takes into account a new interpolation element and its score.
	Push(*Element, []float64)

	// Forward selects an active index for refinement and returns its admissible
	// neighbors from the forward neighborhood.
	Forward() []uint64
}

type defaultStrategy struct {
	internal.Active

	ni uint
	no uint

	εt float64
	εl float64

	k uint

	global []float64
	local  []float64
}

func newStrategy(ni, no uint, config *Config) *defaultStrategy {
	return &defaultStrategy{
		Active: *internal.NewActive(ni, config.MaxLevel, config.MaxIndices),

		ni: ni,
		no: no,

		εt: config.TotalError,
		εl: config.LocalError,

		k: ^uint(0),
	}
}

func (self *defaultStrategy) Continue() bool {
	total := 0.0
	for i := range self.Positions {
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

func (self *defaultStrategy) Forward() []uint64 {
	self.Remove(self.k)
	self.k = internal.LocateMaxFloat64s(self.global, self.Positions)
	return self.Active.Forward(self.k)
}
