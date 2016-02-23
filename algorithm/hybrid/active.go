package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

type active struct {
	internal.Active

	k uint

	scores []float64
}

func newActive(ni uint, config *Config) *active {
	return &active{
		Active: *internal.NewActive(ni, config.MaxLevel, config.MaxIndices),

		k: ^uint(0),
	}
}

func (self *active) pull() []uint64 {
	self.k = internal.LocateMaxFloat64s(self.scores, self.Positions)
	return self.Advance(self.k)
}

func (self *active) push(scores []float64) {
	self.Forget(self.k)
	self.scores = append(self.scores, scores...)
}

func (self *active) stats() (uint, uint) {
	return self.Current() - 1, self.Previous() + 1
}
