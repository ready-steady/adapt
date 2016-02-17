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

	initialized bool
}

func newTracker(ni uint, config *Config) *tracker {
	return &tracker{
		Tracker: *internal.NewTracker(ni, config.MaxLevel, config.MaxIndices),

		ni:   ni,
		rate: config.AdaptivityRate,

		norms:  make([]uint64, 1),
		scores: make([]float64, 0),
	}
}

func (self *tracker) pull() []uint64 {
	if !self.initialized {
		self.initialized = true
		return self.Forward(^uint(0))
	}

	min, k := minUint64Set(self.norms, self.Active)
	max := maxUint64(self.norms)
	if float64(min) > (1.0-self.rate)*float64(max) {
		_, k = maxFloat64Set(self.scores, self.Active)
	}

	indices := self.Forward(k)

	nn := uint(len(indices)) / self.ni
	norm := self.norms[k] + 1
	for i := uint(0); i < nn; i++ {
		self.norms = append(self.norms, norm)
	}

	return indices
}

func (self *tracker) push(score float64) {
	self.scores = append(self.scores, score)
}
