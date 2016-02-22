package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

type tracker struct {
	internal.Tracker

	scores []float64
}

func newTracker(ni uint, config *Config) *tracker {
	return &tracker{
		Tracker: *internal.NewTracker(ni, config.MaxLevel, config.MaxIndices),
	}
}

func (self *tracker) pull() []uint64 {
	return self.Forward(internal.LocateMaxFloat64s(self.scores, self.Active))
}

func (self *tracker) push(scores []float64) {
	self.scores = append(self.scores, scores...)
}
