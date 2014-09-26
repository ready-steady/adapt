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
	value := 1.0

	for i := uint16(0); i < c.dc; i++ {
		if point[i] < 0 || 1 < point[i] {
			return 0
		}

		level := uint32(index[i])

		if level == 0 {
			continue
		}

		order := uint32(index[i] >> 32)

		scale := float64(uint32(2) << (level - 1))
		distance := math.Abs(point[i] - float64(order)/scale)

		if distance >= 1/scale {
			return 0
		}

		value *= 1 - scale*distance
	}

	return value
}
