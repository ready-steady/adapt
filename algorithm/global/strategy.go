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
	δ := make([]float64, no)
	for i := uint(0); i < no; i++ {
		δ[i] = math.Max(self.εr*(self.upper[i]-self.lower[i]), self.εa)
	}
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
	if current == nil {
		next := &state{}
		next.Lindices = self.Active.Next(self.k)
		next.Indices, next.Counts = internal.Index(self.grid, next.Lindices, self.ni)
		return next
	}

	self.consume(current)

	for {
		self.Active.Drop(self.k)
		if len(self.Positions) == 0 {
			return nil
		}
		self.k = internal.LocateMaxFloat64s(self.global, self.Positions)
		if self.Norms[self.k] < uint64(self.lmax) {
			break
		}
	}

	next := &state{}
	next.Lindices = self.Active.Next(self.k)
	next.Indices, next.Counts = internal.Index(self.grid, next.Lindices, self.ni)
	return next
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
