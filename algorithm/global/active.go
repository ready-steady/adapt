package global

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Active is a book-keeper of active level indices.
type Active struct {
	external.Active

	Norms  []uint64  // Norms of level indices
	Scores []float64 // Scores of level indices

	k    uint
	ni   uint
	rate float64
}

func newActive(ni uint, config *Config) *Active {
	return &Active{
		Active: *external.NewActive(ni, config.MaxLevel, config.MaxIndices),

		k:    ^uint(0),
		ni:   ni,
		rate: config.AdaptivityRate,

		Norms:  []uint64{1},
		Scores: []float64{0.0},
	}
}

func (self *Active) pull() []uint64 {
	k := internal.LocateMinUint64s(self.Norms, self.Positions)
	min, max := self.Norms[k], internal.MaxUint64s(self.Norms)
	if float64(min) > (1.0-self.rate)*float64(max) {
		k = internal.LocateMaxFloat64s(self.Scores, self.Positions)
	}

	self.k = k
	norm := self.Norms[k] + 1
	indices := self.Advance(k)

	nn := uint(len(indices)) / self.ni
	for i := uint(0); i < nn; i++ {
		self.Norms = append(self.Norms, norm)
		self.Scores = append(self.Scores, infinity)
	}

	return indices
}

func (self *Active) push(scores []float64) {
	copy(self.Scores[len(self.Scores)-len(scores):], scores)
	self.Forget(self.k)
}
