// Package linhat provides means for working with multi-dimensional bases
// formed by the linear hat function.
//
// Each function in the basis is identified by a sequence of levels and orders.
// Throughout the package, such a sequence is encoded as a sequence of uint64s,
// referred to as an index, where each uint64 is (level|order<<32).
package linhat

import (
	"math"
)

// Self represents a particular instantiation of the basis.
type Self struct {
	dc uint16
}

// New creates an instance of the basis.
func New(dimensions uint16) *Self {
	return &Self{dimensions}
}

// Evaluate computes the value of the multi-dimensional basis function
// corresponding to the given index at the given point.
func (self *Self) Evaluate(index []uint64, point []float64) float64 {
	var value, limit, delta float64 = 1, 0, 0

	for i := uint16(0); i < self.dc; i++ {
		if uint32(index[i]) == 0 {
			if math.Abs(point[i]-0.5) > 0.5 {
				return 0
			} else {
				continue
			}
		}

		limit = float64(uint32(2) << (uint32(index[i]) - 1))
		delta = math.Abs(point[i] - float64(uint32(index[i]>>32))/limit)

		if delta >= 1/limit {
			return 0
		}

		value *= 1 - limit*delta
	}

	return value
}
