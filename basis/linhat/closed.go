package linhat

// Closed represents an instance of the basis on [0, 1]^n.
type Closed struct {
	ic uint16
	oc uint16
}

// NewClosed creates an instance of the basis on [0, 1]^n.
func NewClosed(inputs, outputs uint16) *Closed {
	return &Closed{inputs, outputs}
}

// Outputs returns the dimensionality of the output.
func (c *Closed) Outputs() uint16 {
	return c.oc
}

// Evaluate computes the value of a multi-dimensional basis function at a point.
func (c *Closed) Evaluate(index []uint64, point []float64) float64 {
	ic := int(c.ic)

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

		scale := float64(uint32(2) << (level - 1))
		distance := point[i] - float64(order)/scale
		if distance < 0 {
			distance = -distance
		}

		if distance >= 1/scale {
			return 0
		}

		value *= 1 - scale*distance
	}

	return value
}
