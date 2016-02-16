// Package global provides an algorithm for hierarchical interpolation with
// global adaptation.
package global

import (
	"github.com/ready-steady/adapt/algorithm/internal"
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
	Index([]uint64) []uint64
}

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	grid   Grid
	basis  Basis
	config Config
}

// Progress contains information about the interpolation process.
type Progress struct {
	Level       uint // Reached level
	Active      uint // The number of active indices
	Passive     uint // The number of passive indices
	Evaluations uint // The number of function evaluations
}

// New creates an interpolator.
func New(grid Grid, basis Basis, config *Config) *Interpolator {
	return &Interpolator{
		grid:   grid,
		basis:  basis,
		config: *config,
	}
}

// Compute constructs an interpolant for a function.
func (self *Interpolator) Compute(target Target, metric Metric) *Surrogate {
	config := &self.config

	ni, no := target.Dimensions()
	nw := config.Workers

	surrogate := newSurrogate(ni, no)
	tracker := newTracker(ni, config)

	progress := &Progress{}
	for {
		lindices := tracker.Pull()

		progress.Active = uint(len(tracker.Active))
		progress.Passive = uint(len(tracker.Indices))/ni - progress.Active
		if level := uint(maxUint64(lindices)); level > progress.Level {
			progress.Level = level
		}

		target.Monitor(progress)

		nn := uint(len(lindices)) / ni
		indices, counts := []uint64(nil), make([]uint, nn)
		for i := uint(0); i < nn; i++ {
			newIndices := self.grid.Index(lindices[i*ni : (i+1)*ni])
			indices = append(indices, newIndices...)
			counts[i] = uint(len(newIndices)) / ni
		}

		progress.Evaluations += uint(len(indices)) / ni
		if progress.Evaluations > config.MaxEvaluations {
			break
		}

		nodes := self.grid.Compute(indices)
		values := internal.Invoke(target.Compute, nodes, ni, no, nw)
		surpluses := internal.Subtract(values, internal.Approximate(self.basis,
			surrogate.Indices, surrogate.Surpluses, nodes, ni, no, nw))

		surrogate.push(indices, surpluses)

		for _, count := range counts {
			offset := count * no
			tracker.Push(metric.Score(&Location{values[:offset], surpluses[:offset]}))
			values, surpluses = values[offset:], surpluses[offset:]
		}

		if metric.Done(tracker.Active) {
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
