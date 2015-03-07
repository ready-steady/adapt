package linhat

// Closed represents an instance of the basis on [0, 1]^n.
type Closed struct {
	ni int
}

// NewClosed creates an instance of the basis on [0, 1]^n.
func NewClosed(inputs uint) *Closed {
	return &Closed{int(inputs)}
}

// Compute evaluates the value of a multidimensional basis function at a point.
func (c *Closed) Compute(index []uint64, point []float64) float64 {
	ni := c.ni

	value := 1.0

	for i := 0; i < ni; i++ {
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
		if scale*distance < 1 {
			value *= 1 - scale*distance
		} else {
			return 0
		}
	}

	return value
}
