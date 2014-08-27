// Package hadapt provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package hadapt

import (
	"fmt"
	"math"
)

const (
	bufferInitSize   = 200
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
	Evaluate(point []float64, levels []uint8, orders []uint32) float64
}

// Self represents a particular instantiation of the algorithm.
type Self struct {
	grid         Grid
	basis        Basis
	minLevel     uint8
	maxLevel     uint8
	absTolerance float64
	relTolerance float64
}

// New creates an instance of the algorithm for the given sparse grid and
// functional basis.
func New(grid Grid, basis Basis) *Self {
	return &Self{
		grid:         grid,
		basis:        basis,
		minLevel:     1,
		maxLevel:     9,
		absTolerance: 1e-4,
		relTolerance: 1e-2,
	}
}

// Surrogate is the result of Construct, which represents an interpolant for a
// function.
type Surrogate struct {
	level     uint8
	inCount   uint16
	nodeCount uint32

	levels    []uint8
	orders    []uint32
	surpluses []float64
}

func (s *Surrogate) initialize(inCount uint16) {
	s.inCount = inCount
	s.nodeCount = bufferInitSize

	s.levels = make([]uint8, bufferInitSize*inCount)
	s.orders = make([]uint32, bufferInitSize*inCount)
	s.surpluses = make([]float64, bufferInitSize)
}

func (s *Surrogate) finalize(level uint8, nodeCount uint32) {
	s.level = level
	s.nodeCount = nodeCount

	s.levels = s.levels[0 : nodeCount*uint32(s.inCount)]
	s.orders = s.orders[0 : nodeCount*uint32(s.inCount)]
	s.surpluses = s.surpluses[0:nodeCount]
}

func (s *Surrogate) resize(nodeCount uint32) {
	if nodeCount <= s.nodeCount {
		return
	}

	if count := bufferGrowFactor * s.nodeCount; count > nodeCount {
		nodeCount = count
	}

	levels := make([]uint8, nodeCount*uint32(s.inCount))
	orders := make([]uint32, nodeCount*uint32(s.inCount))
	surpluses := make([]float64, nodeCount)

	copy(levels, s.levels[0:s.nodeCount*uint32(s.inCount)])
	copy(orders, s.orders[0:s.nodeCount*uint32(s.inCount)])
	copy(surpluses, s.surpluses[0:s.nodeCount])

	s.nodeCount = nodeCount
	s.levels = levels
	s.orders = orders
	s.surpluses = surpluses
}

// String returns a string containing human-friendly information about the
// surrogate.
func (s *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{ inputs: %d, levels: %d, nodes: %d }",
		s.inCount, s.level, s.nodeCount)
}

// Construct takes a function and yields a surrogate/interpolant for it, which
// can be further fed to Evaluate for the actual interpolation.
func (self *Self) Construct(target func([]float64) []float64) *Surrogate {
	inc := uint32(self.grid.Dimensionality())

	surrogate := new(Surrogate)
	surrogate.initialize(uint16(inc))

	level := uint8(0)
	nodeCount := uint32(0)

	minValue := math.Inf(1)
	maxValue := math.Inf(-1)

	newc := uint32(1)
	oldc := uint32(0)

	levels := make([]uint8, newc*inc)
	orders := make([]uint32, newc*inc)

	for {
		surrogate.resize(oldc + newc)

		copy(surrogate.levels[oldc*inc:], levels)
		copy(surrogate.orders[oldc*inc:], orders)

		nodes := self.grid.ComputeNodes(levels, orders)
		values := target(nodes)

		for i := uint32(0); i < newc; i++ {
			surrogate.surpluses[oldc+i] = values[i] -
				self.evaluate(inc, oldc,
					nodes[i*inc:(i+1)*inc],
					surrogate.levels[0:oldc*inc],
					surrogate.orders[0:oldc*inc],
					surrogate.surpluses[0:oldc])
		}

		nodeCount += newc

		if level >= self.maxLevel {
			break
		}

		for i := range values {
			if values[i] < minValue {
				minValue = values[i]
			}
			if values[i] > maxValue {
				maxValue = values[i]
			}
		}

		if level >= self.minLevel {
			k := 0

			for i := uint32(0); i < newc; i++ {
				absError := math.Abs(surrogate.surpluses[oldc+i])
				relError := absError / (maxValue - minValue)

				if absError <= self.absTolerance &&
					relError <= self.relTolerance {

					continue
				}

				for j := uint32(0); j < inc; j++ {
					levels[k] = levels[i*inc+j]
					orders[k] = orders[i*inc+j]
					k++
				}
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

func (self *Self) evaluate(inc uint32, sc uint32, point []float64,
	levels []uint8, orders []uint32, surpluses []float64) (value float64) {

	for i := uint32(0); i < sc; i++ {
		value += surpluses[i] * self.basis.Evaluate(point,
			levels[i*inc:(i+1)*inc], orders[i*inc:(i+1)*inc])
	}

	return value
}

// Evaluate takes a surrogate produced by Construct and evaluates it at the
// given points.
func (self *Self) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	inc := uint32(self.grid.Dimensionality())
	pc := uint32(len(points)) / inc

	values := make([]float64, pc)

	for i := uint32(0); i < pc; i++ {
		values[i] = self.evaluate(inc, surrogate.nodeCount, points[i*inc:(i+1)*inc],
			surrogate.levels, surrogate.orders, surrogate.surpluses)
	}

	return values
}
