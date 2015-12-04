package linear

// Open represents an instance of the basis in (0, 1)^n.
type Open struct {
	nd int
}

// NewOpen creates an instance of the basis in (0, 1)^n.
func NewOpen(dimensions uint) *Open {
	return &Open{int(dimensions)}
}

// Compute evaluates the value of a basis function.
func (self *Open) Compute(index []uint64, point []float64) float64 {
	nd := self.nd

	value := 1.0

	for i := 0; i < nd; i++ {
		level := LEVEL_MASK & index[i]
		if level == 0 {
			continue // value *= 1.0
		}

		order := index[i] >> LEVEL_SIZE
		count := uint64(2)<<level - 1

		switch order {
		case 0:
			scale := float64(count + 1)
			if scale*point[i] < 2.0 {
				value *= 2.0 - scale*point[i]
			} else {
				return 0.0 // value *= 0.0
			}
		case count - 1:
			scale1, scale2 := float64(count-1), float64(count+1)
			if scale2*point[i] > scale1 {
				value *= scale2*point[i] - scale1
			} else {
				return 0.0 // value *= 0.0
			}
		default:
			scale := float64(count + 1)
			distance := point[i] - float64(order+1)/scale
			if distance < 0.0 {
				distance = -distance
			}
			if scale*distance < 1.0 {
				value *= 1.0 - scale*distance
			} else {
				return 0.0 // value *= 0.0
			}
		}
	}

	return value
}

// Integrate computes the integral of a basis function.
func (self *Open) Integrate(index []uint64) float64 {
	nd := self.nd

	value := 1.0

	for i := 0; i < nd; i++ {
		level := LEVEL_MASK & index[i]
		if level == 0 {
			continue // value *= 1.0
		}

		order := index[i] >> LEVEL_SIZE
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
