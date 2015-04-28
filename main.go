// Package adapt provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package adapt

import (
	"runtime"
)

// Grid is a sparse grid in [0, 1]^n.
type Grid interface {
	Compute(indices []uint64) []float64
	Refine(indices []uint64) []uint64
	Parent(index []uint64, i uint)
	Sibling(index []uint64, i uint)
}

// Basis is a functional basis in [0, 1]^n.
type Basis interface {
	Compute(index []uint64, point []float64) float64
	Integrate(index []uint64) float64
}

// Interpolator represents a particular instantiation of the algorithm.
type Interpolator struct {
	grid   Grid
	basis  Basis
	config Config
}

// New creates an instance of the algorithm for the given configuration.
func New(grid Grid, basis Basis, config *Config) *Interpolator {
	interpolator := &Interpolator{
		grid:   grid,
		basis:  basis,
		config: *config,
	}

	config = &interpolator.config
	if config.Workers == 0 {
		config.Workers = uint(runtime.GOMAXPROCS(0))
	}
	if config.Rate == 0 {
		config.Rate = 1
	}

	return interpolator
}

// Compute constructs an interpolant for a quantity of interest.
func (self *Interpolator) Compute(target Target) *Surrogate {
	config := &self.config

	ni, no := target.Dimensions()
	nw := config.Workers

	surrogate := newSurrogate(ni, no)
	tracker := newQueue(ni, config)
	history := newHash(ni)

	na, nr, nc := uint(0), uint(0), uint(1)

	indices := make([]uint64, nc*ni)
	nodes := self.grid.Compute(indices)

	for k := uint(0); nc > 0; k++ {
		target.Monitor(k, na, nr, nc)

		values := invoke(target.Compute, nodes, ni, no, nw)

		approximations := approximate(self.basis, surrogate.Indices,
			surrogate.Surpluses, nodes, ni, no, nw)

		surpluses := make([]float64, nc*no)
		for i := uint(0); i < nc*no; i++ {
			surpluses[i] = values[i] - approximations[i]
		}

		scores := measure(self.basis, indices, ni)
		for i := uint(0); i < nc; i++ {
			scores[i] = target.Score(nodes[i*ni:(i+1)*ni],
				surpluses[i*no:(i+1)*no], scores[i])
		}

		accept := tracker.push(indices, scores)

		nn := surrogate.push(indices, surpluses, accept)
		na, nr = na+nn, nr+nc-nn

		indices = tracker.pull()
		indices = self.grid.Refine(indices)
		indices = history.unseen(indices)

		if config.Balance {
			// NOTE: Rejection might play against balancing.
			indices = append(indices, balance(self.grid, history, indices)...)
		}

		nodes = self.grid.Compute(indices)

		nc = uint(len(indices)) / ni
	}

	surrogate.Level = tracker.lnow

	return surrogate
}

// Evaluate computes the values of a surrogate at a set of points.
func (self *Interpolator) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	return approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
}

// Integrate computes the integral of a surrogate over [0, 1]^n.
func (self *Interpolator) Integrate(surrogate *Surrogate) []float64 {
	return integrate(self.basis, surrogate.Indices, surrogate.Surpluses,
		surrogate.Inputs, surrogate.Outputs)
}
