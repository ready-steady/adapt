package polynomial

import (
	"math"
)

// Open is a basis in (0, 1)^n.
type Open struct {
	nd uint
}

// NewOpen creates a basis in (0, 1)^n.
func NewOpen(dimensions, order uint) *Open {
	if order != 1 {
		panic("not implemented")
	}
	return &Open{dimensions}
}

// Compute evaluates a basis function.
func (self *Open) Compute(index []uint64, point []float64) float64 {
	nd, value := self.nd, 1.0
	for i := uint(0); i < nd && value != 0.0; i++ {
		value *= openCompute(levelMask&index[i], index[i]>>levelSize, point[i])
	}
	return value
}

// Integrate computes the integral of a basis function.
func (self *Open) Integrate(index []uint64) float64 {
	nd, value := self.nd, 1.0
	for i := uint(0); i < nd && value != 0.0; i++ {
		value *= openIntegrate(levelMask&index[i], index[i]>>levelSize)
	}
	return value
}

func openCompute(level, order uint64, x float64) float64 {
	if level == 0 {
		return 1.0
	}

	count := uint64(2)<<level - 1

	switch order {
	case 0:
		step := 1.0 / float64(count+1)
		if x >= 2.0*step {
			return 0.0
		}
		return 2.0 - x/step
	case count - 1:
		step, left := 1.0/float64(count+1), float64(count-1)
		if x <= left*step {
			return 0.0
		}
		return x/step - left
	default:
		xi, step := openNode(level, order)
		delta := math.Abs(x - xi)
		if delta >= step {
			return 0.0
		}
		return 1.0 - delta/step
	}
}

func openIntegrate(level, order uint64) float64 {
	if level == 0 {
		return 1.0
	}

	count := uint64(2)<<level - 1

	switch order {
	case 0, count - 1:
		return 2.0 / float64(count+1)
	default:
		return 1.0 / float64(count+1)
	}
}

func openNode(level, order uint64) (x, step float64) {
	step = 1.0 / float64(uint64(2)<<level)
	x = float64(order+1) * step
	return
}
