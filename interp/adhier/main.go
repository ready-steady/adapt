// Package adhier provides an algorithm for adaptive hierarchical interpolation
// with local refinements.
package adhier

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
	s.inCount = inCount
	s.outCount = outCount
	s.nodeCount = bufferInitSize

	s.levels = make([]uint8, bufferInitSize*inCount)
	s.orders = make([]uint32, bufferInitSize*inCount)
	s.surpluses = make([]float64, bufferInitSize*outCount)
}

func (s *Surrogate) finalize(level uint8, nodeCount uint32) {
	s.level = level
	s.nodeCount = nodeCount

	s.levels = s.levels[0 : nodeCount*uint32(s.inCount)]
	s.orders = s.orders[0 : nodeCount*uint32(s.inCount)]
	s.surpluses = s.surpluses[0 : nodeCount*uint32(s.outCount)]
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
	surpluses := make([]float64, nodeCount*uint32(s.outCount))

	copy(levels, s.levels[0:s.nodeCount*uint32(s.inCount)])
	copy(orders, s.orders[0:s.nodeCount*uint32(s.inCount)])
	copy(surpluses, s.surpluses[0:s.nodeCount*uint32(s.outCount)])

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

	level := uint8(0)
	nodeCount := uint32(0)

	minValue := make([]float64, outc)
	maxValue := make([]float64, outc)
	for i := uint32(0); i < outc; i++ {
		minValue[i] = math.Inf(1)
		maxValue[i] = math.Inf(-1)
	}

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

		// Compute the surpluses corresponding to the active nodes.
		for i, k := uint32(0), oldc*outc; i < newc; i++ {
			value := evaluate(inc, self.basis, outc, nodes[i*inc:(i+1)*inc],
				oldc, surrogate.levels, surrogate.orders, surrogate.surpluses)

			for j := uint32(0); j < outc; j++ {
				surrogate.surpluses[k] = values[i*outc+j] - value[j]
				k++
			}
		}

		nodeCount += newc

		if level >= self.maxLevel {
			break
		}

		// Keep track of the maximal and minimal values of the function.
		for i, k := uint32(0), uint32(0); i < newc; i++ {
			for j := uint32(0); j < outc; j++ {
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
			k := uint32(0)

			for i := uint32(0); i < newc; i++ {
				refine := false

				for j := uint32(0); j < outc; j++ {
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

// Evaluate takes a surrogate produced by Construct and evaluates it at the
// given points.
func (self *Self) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	inc := uint32(self.grid.Dimensionality())
	outc := uint32(self.outCount)

	pointCount := uint32(len(points)) / inc

	values := make([]float64, pointCount*outc)

	for i, k := uint32(0), uint32(0); i < pointCount; i++ {
		value := evaluate(inc, self.basis, outc, points[i*inc:(i+1)*inc],
			surrogate.nodeCount, surrogate.levels, surrogate.orders,
			surrogate.surpluses)

		for j := uint32(0); j < outc; j++ {
			values[k] = value[j]
			k++
		}
	}

	return values
}

func evaluate(inc uint32, basis Basis, outc uint32, point []float64, nodec uint32,
	levels []uint8, orders []uint32, surpluses []float64) []float64 {

	value := make([]float64, outc)

	for i := uint32(0); i < nodec; i++ {
		weight := basis.Evaluate(point, levels[i*inc:(i+1)*inc], orders[i*inc:(i+1)*inc])
		for j := uint32(0); j < outc; j++ {
			value[j] += surpluses[i*outc+j] * weight
		}
	}

	return value
}
