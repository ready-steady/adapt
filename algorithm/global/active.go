package global

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Active is a book-keeper of active level indices.
type Active struct {
	external.Active

	Scores []float64 // Scores of level indices

	k  uint
	ni uint
}

func newActive(ni uint, config *Config) *Active {
	return &Active{
		Active: *external.NewActive(ni, config.MaxLevel, config.MaxIndices),

		k:  ^uint(0),
		ni: ni,

		Scores: []float64{0.0},
	}
}

func (self *Active) pull() []uint64 {
	self.k = internal.LocateMaxFloat64s(self.Scores, self.Positions)
	indices := self.Advance(self.k)

	nn := uint(len(indices)) / self.ni
	for i := uint(0); i < nn; i++ {
		self.Scores = append(self.Scores, infinity)
	}

	return indices
}

func (self *Active) push(scores []float64) {
	copy(self.Scores[len(self.Scores)-len(scores):], scores)
	self.Forget(self.k)
}
