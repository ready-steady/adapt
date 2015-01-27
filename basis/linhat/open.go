package linhat

import (
	"math"
)

// Open represents an instance of the basis on (0, 1)^n.
type Open struct {
	dc uint16
}

// NewOpen creates an instance of the basis on (0, 1)^n.
func NewOpen(dimensions uint16) *Open {
	return &Open{dimensions}
}

// Evaluate computes the value of the multi-dimensional basis function
// corresponding to the given index at the given point.
func (o *Open) Evaluate(index []uint64, point []float64) float64 {
	result := []float64{0}
	o.EvaluateComposite(index, []float64{1}, point, result)
	return result[0]
}

// EvaluateCoposite computes a vector-valued weighted sum wherein each term is a
// vector of weights multiplied by a multi-dimensional basis function evaluated
// at a point.
func (o *Open) EvaluateComposite(indices []uint64, weights, point, result []float64) {
	dc := int(o.dc)
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
