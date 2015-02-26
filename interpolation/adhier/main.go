// Package adhier provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package adhier

import (
	"errors"
	"math"
	"runtime"
	"sync"
)

// Grid is the interface that an sparse grid should satisfy in order to be used
// in the algorithm.
type Grid interface {
	ComputeNodes(indices []uint64) []float64
	ComputeChildren(indices []uint64) []uint64
}

// Basis is the interface that a functional basis should satisfy in order to be
// used in the algorithm.
type Basis interface {
	Evaluate(index []uint64, point []float64) float64
}

// Interpolator represents a particular instantiation of the algorithm.
type Interpolator struct {
	grid   Grid
	basis  Basis
	config *Config

	ic uint
	oc uint

	wc uint
}

// New creates an instance of the algorithm for the given configuration.
func New(grid Grid, basis Basis, config *Config) (*Interpolator, error) {
	if config.Inputs == 0 {
		return nil, errors.New("the number of inputs should be positive")
	}
	if config.Outputs == 0 {
		return nil, errors.New("the number of outputs should be positive")
	}
	if config.AbsError < 0 {
		return nil, errors.New("the absolute error tolerance should be nonnegative")
	}
	if config.RelError < 0 {
		return nil, errors.New("the relative error tolerance should be nonnegative")
	}

	wc := config.Workers
	if wc == 0 {
		wc = uint(runtime.GOMAXPROCS(0))
	}

	interpolator := &Interpolator{
		grid:   grid,
		basis:  basis,
		config: config,

		ic: config.Inputs,
		oc: config.Outputs,

		wc: wc,
	}

	return interpolator, nil
}

// Compute takes a target function and produces an interpolant for it. The
// interpolant can then be fed to Evaluate for approximating the target function
// at arbitrary points.
//
// The second argument of Compute is an optional function that can be used for
// monitoring the progress of interpolation. The progress function is called
// once for each level before evaluating the target function at the nodes of
// that level. The signature of the progress function is func(uint32, uint,
// uint) where the arguments are the current level, number of active nodes,
// and total number of nodes, respectively.
func (self *Interpolator) Compute(target func([]float64, []float64, []uint64),
	arguments ...interface{}) *Surrogate {

	var progress func(uint32, uint, uint)
	if len(arguments) > 0 {
		progress = arguments[0].(func(uint32, uint, uint))
	}

	ic, oc := self.ic, self.oc
	config := self.config

	surrogate := new(Surrogate)
	surrogate.initialize(ic, oc)

	// Level 0 is assumed to have only one node, and the order of that node is
	// assumed to be zero.
	level := uint32(0)

	ac := uint(1) // active
	pc := uint(0) // passive

	indices := make([]uint64, ac*ic)

	var i, j, k, l uint
	var nodes, values, approximations []float64

	min := make([]float64, oc)
	max := make([]float64, oc)

	min[0], max[0] = math.Inf(1), math.Inf(-1)
	for i = 1; i < oc; i++ {
		min[i], max[i] = min[0], max[0]
	}

	for {
		if progress != nil {
			progress(level, ac, pc+ac)
		}

		surrogate.resize(pc + ac)
		copy(surrogate.Indices[pc*ic:], indices)

		nodes = self.grid.ComputeNodes(indices)

		// NOTE: Assuming that target might have some logic based on the indices
		// passed to it (for instance, caching), the indices variable should not
		// be used here as it gets modified later on.
		values = self.invoke(target, nodes, surrogate.Indices[pc*ic:(pc+ac)*ic])

		// Compute the surpluses corresponding to the active nodes.
		if level == 0 {
			// The surrogate does not have any nodes yet.
			copy(surrogate.Surpluses, values)
			goto refineLevel
		}

		approximations = self.approximate(surrogate.Indices[:pc*ic],
			surrogate.Surpluses[:pc*oc], nodes)
		for i, k = 0, pc*oc; i < ac; i++ {
			for j = 0; j < oc; j++ {
				surrogate.Surpluses[k] = values[i*oc+j] - approximations[i*oc+j]
				k++
			}
		}

	refineLevel:
		if level >= config.MaxLevel || (pc+ac) >= config.MaxNodes {
			break
		}

		// Keep track of the maximal and minimal values of the function.
		for i, k = 0, 0; i < ac; i++ {
			for j = 0; j < oc; j++ {
				if values[k] < min[j] {
					min[j] = values[k]
				}
				if values[k] > max[j] {
					max[j] = values[k]
				}
				k++
			}
		}

		if level < config.MinLevel {
			goto updateIndices
		}

		k, l = 0, 0

		for i = 0; i < ac; i++ {
			refine := false

			for j = 0; j < oc; j++ {
				absError := surrogate.Surpluses[(pc+i)*oc+j]
				if absError < 0 {
					absError = -absError
				}

				if absError > config.AbsError {
					refine = true
					break
				}

				relError := absError / (max[j] - min[j])

				if relError > config.RelError {
					refine = true
					break
				}
			}

			if !refine {
				l += ic
				continue
			}

			if k != l {
				// Shift everything, assuming a lot of refinements.
				copy(indices[k:], indices[l:])
				l = k
			}

			k += ic
			l += ic
		}

		indices = indices[:k]

	updateIndices:
		indices = self.grid.ComputeChildren(indices)

		pc += ac
		ac = uint(len(indices)) / ic

		if Δ := int32(pc+ac) - int32(config.MaxNodes); Δ > 0 {
			ac -= uint(Δ)
			indices = indices[:ac*ic]
		}
		if ac == 0 {
			break
		}

		level++
	}

	surrogate.finalize(level, pc+ac)
	return surrogate
}

// Evaluate takes a surrogate produced by Compute and evaluates it at a set of
// points.
func (self *Interpolator) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	return self.approximate(surrogate.Indices, surrogate.Surpluses, points)
}

func (self *Interpolator) approximate(indices []uint64, surpluses, points []float64) []float64 {
	ic, oc, wc := self.ic, self.oc, self.wc
	nc := uint(len(indices)) / ic
	pc := uint(len(points)) / ic

	basis := self.basis

	values := make([]float64, pc*oc)

	jobs := make(chan uint, pc)
	group := sync.WaitGroup{}
	group.Add(int(pc))

	for i := uint(0); i < wc; i++ {
		go func() {
			for j := range jobs {
				point := points[j*ic : (j+1)*ic]
				value := values[j*oc : (j+1)*oc]

				for k := uint(0); k < nc; k++ {
					weight := basis.Evaluate(indices[k*ic:(k+1)*ic], point)
					if weight == 0 {
						continue
					}
					for l := uint(0); l < oc; l++ {
						value[l] += surpluses[k*oc+l] * weight
					}
				}

				group.Done()
			}
		}()
	}

	for i := uint(0); i < pc; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
}

func (self *Interpolator) invoke(target func([]float64, []float64, []uint64),
	nodes []float64, indices []uint64) []float64 {

	ic, oc, wc := self.ic, self.oc, self.wc
	nc := uint(len(nodes)) / ic

	values := make([]float64, nc*oc)

	jobs := make(chan uint, nc)
	group := sync.WaitGroup{}
	group.Add(int(nc))

	for i := uint(0); i < wc; i++ {
		go func() {
			for j := range jobs {
				target(nodes[j*ic:(j+1)*ic], values[j*oc:(j+1)*oc], indices[j*ic:(j+1)*ic])
				group.Done()
			}
		}()
	}

	for i := uint(0); i < nc; i++ {
		jobs <- i
	}

	group.Wait()
	close(jobs)

	return values
}
