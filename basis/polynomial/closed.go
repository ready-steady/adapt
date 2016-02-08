package polynomial

// Closed is a basis in [0, 1]^n.
type Closed struct {
	nd uint
	np uint
}

// New creates a basis in [0, 1]^n.
func New(dimensions uint, polynomialOrder uint) *Closed {
	return &Closed{dimensions, polynomialOrder}
}

// Compute evaluates a basis function.
func (self *Closed) Compute(index []uint64, point []float64) float64 {
	nd, np := self.nd, uint64(self.np)

	value := 1.0
	for i := uint(0); i < nd; i++ {
		level, power := levelMask&index[i], np
		if level < power {
			power = level
		}
		if power == 0 {
			continue // value *= 1.0
		}

		order := index[i] >> levelSize

		scale := float64(uint64(2) << (level - 1))
		distance := point[i] - float64(order)/scale
		if distance < 0.0 {
			distance = -distance
		}
		if scale*distance >= 1.0 {
			return 0.0 // value *= 0.0
		}

		switch power {
		case 1:
			value *= 1.0 - scale*distance
		default:
			panic("Not implemented yet")
		}
	}

	return value
}

// Integrate computes the integral of a basis function.
func (_ *Closed) Integrate(_ []uint64) float64 {
	return 0.0
}
