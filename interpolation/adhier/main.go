// Package adhier provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package adhier

import (
	"runtime"
	"sync"
)

// Grid is a sparse grid in [0, 1]^n.
type Grid interface {
	Compute(indices []uint64) []float64
	Refine(indices []uint64) []uint64
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

	return interpolator
}

// Compute constructs an interpolant for a quantity of interest.
func (self *Interpolator) Compute(target Target) *Surrogate {
	config := &self.config

	ni, no := target.Dimensions()

	surrogate := newSurrogate(ni, no)

	hash := newHash(ni)
	queue := newQueue(ni, config.MinLevel, config.MaxLevel, config.Rate)

	// The zeroth level is assumed to have only one node, and the order of that
	// node is assumed to be zero.
	na, np := uint(1), uint(0)

	indices := make([]uint64, na*ni)

	for k := uint(0); ; k++ {
		target.Monitor(k, np, na)

		surrogate.resize(np + na)
		copy(surrogate.Indices[np*ni:], indices)

		nodes := self.grid.Compute(indices)

		values := invoke(target.Compute, nodes, ni, no, config.Workers)
		approximations := approximate(self.basis, surrogate.Indices[:np*ni],
			surrogate.Surpluses[:np*no], nodes, ni, no, config.Workers)

		surpluses := surrogate.Surpluses[np*no : (np+na)*no]
		for i := uint(0); i < na*no; i++ {
			surpluses[i] = values[i] - approximations[i]
		}

		if np+na >= config.MaxNodes {
			break
		}

		scores := make([]float64, na)
		volumes := measure(self.basis, indices, surpluses, ni, no)
		for i := uint(0); i < na; i++ {
			scores[i] = target.Refine(nodes[i*ni:(i+1)*ni], surpluses[i*no:(i+1)*no],
				volumes[i*no:(i+1)*no])
		}

		queue.push(indices, scores)
		indices = queue.pull()
		if len(indices) == 0 {
			break
		}

		indices = hash.unique(self.grid.Refine(indices))

		np += na
		na = uint(len(indices)) / ni

		// Trim if there are excessive nodes.
		if np+na > config.MaxNodes {
			na = config.MaxNodes - np
			indices = indices[:na*ni]
		}
	}

	surrogate.finalize(np + na)
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

func approximate(basis Basis, indices []uint64, surpluses, points []float64,
	ni, no, nw uint) []float64 {

	nn, np := uint(len(indices))/ni, uint(len(points))/ni

	values := make([]float64, np*no)

	jobs := make(chan uint, np)
	group := sync.WaitGroup{}
	group.Add(int(np))

	for i := uint(0); i < nw; i++ {
		go func() {
			for j := range jobs {
				point := points[j*ni : (j+1)*ni]
				value := values[j*no : (j+1)*no]

				for k := uint(0); k < nn; k++ {
					weight := basis.Compute(indices[k*ni:(k+1)*ni], point)
					if weight == 0 {
						continue
					}
					for l := uint(0); l < no; l++ {
						value[l] += weight * surpluses[k*no+l]
					}
				}

				group.Done()
			}
		}()
	}

	for i := uint(0); i < np; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
}

func invoke(compute func([]float64, []float64), nodes []float64, ni, no, nw uint) []float64 {
	nn := uint(len(nodes)) / ni

	values := make([]float64, nn*no)

	jobs := make(chan uint, nn)
	group := sync.WaitGroup{}
	group.Add(int(nn))

	for i := uint(0); i < nw; i++ {
		go func() {
			for j := range jobs {
				compute(nodes[j*ni:(j+1)*ni], values[j*no:(j+1)*no])
				group.Done()
			}
		}()
	}

	for i := uint(0); i < nn; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
}

func integrate(basis Basis, indices []uint64, surpluses []float64, ni, no uint) []float64 {
	nn := uint(len(indices)) / ni

	value := make([]float64, no)

	for i := uint(0); i < nn; i++ {
		weight := basis.Integrate(indices[i*ni : (i+1)*ni])
		for j := uint(0); j < no; j++ {
			value[j] += weight * surpluses[i*no+j]
		}
	}

	return value
}

func measure(basis Basis, indices []uint64, surpluses []float64, ni, no uint) []float64 {
	nn := uint(len(indices)) / ni

	volumes := make([]float64, nn*no)

	for i := uint(0); i < nn; i++ {
		weight := basis.Integrate(indices[i*ni : (i+1)*ni])
		for j := uint(0); j < no; j++ {
			volumes[i*no+j] = weight * surpluses[i*no+j]
		}
	}

	return volumes
}
