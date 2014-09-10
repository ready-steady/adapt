// Package adhier provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package adhier

import (
	"math"
)

// Grid is the interface that an sparse grid should satisfy in order to be used
// in the algorithm.
type Grid interface {
	Dimensionality() uint16
	ComputeNodes(levels []uint8, orders []uint32) []float64
	ComputeChildren(levels []uint8, orders []uint32) ([]uint8, []uint32)
}

// Basis is the interface that a functional basis should satisfy in order to be
// used in the algorithm.
type Basis interface {
	Evaluate(levels []uint8, orders []uint32, point []float64) float64
}

// Self represents a particular instantiation of the algorithm.
type Self struct {
	grid  Grid
	basis Basis

	oc uint16

	minLevel uint8
	maxLevel uint8

	absError float64
	relError float64
}

// New creates an instance of the algorithm for the given sparse grid and
// functional basis.
func New(grid Grid, basis Basis, outputs uint16) *Self {
	return &Self{
		grid:  grid,
		basis: basis,

		oc: outputs,

		minLevel: 1,
		maxLevel: 9,

		absError: 1e-4,
		relError: 1e-2,
	}
}

// Compute takes a function and yields a surrogate/interpolant for it, which
// can be further fed to Evaluate for the actual interpolation.
func (self *Self) Compute(target func([]float64) []float64) *Surrogate {
	ic := uint32(self.grid.Dimensionality())
	oc := uint32(self.oc)

	surrogate := new(Surrogate)
	surrogate.initialize(uint16(ic), uint16(oc))

	// Assume level 0 has only one node, and its order is 0.
	level := uint8(0)
	newc := uint32(1)
	levels := make([]uint8, newc*ic)
	orders := make([]uint32, newc*ic)

	oldc := uint32(0)
	nc := uint32(0)

	var i, j, k, l uint32

	value := make([]float64, oc)

	minValue := make([]float64, oc)
	maxValue := make([]float64, oc)

	minValue[0] = math.Inf(1)
	maxValue[0] = math.Inf(-1)
	for i = 1; i < oc; i++ {
		minValue[i] = minValue[i-1]
		maxValue[i] = maxValue[i-1]
	}

	for {
		surrogate.resize(oldc + newc)

		copy(surrogate.levels[oldc*ic:], levels)
		copy(surrogate.orders[oldc*ic:], orders)

		nodes := self.grid.ComputeNodes(levels, orders)
		values := target(nodes)

		// Compute the surpluses corresponding to the active nodes.
		for i, k = 0, oldc*oc; i < newc; i++ {
			evaluate(self.basis, surrogate, ic, oc, oldc,
				nodes[i*ic:(i+1)*ic], value)

			for j = 0; j < oc; j++ {
				surrogate.surpluses[k] = values[i*oc+j] - value[j]
				k++
			}
		}

		nc += newc

		if level >= self.maxLevel {
			break
		}

		// Keep track of the maximal and minimal values of the function.
		for i, k = 0, 0; i < newc; i++ {
			for j = 0; j < oc; j++ {
				if values[k] < minValue[j] {
					minValue[j] = values[k]
				}
				if values[k] > maxValue[j] {
					maxValue[j] = values[k]
				}
				k++
			}
		}

		if level >= self.minLevel {
			k, l = 0, 0

			for i = 0; i < newc; i++ {
				refine := false

				for j = 0; j < oc; j++ {
					absError := math.Abs(surrogate.surpluses[(oldc+i)*oc+j])

					if absError > self.absError {
						refine = true
						break
					}

					relError := absError / (maxValue[j] - minValue[j])

					if relError > self.relError {
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
					copy(levels[k:], levels[l:])
					copy(orders[k:], orders[l:])
					l = k
				}

				k += ic
				l += ic
			}

			levels = levels[0:k]
			orders = orders[0:k]
		}

		levels, orders = self.grid.ComputeChildren(levels, orders)

		oldc += newc
		newc = uint32(len(levels)) / ic

		if newc == 0 {
			break
		}

		level++
	}

	surrogate.finalize(level, nc)
	return surrogate
}

// Evaluate takes a surrogate produced by Compute and evaluates it at the
// given points.
func (self *Self) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	ic := uint32(surrogate.ic)
	oc := uint32(surrogate.oc)

	pointCount := uint32(len(points)) / ic

	values := make([]float64, pointCount*oc)

	for i := uint32(0); i < pointCount; i++ {
		evaluate(self.basis, surrogate, ic, oc, surrogate.nc,
			points[i*ic:], values[i*oc:])
	}

	return values
}

func evaluate(basis Basis, surrogate *Surrogate, ic, oc, nc uint32,
	point []float64, value []float64) {

	var i, j uint32 = 1, 0

	// Rewrite value in case it is dirty (not zeroed).
	weight := basis.Evaluate(surrogate.levels, surrogate.orders, point)
	for ; j < oc; j++ {
		value[j] = surrogate.surpluses[j] * weight
	}

	for ; i < nc; i++ {
		weight = basis.Evaluate(surrogate.levels[i*ic:], surrogate.orders[i*ic:], point)
		for j = 0; j < oc; j++ {
			value[j] += surrogate.surpluses[i*oc+j] * weight
		}
	}
}
