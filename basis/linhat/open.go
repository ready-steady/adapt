package linhat

import (
	"math"
)

// Open represents an instance of the basis on (0, 1)^n.
type Open struct {
	ic uint16
	oc uint16
}

// NewOpen creates an instance of the basis on (0, 1)^n.
func NewOpen(inputs, outputs uint16) *Open {
	return &Open{inputs, outputs}
}

// Outputs returns the dimensionality of the output.
func (o *Open) Outputs() uint16 {
	return o.oc
}

// EvaluateComposite computes a vector-valued weighted sum wherein each term is
// a vector of weights multiplied by a multi-dimensional basis function
// evaluated at a point.
func (o *Open) EvaluateComposite(indices []uint64, weights, point, result []float64) {
	ic, oc := int(o.ic), int(o.oc)
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
			count := uint32(2)<<level - 1

			switch order {
			case 0:
				scale := float64(count + 1)
				if point[j] >= 2/scale {
					continue outer
				}
				value *= 2 - scale*point[j]
			case count - 1:
				scale1, scale2 := float64(count-1), float64(count+1)
				if point[j] <= scale1/scale2 {
					continue outer
				}
				value *= scale2*point[j] - scale1
			default:
				node := float64(order+1) / float64(count+1)
				scale, distance := float64(count+1), math.Abs(point[j]-node)
				if distance >= 1/scale {
					continue outer
				}
				value *= 1 - scale*distance
			}
		}

		for j := 0; j < oc; j++ {
			result[j] += weights[i*oc+j] * value
		}
	}
}
