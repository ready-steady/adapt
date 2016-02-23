package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

type tracker struct {
	internal.Active

	k uint

	scores []float64
}

func newTracker(ni uint, config *Config) *tracker {
	return &tracker{
		Active: *internal.NewActive(ni, config.MaxLevel, config.MaxIndices),

		k: ^uint(0),
	}
}

func (self *tracker) pull() []uint64 {
	self.k = internal.LocateMaxFloat64s(self.scores, self.Positions)
	return self.Advance(self.k)
}

func (self *tracker) push(scores []float64) {
	self.Forget(self.k)
	self.scores = append(self.scores, scores...)
}

func (self *tracker) stats() (uint, uint) {
	return self.Current() - 1, self.Previous() + 1
}
