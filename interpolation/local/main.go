package local

import (
	"github.com/gomath/numerical/basis"
)

const (
	initialBufferSize = 200
	bufferGrowFactor = 2
)

type Instance struct {
	basis             basis.Interface
	minimalLevel	  uint8
	maximalLevel      uint8
	absoluteTolerance float64
	relativeTolerance float64
}

func New(basis basis.Interface) Instance {
	return Instance{
		basis: basis,
		minimalLevel: 2,
		maximalLevel: 10,
		absoluteTolerance: 1e-4,
		relativeTolerance: 1e-2,
	}
}

type Surrogate struct {
	levels    []uint8
	orders    []uint32
	surpluses []float64
	nodeCount uint32
}

func newSurrogate() Surrogate {
	return Surrogate{
		levels: make([]uint8, initialBufferSize),
		orders: make([]uint32, initialBufferSize),
		surpluses: make([]float64, initialBufferSize),
	}
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

func (instance *Surrogate) trim(nodeCount uint32) {
	instance.levels = instance.levels[0:nodeCount]
	instance.orders = instance.orders[0:nodeCount]
	instance.surpluses = instance.surpluses[0:nodeCount]
	instance.nodeCount = nodeCount
}

func (instance Instance) Construct(target func([]float64) []float64) Surrogate {
	surrogate := newSurrogate()

	level := uint8(0)
	nodeCount := uint32(0)

	orders := instance.basis.ComputeOrders(level)
	levels := make([]uint8, len(orders))

	for {
		count := uint32(len(levels))

		if count == 0 {
			break
		}

		surrogate.resize(nodeCount + count)

		copy(surrogate.levels[nodeCount:], levels)
		copy(surrogate.orders[nodeCount:], orders)

		nodes := instance.basis.ComputeNodes(levels, orders)
		values := target(nodes)

		copy(surrogate.surpluses[nodeCount:], values)

		nodeCount += count

		if level >= instance.maximalLevel {
			break
		}

		levels, orders = instance.basis.ComputeChildren(levels, orders)

		level++
	}

	surrogate.trim(nodeCount)

	return surrogate
}
