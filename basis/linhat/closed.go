package linhat

import (
	"math"
)

// Closed represents an instance of the basis on [0, 1]^n.
type Closed struct {
	dc uint16
}

// NewClosed creates an instance of the basis on [0, 1]^n.
func NewClosed(dimensions uint16) *Closed {
	return &Closed{dimensions}
}

// Evaluate computes the value of the multi-dimensional basis function
// corresponding to the given index at the given point.
func (c *Closed) Evaluate(index []uint64, point []float64) float64 {
	var value, limit, delta float64 = 1, 0, 0

	for i := uint16(0); i < c.dc; i++ {
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
