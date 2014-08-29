// Package linhat provides means for working with multi-dimensional bases
// formed by the linear hat function.
package linhat

import (
	"math"
)

// Self represents a particular instantiation of the basis.
type Self struct {
	dimensionCount uint16
}

// New creates an instance of the basis.
func New(dimensionCount uint16) *Self {
	return &Self{dimensionCount}
}

// Evaluate computes the value of a multi-dimensional basis function,
// identified by the given levels and orders, at the given point.
func (self *Self) Evaluate(levels []uint8, orders []uint32, point []float64) float64 {
	var value, limit, delta float64 = 1, 0, 0

	for i := uint16(0); i < self.dimensionCount; i++ {
		if levels[i] == 0 {
			if math.Abs(point[i]-0.5) > 0.5 {
				return 0
			} else {
				continue
			}
		}

		limit = float64(uint32(2) << (levels[i] - 1))
		delta = math.Abs(point[i] - float64(orders[i])/limit)

		if delta >= 1/limit {
			return 0
		}

		value *= 1 - limit*delta
	}

	return value
}
