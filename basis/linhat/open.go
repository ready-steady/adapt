package linhat

// Open represents an instance of the basis in (0, 1)^n.
type Open struct {
	nd int
}

// NewOpen creates an instance of the basis in (0, 1)^n.
func NewOpen(dimensions uint) *Open {
	return &Open{int(dimensions)}
}

// Compute evaluates the value of a basis function at a point.
func (self *Open) Compute(index []uint64, point []float64) float64 {
	nd := self.nd

	value := 1.0

	for i := 0; i < nd; i++ {
		level := 0xFFFFFFFF & index[i]
		if level == 0 {
			continue // value *= 1
		}

		order := index[i] >> 32
		count := uint64(2)<<level - 1

		switch order {
		case 0:
			scale := float64(count + 1)
			if scale*point[i] < 2 {
				value *= 2 - scale*point[i]
			} else {
				return 0 // value *= 0
			}
		case count - 1:
			scale1, scale2 := float64(count-1), float64(count+1)
			if scale2*point[i] > scale1 {
				value *= scale2*point[i] - scale1
			} else {
				return 0 // value *= 0
			}
		default:
			scale := float64(count + 1)
			distance := point[i] - float64(order+1)/scale
			if distance < 0 {
				distance = -distance
			}
			if scale*distance < 1 {
				value *= 1 - scale*distance
			} else {
				return 0 // value *= 0
			}
		}
	}

	return value
}

// Integrate computes the integral of a basis function over the whole domain.
func (self *Open) Integrate(index []uint64) float64 {
	nd := self.nd

	value := 1.0

	for i := 0; i < nd; i++ {
		level := 0xFFFFFFFF & index[i]
		if level == 0 {
			continue // value *= 1
		}

		order := index[i] >> 32
		count := uint64(2)<<level - 1

		switch order {
		case 0, count - 1:
			value *= 2 / float64(count+1)
		default:
			value *= 1 / float64(count+1)
		}
	}

	return value
}
