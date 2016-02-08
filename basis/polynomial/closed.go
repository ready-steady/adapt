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

		x := point[i]
		xi, h := node(level, order)
		d := x - xi
		if d < 0.0 {
			d = -d
		}
		if d >= h {
			return 0.0 // value *= 0.0
		}

		for j := uint64(0); j < power; j++ {
			level, order = parent(level, order)
			xj, _ := node(level, order)
			value *= (x - xj) / (xi - xj)
		}
	}

	return value
}

// Integrate computes the integral of a basis function.
func (_ *Closed) Integrate(_ []uint64) float64 {
	return 0.0
}

func node(level, order uint64) (x, h float64) {
	if level == 0 {
		x, h = 0.5, 1.0
	} else {
		h = 1.0 / float64(uint64(2)<<(level-1))
		x = h * float64(order)
	}
	return
}

func parent(level, order uint64) (uint64, uint64) {
	switch level {
	case 0:
		panic("the root does not have a parent")
	case 1:
		level = 0
		order = 0
	case 2:
		level = 1
		order -= 1
	default:
		level -= 1
		if ((order-1)/2)%2 == 0 {
			order = (order + 1) / 2
		} else {
			order = (order - 1) / 2
		}
	}
	return level, order
}
