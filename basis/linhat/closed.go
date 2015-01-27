package linhat

import (
	"math"
)

// Closed represents an instance of the basis on [0, 1]^n.
type Closed struct {
	ic uint16
}

// NewClosed creates an instance of the basis on [0, 1]^n.
func NewClosed(dimensions uint16) *Closed {
	return &Closed{dimensions}
}

// EvaluateComposite computes a vector-valued weighted sum wherein each term is
// a vector of weights multiplied by a multi-dimensional basis function
// evaluated at a point.
func (c *Closed) EvaluateComposite(indices []uint64, weights, point, result []float64) {
	ic := int(c.ic)
	oc := len(result)
	nc := len(indices) / ic

	for i := 0; i < oc; i++ {
		result[i] = 0
	}

outer:
	for i := 0; i < nc; i++ {
		value := 1.0

		for j := 0; j < ic; j++ {
			if point[j] < 0 || 1 < point[j] {
				continue outer
			}

			level := uint32(indices[i*ic+j])

			if level == 0 {
				continue
			}

			order := uint32(indices[i*ic+j] >> 32)

			scale := float64(uint32(2) << (level - 1))
			distance := math.Abs(point[j] - float64(order)/scale)

			if distance >= 1/scale {
				continue outer
			}

			value *= 1 - scale*distance
		}

		for j := 0; j < oc; j++ {
			result[j] += weights[i*oc+j] * value
		}
	}
}
