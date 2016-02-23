package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

type Active struct {
	internal.Active

	k uint

	scores []float64
}

func newActive(ni uint, config *Config) *Active {
	return &Active{
		Active: *internal.NewActive(ni, config.MaxLevel, config.MaxIndices),

		k: ^uint(0),
	}
}

func (self *Active) pull() []uint64 {
	self.k = internal.LocateMaxFloat64s(self.scores, self.Positions)
	return self.Advance(self.k)
}

func (self *Active) push(scores []float64) {
	self.Forget(self.k)
	self.scores = append(self.scores, scores...)
}

func (self *Active) stats() (uint, uint) {
	return self.Current() - 1, self.Previous() + 1
}
