// Package adapt provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package adapt

import (
	"runtime"
)

// Grid is a sparse grid in [0, 1]^n or (0, 1)^n.
type Grid interface {
	// Compute returns the nodes corresponding to the given indices.
	Compute([]uint64) []float64

	// Refine returns the child indices corresponding to a set of parent
	// indices.
	Refine([]uint64) []uint64

	// Parent transforms an index into its parent index in the ith dimension.
	Parent([]uint64, uint)

	// Sibling transforms an index into its sibling index in the ith dimension.
	Sibling([]uint64, uint)
}

// Basis is a functional basis in [0, 1]^n or (0, 1)^n.
type Basis interface {
	// Compute evaluates the value of a basis function at a point.
	Compute([]uint64, []float64) float64

	// Integrate computes the integral of a basis function over the whole
	// domain.
	Integrate([]uint64) float64
}

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	grid   Grid
	basis  Basis
	config Config
}

// Progress contains information about the interpolation process.
type Progress struct {
	Level     uint      // Reached level
	Iteration uint      // Iteration number
	Accepted  uint      // Number of accepted nodes
	Rejected  uint      // Number of rejected nodes
	Current   uint      // Number of current nodes
	Integral  []float64 // Integral over the whole domain
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

	integral := make([]float64, no)

	for k := uint(0); nc > 0; k++ {
		progress := Progress{
			Level:     tracker.lnow,
			Iteration: k,
			Accepted:  na,
			Rejected:  nr,
			Current:   nc,
			Integral:  integral,
		}

		target.Monitor(&progress)

		surpluses := subtract(
			invoke(target.Compute, nodes, ni, no, nw),
			approximate(self.basis, surrogate.Indices, surrogate.Surpluses, nodes, ni, no, nw),
		)

		location := Location{}
		scores := measure(self.basis, indices, ni)
		for i := uint(0); i < nc; i++ {
			location = Location{
				Node:    nodes[i*ni : (i+1)*ni],
				Surplus: surpluses[i*no : (i+1)*no],
				Volume:  scores[i],
			}
			scores[i] = target.Score(&location, &progress)
		}

		indices, surpluses, scores = compact(indices, surpluses, scores, ni, no, nc)

		nn := uint(len(scores))
		na, nr = na+nn, nr+nc-nn

		tracker.push(indices, scores)
		surrogate.push(indices, surpluses)
		surrogate.step(tracker.lnow, nn, nc-nn)

		cumulate(self.basis, indices, surpluses, ni, no, nn, integral)

		indices = history.unseen(self.grid.Refine(tracker.pull()))
		if config.Balance {
			indices = append(indices, balance(self.grid, history, indices)...)
		}

		nodes = self.grid.Compute(indices)

		nc = uint(len(indices)) / ni
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

	integral := make([]float64, no)
	cumulate(self.basis, surrogate.Indices, surrogate.Surpluses, ni, no, nn, integral)

	return integral
}
