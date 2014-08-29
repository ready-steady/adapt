// Package adhier provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package adhier

import (
	"fmt"
	"math"
)

const (
	bufferInitCount  = 200
	bufferGrowFactor = 2
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
	grid     Grid
	basis    Basis
	outCount uint16
	minLevel uint8
	maxLevel uint8
	absError float64
	relError float64
}

// New creates an instance of the algorithm for the given sparse grid and
// functional basis.
func New(grid Grid, basis Basis, outCount uint16) *Self {
	return &Self{
		grid:     grid,
		basis:    basis,
		outCount: outCount,
		minLevel: 1,
		maxLevel: 9,
		absError: 1e-4,
		relError: 1e-2,
	}
}

// Surrogate is the result of Construct, which represents an interpolant for a
// function.
type Surrogate struct {
	level     uint8
	inCount   uint16
	outCount  uint16
	nodeCount uint32

	levels    []uint8
	orders    []uint32
	surpluses []float64
}

func (s *Surrogate) initialize(inCount, outCount uint16) {
	ins := bufferInitCount*inCount
	outs := bufferInitCount*outCount

	s.inCount = inCount
	s.outCount = outCount
	s.nodeCount = bufferInitCount

	s.levels = make([]uint8, ins)
	s.orders = make([]uint32, ins)
	s.surpluses = make([]float64, outs)
}

func (s *Surrogate) finalize(level uint8, nodeCount uint32) {
	ins := nodeCount*uint32(s.inCount)
	outs := nodeCount*uint32(s.outCount)

	s.level = level
	s.nodeCount = nodeCount

	s.levels = s.levels[0:ins]
	s.orders = s.orders[0:ins]
	s.surpluses = s.surpluses[0:outs]
}

func (s *Surrogate) resize(nodeCount uint32) {
	if nodeCount <= s.nodeCount {
		return
	}

	if count := bufferGrowFactor * s.nodeCount; count > nodeCount {
		nodeCount = count
	}

	// New sizes
	ins := nodeCount*uint32(s.inCount)
	outs := nodeCount*uint32(s.outCount)

	levels := make([]uint8, ins)
	orders := make([]uint32, ins)
	surpluses := make([]float64, outs)

	// Old sizes
	ins = s.nodeCount*uint32(s.inCount)
	outs = s.nodeCount*uint32(s.outCount)

	copy(levels, s.levels[0:ins])
	copy(orders, s.orders[0:ins])
	copy(surpluses, s.surpluses[0:outs])

	s.nodeCount = nodeCount

	s.levels = levels
	s.orders = orders
	s.surpluses = surpluses
}

// String returns a string containing human-friendly information about the
// surrogate.
func (s *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{ inputs: %d, outputs: %d, levels: %d, nodes: %d }",
		s.inCount, s.outCount, s.level, s.nodeCount)
}

// Construct takes a function and yields a surrogate/interpolant for it, which
// can be further fed to Evaluate for the actual interpolation.
func (self *Self) Construct(target func([]float64) []float64) *Surrogate {
	inc := uint32(self.grid.Dimensionality())
	outc := uint32(self.outCount)

	surrogate := new(Surrogate)
	surrogate.initialize(uint16(inc), uint16(outc))

	// Assume level 0 has only one node, and its order is 0.
	level := uint8(0)
	newc := uint32(1)
	levels := make([]uint8, newc*inc)
	orders := make([]uint32, newc*inc)

	oldc := uint32(0)
	nodeCount := uint32(0)

	var i, j, k, l uint32

	value := make([]float64, outc)

	minValue := make([]float64, outc)
	maxValue := make([]float64, outc)

	minValue[0] = math.Inf(1)
	maxValue[0] = math.Inf(-1)
	for i = 1; i < outc; i++ {
		minValue[i] = minValue[i-1]
		maxValue[i] = maxValue[i-1]
	}

	for {
		surrogate.resize(oldc + newc)

		copy(surrogate.levels[oldc*inc:], levels)
		copy(surrogate.orders[oldc*inc:], orders)

		nodes := self.grid.ComputeNodes(levels, orders)
		values := target(nodes)

		// Compute the surpluses corresponding to the active nodes.
		for i, k = 0, oldc*outc; i < newc; i++ {
			evaluate(self.basis, surrogate, inc, outc, oldc,
				nodes[i*inc:(i+1)*inc], value)

			for j = 0; j < outc; j++ {
				surrogate.surpluses[k] = values[i*outc+j] - value[j]
				k++
			}
		}

		nodeCount += newc

		if level >= self.maxLevel {
			break
		}

		// Keep track of the maximal and minimal values of the function.
		for i, k = 0, 0; i < newc; i++ {
			for j = 0; j < outc; j++ {
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

				for j = 0; j < outc; j++ {
					absError := math.Abs(surrogate.surpluses[(oldc+i)*outc+j])

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
					l += inc
					continue
				}

				if k != l {
					// Shift everything, assuming a lot of refinements.
					copy(levels[k:], levels[l:])
					copy(orders[k:], orders[l:])
					l = k
				}

				k += inc
				l += inc
			}

			levels = levels[0:k]
			orders = orders[0:k]
		}

		levels, orders = self.grid.ComputeChildren(levels, orders)

		oldc += newc
		newc = uint32(len(levels)) / inc

		if newc == 0 {
			break
		}

		level++
	}

	surrogate.finalize(level, nodeCount)
	return surrogate
}

// Evaluate takes a surrogate produced by Construct and evaluates it at the
// given points.
func (self *Self) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	inc := uint32(surrogate.inCount)
	outc := uint32(surrogate.outCount)

	pointCount := uint32(len(points)) / inc

	values := make([]float64, pointCount*outc)

	for i := uint32(0); i < pointCount; i++ {
		evaluate(self.basis, surrogate, inc, outc, surrogate.nodeCount,
			points[i*inc:], values[i*outc:])
	}

	return values
}

func evaluate(basis Basis, surrogate *Surrogate, inc, outc, nodeCount uint32,
	point []float64, value []float64) {

	var i, j uint32 = 1, 0

	// Rewrite value in case it is dirty (not zeroed).
	weight := basis.Evaluate(surrogate.levels, surrogate.orders, point)
	for ; j < outc; j++ {
		value[j] = surrogate.surpluses[j] * weight
	}

	for ; i < nodeCount; i++ {
		weight = basis.Evaluate(surrogate.levels[i*inc:], surrogate.orders[i*inc:], point)
		for j = 0; j < outc; j++ {
			value[j] += surrogate.surpluses[i*outc+j] * weight
		}
	}
}
