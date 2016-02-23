package global

import (
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Active is a book-keeper of active level indices.
type Active struct {
	internal.Active

	k    uint
	ni   uint
	rate float64

	norms  []uint64
	scores []float64
}

func newActive(ni uint, config *Config) *Active {
	return &Active{
		Active: *internal.NewActive(ni, config.MaxLevel, config.MaxIndices),

		k:    ^uint(0),
		ni:   ni,
		rate: config.AdaptivityRate,

		norms: []uint64{1},
	}
}

func (self *Active) pull() []uint64 {
	k := internal.LocateMinUint64s(self.norms, self.Positions)
	min, max := self.norms[k], internal.MaxUint64s(self.norms)
	if float64(min) > (1.0-self.rate)*float64(max) {
		k = internal.LocateMaxFloat64s(self.scores, self.Positions)
	}

	self.k = k
	norm := self.norms[k] + 1
	indices := self.Advance(k)

	nn := uint(len(indices)) / self.ni
	for i := uint(0); i < nn; i++ {
		self.norms = append(self.norms, norm)
	}

	return indices
}

func (self *Active) push(scores []float64) {
	self.Forget(self.k)
	self.scores = append(self.scores, scores...)
}

func (self *Active) stats() (uint, uint) {
	return self.Current() - 1, self.Previous() + 1
}
