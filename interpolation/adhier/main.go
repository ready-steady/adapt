// Package adhier provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package adhier

import (
	"errors"
	"math"
)

// Grid is the interface that an sparse grid should satisfy in order to be used
// in the algorithm.
type Grid interface {
	Dimensions() uint16
	ComputeNodes(indices []uint64) []float64
	ComputeChildren(indices []uint64) []uint64
}

// Basis is the interface that a functional basis should satisfy in order to be
// used in the algorithm.
type Basis interface {
	Outputs() uint16
	EvaluateComposite(indices []uint64, weights, point, result []float64)
}

// Interpolator represents a particular instantiation of the algorithm.
type Interpolator struct {
	grid   Grid
	basis  Basis
	config Config

	ic uint32
	oc uint32
}

// New creates an instance of the algorithm for the given configuration.
func New(grid Grid, basis Basis, config Config) (*Interpolator, error) {
	if config.AbsError <= 0 {
		return nil, errors.New("the absolute error is invalid")
	}
	if config.RelError <= 0 {
		return nil, errors.New("the relative error is invalid")
	}

	interpolator := &Interpolator{
		grid:   grid,
		basis:  basis,
		config: config,

		ic: uint32(grid.Dimensions()),
		oc: uint32(basis.Outputs()),
	}

	return interpolator, nil
}

// Compute takes a function and yields a surrogate for it, which can be further
// fed to Evaluate for the actual interpolation.
func (self *Interpolator) Compute(target func([]float64, []uint64) []float64) *Surrogate {
	ic, oc := self.ic, self.oc

	basis := self.basis

	surrogate := new(Surrogate)
	surrogate.initialize(ic, oc)

	// Level 0 is assumed to have only one node, and the order of that node is
	// assumed to be zero.
	level := uint8(0)

	ac := uint32(1) // active
	pc := uint32(0) // passive
	nc := uint32(0) // total

	indices := make([]uint64, ac*ic)

	var i, j, k, l uint32

	value := make([]float64, oc)

	min := make([]float64, oc)
	max := make([]float64, oc)

	min[0], max[0] = math.Inf(1), math.Inf(-1)
	for i = 1; i < oc; i++ {
		min[i], max[i] = min[0], max[0]
	}

	for {
		surrogate.resize(pc + ac)
		copy(surrogate.indices[pc*ic:], indices)

		nodes := self.grid.ComputeNodes(indices)
		values := target(nodes, indices)

		// Compute the surpluses corresponding to the active nodes.
		passiveIndices := surrogate.indices[:pc*ic]
		passiveSurpluses := surrogate.surpluses[:pc*oc]
		for i, k = 0, pc*oc; i < ac; i++ {
			basis.EvaluateComposite(passiveIndices, passiveSurpluses,
				nodes[i*ic:(i+1)*ic], value)
			for j = 0; j < oc; j++ {
				surrogate.surpluses[k] = values[i*oc+j] - value[j]
				k++
			}
		}

		nc += ac

		if level >= self.config.MaxLevel || nc >= self.config.MaxNodes {
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

		if level < self.config.MinLevel {
			goto next
		}

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
				copy(indices[k:], indices[l:])
				l = k
			}

			k += ic
			l += ic
		}

		indices = indices[:k]

	next:
		indices = self.grid.ComputeChildren(indices)

		pc += ac
		ac = uint32(len(indices)) / ic

		if δ := int32(nc+ac) - int32(self.config.MaxNodes); δ > 0 {
			ac -= uint32(δ)
			indices = indices[:ac*ic]
		}
		if ac == 0 {
			break
		}

		level++
	}

	surrogate.finalize(level, nc)
	return surrogate
}

// Evaluate takes a surrogate produced by Compute and evaluates it at the given
// points.
func (self *Interpolator) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	ic, oc, nc := surrogate.ic, surrogate.oc, surrogate.nc
	pc := uint32(len(points)) / ic

	basis := self.basis
	indices := surrogate.indices[:nc*ic]
	surpluses := surrogate.surpluses[:nc*oc]

	values := make([]float64, pc*oc)
	for i := uint32(0); i < pc; i++ {
		basis.EvaluateComposite(indices, surpluses,
			points[i*ic:(i+1)*ic], values[i*oc:(i+1)*oc])
	}

	return values
}
