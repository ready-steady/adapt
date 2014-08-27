package newtoncotes

import (
	"math"
)

type Instance struct {
}

func New() *Instance {
	return new(Instance)
}

func computeNode(level uint8, order uint32) float64 {
	if level == 0 {
		return 0.5
	} else {
		return float64(order) / float64(uint32(2)<<(level-1))
	}
}

func (_ *Instance) ComputeNodes(levels []uint8, orders []uint32) []float64 {
	nodes := make([]float64, len(levels))

	for i := range nodes {
		nodes[i] = computeNode(levels[i], orders[i])
	}

	return nodes
}

func (_ *Instance) ComputeChildren(levels []uint8, orders []uint32) ([]uint8, []uint32) {
	count := len(levels)

	childLevels := make([]uint8, 2*count)
	childOrders := make([]uint32, 2*count)

	k := 0

	for i := 0; i < count; i++ {
		switch levels[i] {
		case 0:
			childLevels[k] = 1
			childOrders[k] = 0
			k++

			childLevels[k] = 1
			childOrders[k] = 2
			k++

		case 1:
			childLevels[k] = 2
			childOrders[k] = orders[i] + 1
			k++

		default:
			childLevels[k] = levels[i] + 1
			childOrders[k] = 2*orders[i] - 1
			k++

			childLevels[k] = levels[i] + 1
			childOrders[k] = 2*orders[i] + 1
			k++
		}
	}

	return childLevels[0:k], childOrders[0:k]
}

func computeWeight(point float64, level uint8, order uint32) float64 {
	if level == 0 {
		if math.Abs(point-0.5) > 0.5 {
			return 0
		} else {
			return 1
		}
	}

	limit := float64(uint32(2) << (level - 1))
	delta := math.Abs(point - float64(order)/limit)

	if delta >= 1/limit {
		return 0
	} else {
		return 1 - limit*delta
	}
}

func (_ *Instance) Evaluate(point float64, levels []uint8, orders []uint32, surpluses []float64) (value float64) {
	for i := range surpluses {
		value += surpluses[i] * computeWeight(point, levels[i], orders[i])
	}

	return value
}
