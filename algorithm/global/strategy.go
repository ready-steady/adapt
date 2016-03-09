package global

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/internal"
)

// Strategy controls the interpolation process.
type strategy interface {
	// Start returns the initial level and nodal indices.
	Start() ([]uint64, []uint64, []uint)

	// Check returns true if the interpolation process should continue.
	Check() bool

	// Push takes into account new information.
	Push(*state)

	// Next returns the level and nodal indices for the next iteration.
	Next() ([]uint64, []uint64, []uint)
}

type basicStrategy struct {
	internal.Active

	ni uint
	no uint

	grid Grid

	εa float64
	εr float64

	k uint

	scores []float64
	errors []float64

	lower []float64
	upper []float64
}

func newStrategy(ni, no uint, grid Grid, config *Config) *basicStrategy {
	return &basicStrategy{
		Active: *internal.NewActive(ni, config.MaxLevel, config.MaxIndices),

		ni: ni,
		no: no,

		grid: grid,

		εa: config.AbsoluteError,
		εr: config.RelativeError,

		k: ^uint(0),

		lower: internal.RepeatFloat64(math.Inf(1.0), no),
		upper: internal.RepeatFloat64(math.Inf(-1.0), no),
	}
}

func (self *basicStrategy) Start() (lindices []uint64, indices []uint64, counts []uint) {
	lindices = self.Active.Start()
	indices, counts = internal.Index(self.grid, lindices, self.ni)
	return
}

func (self *basicStrategy) Check() bool {
	no, errors := self.no, self.errors
	ne := uint(len(errors)) / no
	if ne == 0 {
		return true
	}
	δ := threshold(self.lower, self.upper, self.εa, self.εr)
	for i := range self.Positions {
		if i >= ne {
			continue
		}
		for j := uint(0); j < no; j++ {
			if errors[i*no+j] > δ[j] {
				return true
			}
		}
	}
	return false
}

func (self *basicStrategy) Push(state *state) {
	self.updateBounds(state.observations)
	self.scores = append(self.scores, state.scores...)
	self.errors = append(self.errors, error(state.surpluses, state.counts, self.no)...)
}

func (self *basicStrategy) Next() ([]uint64, []uint64, []uint) {
	self.Remove(self.k)
	self.k = internal.LocateMaxFloat64s(self.scores, self.Positions)
	lindices := self.Active.Next(self.k)
	indices, counts := internal.Index(self.grid, lindices, self.ni)
	return lindices, indices, counts
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
