package linear

// Open is a basis in (0, 1)^n.
type Open struct {
	nd int
}

// NewOpen creates a basis in (0, 1)^n.
func NewOpen(dimensions uint) *Open {
	return &Open{int(dimensions)}
}

// Compute evaluates a basis function.
func (self *Open) Compute(index []uint64, point []float64) float64 {
	nd := self.nd

	value := 1.0
	for i := 0; i < nd; i++ {
		level := levelMask & index[i]
		if level == 0 {
			continue // value *= 1.0
		}

		order := index[i] >> levelSize
		count := uint64(2)<<level - 1

		x := point[i]
		switch order {
		case 0:
			step := 1.0 / float64(count+1)
			if x >= 2.0*step {
				return 0.0 // value *= 0.0
			}
			value *= 2.0 - x/step
		case count - 1:
			step, left := 1.0/float64(count+1), float64(count-1)
			if x <= left*step {
				return 0.0 // value *= 0.0
			}
			value *= x/step - left
		default:
			step := 1.0 / float64(count+1)
			delta := x - float64(order+1)*step
			if delta < 0.0 {
				delta = -delta
			}
			if delta >= step {
				return 0.0 // value *= 0.0
			}
			value *= 1.0 - delta/step
		}
	}

	return value
}

// Integrate computes the integral of a basis function.
func (self *Open) Integrate(index []uint64) float64 {
	nd := self.nd

	value := 1.0
	for i := 0; i < nd; i++ {
		level := levelMask & index[i]
		if level == 0 {
			continue // value *= 1.0
		}

		order := index[i] >> levelSize
		count := uint64(2)<<level - 1

		switch order {
		case 0, count - 1:
			value *= 2.0 / float64(count+1)
		default:
			value *= 1.0 / float64(count+1)
		}
	}

	return value
}
