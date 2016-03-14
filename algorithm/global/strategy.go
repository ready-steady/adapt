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

	// Next sets up the level and nodal indices for the next iteration.
	Next(*state)

	// Push consumes the result of the current iteration.
	Push(*state)
}

type basicStrategy struct {
	internal.Active

	ni uint
	no uint

	grid      Grid
	surrogate *external.Surrogate

	εa float64
	εr float64

	k uint

	scores []float64
	errors []float64

	lower []float64
	upper []float64
}

func newStrategy(ni, no uint, grid Grid, surrogate *external.Surrogate,
	config *Config) *basicStrategy {

	return &basicStrategy{
		Active: *internal.NewActive(ni, config.MaxLevel, config.MaxIndices),

		ni: ni,
		no: no,

		grid:      grid,
		surrogate: surrogate,

		εa: config.AbsoluteError,
		εr: config.RelativeError,

		k: ^uint(0),

		lower: internal.RepeatFloat64(math.Inf(1.0), no),
		upper: internal.RepeatFloat64(math.Inf(-1.0), no),
	}
}

func (self *basicStrategy) Done() bool {
	no, errors := self.no, self.errors
	ne := uint(len(errors)) / no
	if ne == 0 {
		return false
	}
	δ := threshold(self.lower, self.upper, self.εa, self.εr)
	for i := range self.Positions {
		if i >= ne {
			continue
		}
		for j := uint(0); j < no; j++ {
			if errors[i*no+j] > δ[j] {
				return false
			}
		}
	}
	return true
}

func (self *basicStrategy) Next(state *state) {
	if state.lindices == nil {
		state.lindices = self.Active.Start()
		state.indices, state.counts = internal.Index(self.grid, state.lindices, self.ni)
	} else {
		self.Remove(self.k)
		self.k = internal.LocateMaxFloat64s(self.scores, self.Positions)
		state.lindices = self.Active.Next(self.k)
		state.indices, state.counts = internal.Index(self.grid, state.lindices, self.ni)
	}
}

func (self *basicStrategy) Push(state *state) {
	self.surrogate.Push(state.indices, state.surpluses, state.volumes)
	self.updateBounds(state.observations)
	self.scores = append(self.scores, state.scores...)
	self.errors = append(self.errors, error(state.surpluses, state.counts, self.no)...)
}

func (self *basicStrategy) updateBounds(observations []float64) {
	no := self.no
	for i, point := range observations {
		j := uint(i) % no
		self.lower[j] = math.Min(self.lower[j], point)
		self.upper[j] = math.Max(self.upper[j], point)
	}
}

func error(surpluses []float64, counts []uint, no uint) []float64 {
	nn := uint(len(counts))
	errors := make([]float64, nn*no)
	for i := uint(0); i < nn; i++ {
		ns := counts[i] * no
		for j := uint(0); j < ns; j++ {
			k := i*no + j%no
			errors[k] = math.Max(errors[k], math.Abs(surpluses[j]))
		}
		surpluses = surpluses[ns:]
	}
	return errors
}

func threshold(lower, upper []float64, εa, εr float64) []float64 {
	threshold := make([]float64, len(lower))
	for i := range threshold {
		threshold[i] = math.Max(εr*(upper[i]-lower[i]), εa)
	}
	return threshold
}
