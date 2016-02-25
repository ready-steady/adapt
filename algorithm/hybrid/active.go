package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Active is a book-keeper of active level indices.
type Active struct {
	external.Active

	Global []float64 // Error indicators of level indices

	k  uint
	ni uint
}

func newActive(ni uint, config *Config) *Active {
	return &Active{
		Active: *external.NewActive(ni, config.MaxLevel, config.MaxIndices),

		k:  ^uint(0),
		ni: ni,
	}
}

func (self *Active) pull() []uint64 {
	self.k = internal.LocateMaxFloat64s(self.Global, self.Positions)

	indices := self.Advance(self.k)

	nn := uint(len(indices)) / self.ni
	for i := uint(0); i < nn; i++ {
		self.Global = append(self.Global, infinity)
	}

	return indices
}

func (self *Active) push(global []float64) {
	copy(self.Global[len(self.Global)-len(global):], global)
	self.Forget(self.k)
}
