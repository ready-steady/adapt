// Package hybrid provides an algorithm for hierarchical interpolation with
// hybrid adaptation.
package hybrid

import (
	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/adapt/algorithm/internal"
)

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
	basis  external.Basis
	config Config
}

// New creates an interpolator.
func New(grid Grid, basis external.Basis, config *Config) *Interpolator {
	return &Interpolator{
		grid:   grid,
		basis:  basis,
		config: *config,
	}
}

// Compute constructs an interpolant for a function.
func (self *Interpolator) Compute(target Target) *external.Surrogate {
	config := &self.config

	ni, no := target.Dimensions()
	nw := config.Workers

	surrogate := external.NewSurrogate(ni, no)
	tracker := newTracker(ni, config)

	progress := &Progress{}
	for {
		lindices := tracker.pull()

		nn := uint(len(lindices)) / ni
		indices, counts := []uint64(nil), make([]uint, nn)
		for i := uint(0); i < nn; i++ {
			newIndices := self.grid.Index(lindices[i*ni : (i+1)*ni])
			indices = append(indices, newIndices...)
			counts[i] = uint(len(newIndices)) / ni
		}

		progress.Level = internal.MaxUint(progress.Level, uint(internal.MaxUint64s(lindices)))
		progress.Active, progress.Passive = tracker.CountActive(), tracker.CountPassive()
		progress.Performed += progress.Requested
		progress.Requested = uint(len(indices)) / ni

		if !target.Before(progress) {
			break
		}

		nodes := self.grid.Compute(indices)
		values := internal.Invoke(target.Compute, nodes, ni, no, nw)
		surpluses := internal.Subtract(values, internal.Approximate(self.basis,
			surrogate.Indices, surrogate.Surpluses, nodes, ni, no, nw))

		surrogate.Push(indices, surpluses)

		for _, count := range counts {
			offset := count * no
			tracker.push(target.Score(&Location{values[:offset], surpluses[:offset]}))
			values, surpluses = values[offset:], surpluses[offset:]
		}

		if !target.After(tracker.Active) {
			break
		}
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *external.Surrogate, points []float64) []float64 {
	return internal.Approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
}
