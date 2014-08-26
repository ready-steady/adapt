package local

import (
	"fmt"
	"math"

	"github.com/gomath/numan/basis"
)

const (
	initialBufferSize = 200
	bufferGrowFactor = 2
)

type Instance struct {
	basis             basis.Interface
	minimalLevel      uint8
	maximalLevel      uint8
	absoluteTolerance float64
	relativeTolerance float64
}

func New(basis basis.Interface) *Instance {
	return &Instance {
		basis: basis,
		minimalLevel: 2 - 1,
		maximalLevel: 10 - 1,
		absoluteTolerance: 1e-4,
		relativeTolerance: 1e-2,
	}
}

type Surrogate struct {
	level     uint8
	nodeCount uint32

	levels    []uint8
	orders    []uint32
	surpluses []float64
}

func (instance *Surrogate) initialize() {
	instance.levels = make([]uint8, initialBufferSize)
	instance.orders = make([]uint32, initialBufferSize)
	instance.surpluses = make([]float64, initialBufferSize)
}

func (instance *Surrogate) finalize(level uint8, nodeCount uint32) {
	instance.level = level
	instance.nodeCount = nodeCount

	instance.levels = instance.levels[0:nodeCount]
	instance.orders = instance.orders[0:nodeCount]
	instance.surpluses = instance.surpluses[0:nodeCount]
}

func (instance *Surrogate) String() string {
	return fmt.Sprintf("Surrogate{ levels: %d, nodes: %d }",
		instance.level + 1, instance.nodeCount)
}

func (instance *Surrogate) resize(size uint32) {
	currentSize := uint32(len(instance.levels))

	if size <= currentSize {
		return
	}
	
	if grownSize := bufferGrowFactor * currentSize; grownSize > size {
		size = grownSize
	}

	levels := make([]uint8, size)
	orders := make([]uint32, size)
	surpluses := make([]float64, size)

	copy(levels, instance.levels[0:currentSize])
	copy(orders, instance.orders[0:currentSize])
	copy(surpluses, instance.surpluses[0:currentSize])

	instance.levels = levels
	instance.orders = orders
	instance.surpluses = surpluses
}

func (instance *Instance) Construct(target func([]float64) []float64) *Surrogate {
	surrogate := new(Surrogate)
	surrogate.initialize()

	level := uint8(0)
	nodeCount := uint32(0)

	minimalValue := math.Inf(1)
	maximalValue := math.Inf(-1)

	orders := instance.basis.ComputeOrders(level)

	newCount := uint32(len(orders))
	oldCount := uint32(0)

	levels := make([]uint8, newCount)

	for {
		surrogate.resize(oldCount + newCount)

		copy(surrogate.levels[oldCount:], levels)
		copy(surrogate.orders[oldCount:], orders)

		nodes := instance.basis.ComputeNodes(levels, orders)
		values := target(nodes)

		for i := uint32(0); i < newCount; i++ {
			surrogate.surpluses[oldCount + i] = values[i] -
				instance.basis.Evaluate(nodes[i], surrogate.levels[0:oldCount],
					surrogate.orders[0:oldCount], surrogate.surpluses[0:oldCount])
		}

		nodeCount += newCount

		if level >= instance.maximalLevel {
			break
		}

		for i := range values {
			if values[i] < minimalValue {
				minimalValue = values[i]
			}
			if values[i] > maximalValue {
				maximalValue = values[i]
			}
		}

		if level >= instance.minimalLevel {
			k := 0

			for i := uint32(0); i < newCount; i++ {
				absoluteError := math.Abs(surrogate.surpluses[oldCount + i])
				relativeError := absoluteError / (maximalValue - minimalValue)

				if absoluteError <= instance.absoluteTolerance &&
					relativeError <= instance.relativeTolerance {

					continue;
				}

				levels[k] = levels[i]
				orders[k] = orders[i]

				k++
			}

			levels = levels[0:k]
			orders = orders[0:k]
		}

		levels, orders = instance.basis.ComputeChildren(levels, orders)

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
