package global

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

type tracker struct {
	internal.Tracker

	ni   uint
	rate float64

	norms  []uint64
	scores []float64
}

func newTracker(ni uint, config *Config) *tracker {
	return &tracker{
		Tracker: *internal.NewTracker(ni, config.MaxLevel, config.MaxIndices),

		ni:   ni,
		rate: config.AdaptivityRate,
	}
}

func (self *tracker) pull() (indices []uint64) {
	var norm uint64

	if self.norms == nil {
		norm, indices = 0, self.Forward(^uint(0))
	} else {
		k := internal.LocateMinUint64s(self.norms, self.Active)
		min, max := self.norms[k], internal.MaxUint64s(self.norms)
		if float64(min) > (1.0-self.rate)*float64(max) {
			k = internal.LocateMaxFloat64s(self.scores, self.Active)
		}
		norm, indices = self.norms[k]+1, self.Forward(k)
	}

	nn := uint(len(indices)) / self.ni
	for i := uint(0); i < nn; i++ {
		self.norms = append(self.norms, norm)
	}

	return
}

func (self *tracker) push(scores []float64) {
	self.scores = append(self.scores, scores...)
}
