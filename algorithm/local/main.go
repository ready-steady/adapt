// Package local provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package adapt

import (
	"runtime"
)

const (
	LEVEL_MASK = 0x3F
	LEVEL_SIZE = 6
)

// Basis is a functional basis.
type Basis interface {
	// Compute evaluates the value of a basis function at a point.
	Compute([]uint64, []float64) float64

	// Integrate computes the integral of a basis function over the whole
	// domain.
	Integrate([]uint64) float64
}

// Grid is a sparse grid.
type Grid interface {
	// Compute returns the nodes corresponding to the given indices.
	Compute([]uint64) []float64

	// Children returns the child indices corresponding to a set of parent
	// indices.
	Children([]uint64) []uint64
}

// Interpolator is an instance of the algorithm.
type Interpolator struct {
	config Config
	basis  Basis
	grid   Grid
}

// Progress contains information about the interpolation process.
type Progress struct {
	Level    uint      // Reached level
	Active   uint      // Number of active nodes
	Passive  uint      // Number of passive nodes
	Refined  uint      // Number of refined nodes
	Integral []float64 // Integral over the whole domain
}

// New creates a new interpolator.
func New(grid Grid, basis Basis, config *Config) *Interpolator {
	interpolator := &Interpolator{
		config: *config,
		basis:  basis,
		grid:   grid,
	}
	config = &interpolator.config
	if config.Workers == 0 {
		config.Workers = uint(runtime.GOMAXPROCS(0))
	}
	return interpolator
}

// Compute constructs an interpolant for a function.
func (self *Interpolator) Compute(target Target) *Surrogate {
	config := &self.config

	ni, no := target.Dimensions()
	nw := config.Workers

	surrogate := newSurrogate(ni, no)
	queue := newQueue(ni, config)
	hash := newHash(ni)

	indices := make([]uint64, 1*ni)
	progress := Progress{Active: 1, Integral: make([]float64, no)}
	for {
		target.Monitor(&progress)

		nodes := self.grid.Compute(indices)
		values := invoke(target.Compute, nodes, ni, no, nw)
		surpluses := subtract(values, approximate(self.basis, surrogate.Indices,
			surrogate.Surpluses, nodes, ni, no, nw))

		surrogate.push(indices, surpluses)
		cumulate(self.basis, indices, surpluses, ni, no, progress.Integral)

		scores := assess(self.basis, target, &progress, indices, nodes, surpluses, ni, no)
		indices = queue.filter(indices, scores)

		progress.Refined += uint(len(indices)) / ni

		indices = hash.filter(self.grid.Children(indices))

		progress.Passive += progress.Active
		progress.Active = uint(len(indices)) / ni

		if progress.Active == 0 || progress.Active+progress.Passive > config.MaxEvaluations {
			break
		}

		progress.Level++
	}

	surrogate.Level = progress.Level
	return surrogate
}

// Evaluate computes the values of an interpolant at a set of points.
func (self *Interpolator) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	return approximate(self.basis, surrogate.Indices, surrogate.Surpluses, points,
		surrogate.Inputs, surrogate.Outputs, self.config.Workers)
}

// Integrate computes the integral of an interpolant over the whole domain.
func (self *Interpolator) Integrate(surrogate *Surrogate) []float64 {
	ni, no := surrogate.Inputs, surrogate.Outputs
	integral := make([]float64, no)
	cumulate(self.basis, surrogate.Indices, surrogate.Surpluses, ni, no, integral)
	return integral
}
