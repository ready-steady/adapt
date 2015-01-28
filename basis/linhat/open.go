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

// EvaluateComposite computes a number of vector-valued weighted sums wherein
// each sum is composed of weight vectors multiplied by multi-dimensional basis
// functions evaluated at a point. The basis functions are identified by
// indices; the weight vectors are stored in weights; and the points
// corresponding to the sums are stored in points.
func (o *Open) EvaluateComposite(indices []uint64, weights, points []float64) []float64 {
	ic, oc := int(o.ic), int(o.oc)
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
				result[k*oc+j] += weights[i*oc+j] * value
			}
		}
	}

	return result
}
