package linhat

// Closed represents an instance of the basis on [0, 1]^n.
type Closed struct {
	ic uint
}

// NewClosed creates an instance of the basis on [0, 1]^n.
func NewClosed(inputs uint) *Closed {
	return &Closed{inputs}
}

// Evaluate computes the value of a multi-dimensional basis function at a point.
func (c *Closed) Evaluate(index []uint64, point []float64) float64 {
	ic := int(c.ic)

	value := 1.0

	for i := 0; i < ic; i++ {
		level := 0xFFFFFFFF & index[i]
		if level == 0 {
			continue
		}

		order := index[i] >> 32

		scale := float64(uint64(2) << (level - 1))
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
