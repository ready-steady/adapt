// Package adhier provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package adhier

import (
	"math"
)

// Grid is the interface that an sparse grid should satisfy in order to be used
// in the algorithm.
type Grid interface {
	Dimensions() uint16
	ComputeNodes(index []uint64) []float64
	ComputeChildren(index []uint64) []uint64
}

// Basis is the interface that a functional basis should satisfy in order to be
// used in the algorithm.
type Basis interface {
	Evaluate(index []uint64, point []float64) float64
}

// Self represents a particular instantiation of the algorithm.
type Self struct {
	grid   Grid
	basis  Basis
	config Config

	ic uint16
	oc uint16
}

// New creates an instance of the algorithm for the given configuration.
func New(grid Grid, basis Basis, config Config, outputs uint16) *Self {
	return &Self{
		grid:   grid,
		basis:  basis,
		config: config,

		ic: grid.Dimensions(),
		oc: outputs,
	}
}

// Compute takes a function and yields a surrogate for it, which can be further
// fed to Evaluate for the actual interpolation.
func (self *Self) Compute(target func([]float64) []float64) *Surrogate {
	ic := uint32(self.ic)
	oc := uint32(self.oc)

	surrogate := new(Surrogate)
	surrogate.initialize(uint16(ic), uint16(oc))

	// Level 0 is assumed to have only one node, and the order of that node is
	// assumed to be zero.
	level := uint8(0)

	ac := uint32(1) // active
	pc := uint32(0) // passive
	nc := uint32(0) // total

	index := make([]uint64, ac*ic)

	var i, j, k, l uint32

	value := make([]float64, oc)

	min := make([]float64, oc)
	max := make([]float64, oc)

	min[0] = math.Inf(1)
	max[0] = math.Inf(-1)
	for i = 1; i < oc; i++ {
		min[i] = min[i-1]
		max[i] = max[i-1]
	}

	for {
		surrogate.resize(pc + ac)

		copy(surrogate.index[pc*ic:], index)

		nodes := self.grid.ComputeNodes(index)
		values := target(nodes)

		// Compute the surpluses corresponding to the active nodes.
		for i, k = 0, pc*oc; i < ac; i++ {
			evaluate(self.basis, surrogate, ic, oc, pc,
				nodes[i*ic:(i+1)*ic], value)

			for j = 0; j < oc; j++ {
				surrogate.surpluses[k] = values[i*oc+j] - value[j]
				k++
			}
		}

		nc += ac

		if level >= self.config.MaxLevel {
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

		if level >= self.config.MinLevel {
			k, l = 0, 0

			for i = 0; i < ac; i++ {
				refine := false

				for j = 0; j < oc; j++ {
					absError := math.Abs(surrogate.surpluses[(pc+i)*oc+j])

					if absError > self.config.AbsError {
						refine = true
						break
					}

					relError := absError / (max[j] - min[j])

					if relError > self.config.RelError {
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
					copy(index[k:], index[l:])
					l = k
				}

				k += ic
				l += ic
			}

			index = index[0:k]
		}

		index = self.grid.ComputeChildren(index)

		pc += ac
		ac = uint32(len(index)) / ic

		if ac == 0 {
			break
		}

		level++
	}

	surrogate.finalize(level, nc)
	return surrogate
}

// Evaluate takes a surrogate produced by Compute and evaluates it at the
// given points.
func (self *Self) Evaluate(s *Surrogate, points []float64) []float64 {
	ic := uint32(s.ic)
	oc := uint32(s.oc)
	pc := uint32(len(points)) / ic

	values := make([]float64, pc*oc)
	for i := uint32(0); i < pc; i++ {
		evaluate(self.basis, s, ic, oc, s.nc, points[i*ic:], values[i*oc:])
	}

	return values
}

func evaluate(b Basis, s *Surrogate, ic, oc, nc uint32, point []float64, value []float64) {
	var i, j uint32 = 1, 0

	// Treat the first separately in case value is not zeroed.
	weight := b.Evaluate(s.index, point)
	for ; j < oc; j++ {
		value[j] = s.surpluses[j] * weight
	}

	for ; i < nc; i++ {
		weight = b.Evaluate(s.index[i*ic:], point)
		for j = 0; j < oc; j++ {
			value[j] += s.surpluses[i*oc+j] * weight
		}
	}
}
