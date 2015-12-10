// Package global provides an algorithm for globally adaptive hierarchical
// interpolation.
package global

import (
	"math"

	"github.com/ready-steady/adapt/algorithm/internal"
)

var (
	infinity = math.Inf(1.0)
)

// Basis is a functional basis.
type Basis interface {
	// Compute evaluates the value of a basis function.
	Compute([]uint64, []float64) float64
}

// Grid is a sparse grid.
type Grid interface {
	// Compute returns the nodes corresponding to a set of indices.
	Compute([]uint64) []float64

	// Index returns the indices of a set of levels.
	Index([]uint8) []uint64
}

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	config Config
	basis  Basis
	grid   Grid
}

// Progress contains information about the interpolation process.
type Progress struct {
	Level       uint8 // Reached level
	Active      uint  // The number of active indices
	Passive     uint  // The number of passive indices
	Evaluations uint  // The number of function evaluations
}

type cursor map[uint]bool

// New creates an interpolator.
func New(grid Grid, basis Basis, config *Config) *Interpolator {
	return &Interpolator{
		config: *config,
		basis:  basis,
		grid:   grid,
	}
}

// Compute constructs an interpolant for a function.
func (self *Interpolator) Compute(target Target) *Surrogate {
	config := &self.config

	ni, no := target.Dimensions()
	nw := config.Workers

	lindices := make([]uint8, 1*ni)
	indices := self.grid.Index(lindices)
	counts := []uint{uint(len(indices)) / ni}

	surrogate := newSurrogate(ni, no)
	accuracy := newAccuracy(no, config)
	tracker := newTracker(ni, config)

	progress := &Progress{Active: 1}
	for {
		target.Monitor(progress)

		nodes := self.grid.Compute(indices)
		values := internal.Invoke(target.Compute, nodes, ni, no, nw)
		surpluses := internal.Subtract(values, internal.Approximate(self.basis,
			surrogate.Indices, surrogate.Surpluses, nodes, ni, no, nw))

		surrogate.push(indices, surpluses)
		accuracy.push(values, surpluses, counts)
		tracker.push(lindices, assess(target.Score, surpluses, counts, no))

		if accuracy.enough(tracker.active) {
			break
		}

		lindices = tracker.pull()

		progress.Active--
		progress.Passive++
		progress.Active += uint(len(lindices)) / ni

		indices, counts = indices[:0], counts[:0]
		for len(lindices) > 0 {
			newIndices := self.grid.Index(lindices[:ni])
			indices = append(indices, newIndices...)
			counts = append(counts, uint(len(newIndices))/ni)
			lindices = lindices[ni:]
		}

		level := maxUint8(lindices)
		if level > progress.Level {
			progress.Level = level
		}

		progress.Evaluations += uint(len(indices)) / ni
		if progress.Evaluations > config.MaxEvaluations {
			break
		}
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
}
