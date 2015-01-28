package linhat

import (
	"math"
)

// Closed represents an instance of the basis on [0, 1]^n.
type Closed struct {
	ic uint16
	oc uint16
}

// NewClosed creates an instance of the basis on [0, 1]^n.
func NewClosed(inputs, outputs uint16) *Closed {
	return &Closed{inputs, outputs}
}

// EvaluateComposite computes a number of vector-valued weighted sums wherein
// each sum is composed of weight vectors multiplied by multi-dimensional basis
// functions evaluated at a point. The basis functions are identified by
// indices; the weight vectors are stored in weights; and the points
// corresponding to the sums are stored in points.
func (c *Closed) EvaluateComposite(indices []uint64, weights, points []float64) []float64 {
	ic, oc := int(c.ic), int(c.oc)
	nc := len(indices) / ic
	pc := len(points) / ic

	result := make([]float64, pc*oc)

	for k := 0; k < pc; k++ {
		point := points[k*ic : (k+1)*ic]

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
				result[k*oc+j] += weights[i*oc+j] * value
			}
		}
	}

	return result
}
