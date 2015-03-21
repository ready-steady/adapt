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
	Balance(indices []uint64, find func([]uint64) bool, push func([]uint64))
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
	ni, no := target.Dimensions()
	nw := self.config.Workers

	surrogate := newSurrogate(ni, no)
	queue := newQueue(ni, &self.config)
	hash := newHash(ni)

	indices := queue.pull()
	nodes := self.grid.Compute(indices)

	na, np := uint(len(indices))/ni, uint(0)

	for k := uint(0); na > 0; k++ {
		target.Monitor(k, np, na)

		values := invoke(target.Compute, nodes, ni, no, nw)

		approximations := approximate(self.basis, surrogate.Indices,
			surrogate.Surpluses, nodes, ni, no, nw)

		surpluses := make([]float64, na*no)
		for i := uint(0); i < na*no; i++ {
			surpluses[i] = values[i] - approximations[i]
		}

		surrogate.push(indices, surpluses)

		scores := measure(self.basis, indices, ni)
		for i := uint(0); i < na; i++ {
			scores[i] = target.Refine(nodes[i*ni:(i+1)*ni],
				surpluses[i*no:(i+1)*no], scores[i])
		}

		queue.push(indices, scores)

		indices = queue.pull()
		indices = self.grid.Refine(indices)
		indices = hash.unique(indices)
		self.grid.Balance(indices, hash.find, func(index []uint64) {
			indices = append(indices, index...)
			hash.tap(index)
		})

		nodes = self.grid.Compute(indices)

		np += na
		na = uint(len(indices)) / ni
	}

	surrogate.Nodes = np
	surrogate.Level = level(surrogate.Indices, ni)

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
		volume := basis.Integrate(indices[i*ni : (i+1)*ni])
		for j := uint(0); j < no; j++ {
			value[j] += surpluses[i*no+j] * volume
		}
	}

	return value
}

func measure(basis Basis, indices []uint64, ni uint) []float64 {
	nn := uint(len(indices)) / ni

	volumes := make([]float64, nn)

	for i := uint(0); i < nn; i++ {
		volumes[i] = basis.Integrate(indices[i*ni : (i+1)*ni])
	}

	return volumes
}

func level(indices []uint64, ni uint) uint {
	nn := uint(len(indices)) / ni

	level := uint(0)
	for i := uint(0); i < nn; i++ {
		l := uint(0)
		for j := uint(0); j < ni; j++ {
			l += uint(0xFFFFFFFF & indices[i*ni+j])
		}
		if l > level {
			level = l
		}
	}

	return level
}
