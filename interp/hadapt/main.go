// Package hadapt provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package hadapt

import (
	"fmt"
	"math"
)

const (
	initialBufferSize = 200
	bufferGrowFactor  = 2
)

// Grid is the interface that an sparse grid should satisfy in order to be used
// in the algorithm.
type Grid interface {
	ComputeNodes(levels []uint8, orders []uint32) []float64
	ComputeChildren(levels []uint8, orders []uint32) ([]uint8, []uint32)
}

// Basis is the interface that a functional basis should satisfy in order to be
// used in the algorithm.
type Basis interface {
	Evaluate(point float64, level uint8, order uint32) float64
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
		minLevel:     2 - 1,
		maxLevel:     10 - 1,
		absTolerance: 1e-4,
		relTolerance: 1e-2,
	}
}

// Surrogate is the result of Construct, which represents an interpolant for a
// function.
type Surrogate struct {
	level     uint8
	nodeCount uint32

	levels    []uint8
	orders    []uint32
	surpluses []float64
}

func (s *Surrogate) initialize() {
	s.levels = make([]uint8, initialBufferSize)
	s.orders = make([]uint32, initialBufferSize)
	s.surpluses = make([]float64, initialBufferSize)
}

func (s *Surrogate) finalize(level uint8, nodeCount uint32) {
	s.level = level
	s.nodeCount = nodeCount

	s.levels = s.levels[0:nodeCount]
	s.orders = s.orders[0:nodeCount]
	s.surpluses = s.surpluses[0:nodeCount]
}

// String returns a string containing human-friendly information about the
// surrogate.
func (s *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{ levels: %d, nodes: %d }", s.level+1, s.nodeCount)
}

func (s *Surrogate) resize(size uint32) {
	currentSize := uint32(len(s.levels))

	if size <= currentSize {
		return
	}

	if grownSize := bufferGrowFactor * currentSize; grownSize > size {
		size = grownSize
	}

	levels := make([]uint8, size)
	orders := make([]uint32, size)
	surpluses := make([]float64, size)

	copy(levels, s.levels[0:currentSize])
	copy(orders, s.orders[0:currentSize])
	copy(surpluses, s.surpluses[0:currentSize])

	s.levels = levels
	s.orders = orders
	s.surpluses = surpluses
}

// Construct takes a function and yields a surrogate/interpolant for it, which
// can be further fed to Evaluate for the actual interpolation.
func (self *Self) Construct(target func([]float64) []float64) *Surrogate {
	surrogate := new(Surrogate)
	surrogate.initialize()

	level := uint8(0)
	nodeCount := uint32(0)

	minValue := math.Inf(1)
	maxValue := math.Inf(-1)

	newCount := uint32(1)
	oldCount := uint32(0)

	levels := make([]uint8, newCount)
	orders := make([]uint32, newCount)

	for {
		surrogate.resize(oldCount + newCount)

		copy(surrogate.levels[oldCount:], levels)
		copy(surrogate.orders[oldCount:], orders)

		nodes := self.grid.ComputeNodes(levels, orders)
		values := target(nodes)

		for i := uint32(0); i < newCount; i++ {
			surrogate.surpluses[oldCount+i] = values[i] -
				self.evaluate(nodes[i], surrogate.levels[0:oldCount],
					surrogate.orders[0:oldCount], surrogate.surpluses[0:oldCount])
		}

		nodeCount += newCount

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

			for i := uint32(0); i < newCount; i++ {
				absError := math.Abs(surrogate.surpluses[oldCount+i])
				relError := absError / (maxValue - minValue)

				if absError <= self.absTolerance &&
					relError <= self.relTolerance {

					continue
				}

				levels[k] = levels[i]
				orders[k] = orders[i]

				k++
			}

			levels = levels[0:k]
			orders = orders[0:k]
		}

		levels, orders = self.grid.ComputeChildren(levels, orders)

		oldCount += newCount
		newCount = uint32(len(levels))

		if newCount == 0 {
			break
		}

		level++
	}

	surrogate.finalize(level, nodeCount)
	return surrogate
}

func (self *Self) evaluate(point float64, levels []uint8, orders []uint32, surpluses []float64) (value float64) {
	for i := range surpluses {
		value += surpluses[i] * self.basis.Evaluate(point, levels[i], orders[i])
	}

	return value
}

// Evaluate takes a surrogate produced by Construct and evaluates it at the
// given points.
func (self *Self) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	values := make([]float64, len(points))

	for i := range values {
		values[i] = self.evaluate(points[i], surrogate.levels,
			surrogate.orders, surrogate.surpluses)
	}

	return values
}
