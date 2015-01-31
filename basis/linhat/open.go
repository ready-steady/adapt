package linhat

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

// Evaluate computes the value of a multi-dimensional basis function at a point.
func (o *Open) Evaluate(index []uint64, point []float64) float64 {
	ic := int(o.ic)

	value := 1.0

	for i := 0; i < ic; i++ {
		if point[i] < 0 || 1 < point[i] {
			return 0
		}

		level := uint32(index[i])
		if level == 0 {
			continue
		}

		order := uint32(index[i] >> 32)
		count := uint32(2)<<level - 1

		switch order {
		case 0:
			scale := float64(count + 1)
			if point[i] >= 2/scale {
				return 0
			}
			value *= 2 - scale*point[i]
		case count - 1:
			scale1, scale2 := float64(count-1), float64(count+1)
			if point[i] <= scale1/scale2 {
				return 0
			}
			value *= scale2*point[i] - scale1
		default:
			node := float64(order+1) / float64(count+1)
			scale, distance := float64(count+1), point[i]-node
			if distance < 0 {
				distance = -distance
			}
			if distance >= 1/scale {
				return 0
			}
			value *= 1 - scale*distance
		}
	}

	return value
}
