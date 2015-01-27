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
	result := []float64{0}
	c.EvaluateComposite(index, []float64{1}, point, result)
	return result[0]
}

// EvaluateCoposite computes a vector-valued weighted sum wherein each term is a
// weight vector multiplied by a multi-dimensional basis evaluated at a point.
func (c *Closed) EvaluateComposite(indices []uint64, weights, point, result []float64) {
	dc := int(c.dc)
	oc := len(result)
	nc := len(indices) / dc

	for i := 0; i < oc; i++ {
		result[i] = 0
	}

outer:
	for i := 0; i < nc; i++ {
		value := 1.0

		for j := 0; j < dc; j++ {
			if point[j] < 0 || 1 < point[j] {
				continue outer
			}

			level := uint32(indices[i*dc+j])

			if level == 0 {
				continue
			}

			order := uint32(indices[i*dc+j] >> 32)

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
