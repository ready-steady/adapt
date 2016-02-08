package linear

// Closed is a basis in [0, 1]^n.
type Closed struct {
	nd uint
}

// NewClosed creates a basis in [0, 1]^n.
func NewClosed(dimensions uint) *Closed {
	return &Closed{dimensions}
}

// Compute evaluates a basis function.
func (self *Closed) Compute(index []uint64, point []float64) float64 {
	nd := self.nd

	value := 1.0
	for i := uint(0); i < nd; i++ {
		level := levelMask & index[i]
		if level == 0 {
			continue // value *= 1.0
		}

		order := index[i] >> levelSize

		step := 1.0 / float64(uint64(2)<<(level-1))

		delta := point[i] - float64(order)*step
		if delta < 0.0 {
			delta = -delta
		}
		if delta >= step {
			return 0.0 // value *= 0.0
		}

		value *= 1.0 - delta/step
	}

	return value
}

// Integrate computes the integral of a basis function.
func (self *Closed) Integrate(index []uint64) float64 {
	nd := self.nd

	value := 1.0
	for i := uint(0); i < nd; i++ {
		level := levelMask & index[i]
		switch level {
		case 0:
			// value *= 1.0
		case 1:
			value *= 0.25
		default:
			value *= 1.0 / float64(uint64(2)<<(level-1))
		}
	}

	return value
}
