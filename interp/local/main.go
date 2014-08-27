package local

import (
	"fmt"
	"math"

	"github.com/gomath/numan/basis"
)

const (
	initialBufferSize = 200
	bufferGrowFactor  = 2
)

type Instance struct {
	basis        basis.Interface
	minLevel     uint8
	maxLevel     uint8
	absTolerance float64
	relTolerance float64
}

func New(basis basis.Interface) *Instance {
	return &Instance{
		basis:        basis,
		minLevel:     2 - 1,
		maxLevel:     10 - 1,
		absTolerance: 1e-4,
		relTolerance: 1e-2,
	}
}

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

func (self *Instance) Construct(target func([]float64) []float64) *Surrogate {
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

		nodes := self.basis.ComputeNodes(levels, orders)
		values := target(nodes)

		for i := uint32(0); i < newCount; i++ {
			surrogate.surpluses[oldCount+i] = values[i] -
				self.basis.Evaluate(nodes[i], surrogate.levels[0:oldCount],
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

		levels, orders = self.basis.ComputeChildren(levels, orders)

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

func (self *Instance) Evaluate(surrogate *Surrogate, points []float64) []float64 {
	values := make([]float64, len(points))

	for i := range values {
		values[i] = self.basis.Evaluate(points[i], surrogate.levels,
			surrogate.orders, surrogate.surpluses)
	}

	return values
}
