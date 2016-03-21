package global

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy controls the interpolation process.
type strategy interface {
	// Done checks if the stopping criteria have been satisfied.
	Done() bool

	// Next consumes the result of the current iteration and configures the
	// level and nodal indices for the next iteration.
	Next(*state, *external.Surrogate) *state
}

type basicStrategy struct {
	internal.Active

	ni uint
	no uint

	grid Grid

	lmax uint

	εa float64
	εr float64

	k uint

	global []float64
	local  []float64
	lower  []float64
	upper  []float64
}

func newStrategy(ni, no uint, grid Grid, config *Config) *basicStrategy {
	return &basicStrategy{
		Active: *internal.NewActive(ni),

		ni: ni,
		no: no,

		grid: grid,

		lmax: config.MaxLevel,

		εa: config.AbsoluteError,
		εr: config.RelativeError,

		k: ^uint(0),

		lower: internal.RepeatFloat64(math.Inf(1.0), no),
		upper: internal.RepeatFloat64(math.Inf(-1.0), no),
	}
}

func (self *basicStrategy) Done() bool {
	no := self.no
	nl := uint(len(self.local)) / no
	if nl == 0 {
		return false
	}
	δ := threshold(self.lower, self.upper, self.εa, self.εr)
	for i := range self.Positions {
		if i >= nl {
			continue
		}
		for j := uint(0); j < no; j++ {
			if self.local[i*no+j] > δ[j] {
				return false
			}
		}
	}
	return true
}

func (self *basicStrategy) Next(current *state, _ *external.Surrogate) *state {
	if current != nil {
		self.consume(current)
		if !self.advance() {
			return nil
		}
	}
	next := &state{}
	next.Lindices = self.Active.Next(self.k)
	next.Indices, next.Counts = internal.Index(self.grid, next.Lindices, self.ni)
	return next
}

func (self *basicStrategy) advance() bool {
	for {
		self.Active.Drop(self.k)
		if len(self.Positions) == 0 {
			return false
		}
		self.k = internal.LocateMaxFloat64s(self.global, self.Positions)
		if self.Norms[self.k] < uint64(self.lmax) {
			return true
		}
	}
}

func (self *basicStrategy) consume(state *state) {
	no, nn := self.no, uint(len(state.Counts))
	local := make([]float64, nn*no)
	for i, offset := uint(0), uint(0); i < nn; i++ {
		ns := state.Counts[i] * no
		for j := uint(0); j < ns; j++ {
			k := i*no + j%no
			local[k] = math.Max(local[k], math.Abs(state.Surpluses[offset+j]))
		}
		offset += state.Counts[i]
	}
	self.global = append(self.global, state.Scores...)
	self.local = append(self.local, local...)
	for i, point := range state.Observations {
		j := uint(i) % no
		self.lower[j] = math.Min(self.lower[j], point)
		self.upper[j] = math.Max(self.upper[j], point)
	}
}

func threshold(lower, upper []float64, εa, εr float64) []float64 {
	threshold := make([]float64, len(lower))
	for i := range threshold {
		threshold[i] = math.Max(εr*(upper[i]-lower[i]), εa)
	}
	return threshold
}
