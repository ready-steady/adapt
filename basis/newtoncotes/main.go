// Package newtoncotes provides functions for working with the Newton–Cotes
// hierarchical basis on a unit hypercube including boundaries.
package newtoncotes

import (
	"math"
)

// Self contains the configuration of a particular instantiation of the basis
// and provides an access point for the accompanying functions.
type Self struct {
	dimensionCount uint16
}

// New creates an instance of the basis for the given dimensionality.
func New(dimensionCount uint16) *Self {
	return &Self{dimensionCount}
}

// ComputeNodes returns the nodes corresponding to the given levels and orders.
func (_ *Self) ComputeNodes(levels []uint8, orders []uint32) []float64 {
	nodes := make([]float64, len(levels))

	for i := range nodes {
		if levels[i] == 0 {
			nodes[i] = 0.5
		} else {
			nodes[i] = float64(orders[i]) / float64(uint32(2)<<(levels[i]-1))
		}
	}

	return nodes
}

// ComputeChildren returns the levels and orders of the child nodes
// corresponding to the parent nodes given by their levels and orders.
func (self *Self) ComputeChildren(levels []uint8, orders []uint32) ([]uint8, []uint32) {
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

func (_ *Self) Evaluate(point float64, levels []uint8, orders []uint32, surpluses []float64) (value float64) {
	for i := range surpluses {
		value += surpluses[i] * computeWeight(point, levels[i], orders[i])
	}

	return value
}
