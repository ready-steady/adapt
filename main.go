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

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	grid   Grid
	basis  Basis
	config Config
}

// New creates a new interpolator.
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

// Compute constructs an interpolant for a function.
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

	integral, compensation := make([]float64, no), make([]float64, no)

	for k := uint(0); nc > 0; k++ {
		global := Global{
			Integral: integral,
		}

		target.Monitor(k, na, nr, nc)

		surpluses := subtract(
			invoke(target.Compute, nodes, ni, no, nw),
			approximate(self.basis, surrogate.Indices,
				surrogate.Surpluses, nodes, ni, no, nw),
		)

		scores := measure(self.basis, indices, ni)
		for i := uint(0); i < nc; i++ {
			local := Local{
				Node:    nodes[i*ni : (i+1)*ni],
				Surplus: surpluses[i*no : (i+1)*no],
				Volume:  scores[i],
			}
			scores[i] = target.Score(local, global)
		}

		indices, surpluses, scores = compact(indices, surpluses, scores, ni, no, nc)

		nn := uint(len(scores))

		tracker.push(indices, scores)
		surrogate.push(indices, surpluses)
		surrogate.step(tracker.lnow, nn, nc-nn)

		cumulate(self.basis, indices, surpluses, ni, no, nn, integral, compensation)

		indices = history.unseen(self.grid.Refine(tracker.pull()))
		if config.Balance {
			indices = append(indices, balance(self.grid, history, indices)...)
		}

		nodes = self.grid.Compute(indices)

		nn = uint(len(indices)) / ni
		na, nr, nc = na+nn, nr+nc-nn, nn
	}

	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	return approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
}

// Integrate computes the integral of an interpolant over [0, 1]^n.
func (self *Interpolator) Integrate(surrogate *Surrogate) []float64 {
	ni, no, nn := surrogate.Inputs, surrogate.Outputs, surrogate.Nodes

	integral, compensation := make([]float64, no), make([]float64, no)
	cumulate(self.basis, surrogate.Indices, surrogate.Surpluses,
		ni, no, nn, integral, compensation)

	return integral
}
