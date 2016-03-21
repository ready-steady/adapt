package polynomial

import (
	"math"

	"github.com/ready-steady/adapt/internal"
)

// Open is a basis in (0, 1)^n.
type Open struct {
	nd uint
}

// NewOpen creates a basis.
func NewOpen(dimensions, power uint) *Open {
	if power != 1 {
		panic("not implemented")
	}
	return &Open{dimensions}
}

// Compute evaluates a basis function.
func (self *Open) Compute(index []uint64, point []float64) float64 {
	nd, value := self.nd, 1.0
	for i := uint(0); i < nd && value != 0.0; i++ {
		value *= openCompute(internal.LEVEL_MASK&index[i],
			index[i]>>internal.LEVEL_SIZE, point[i])
	}
	return value
}

// Integrate computes the integral of a basis function.
func (self *Open) Integrate(index []uint64) float64 {
	nd, value := self.nd, 1.0
	for i := uint(0); i < nd && value != 0.0; i++ {
		value *= openIntegrate(internal.LEVEL_MASK&index[i],
			index[i]>>internal.LEVEL_SIZE)
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
		h := 1.0 / float64(count+1)
		if x >= 2.0*h {
			return 0.0
		}
		return 2.0 - x/h
	case count - 1:
		h, left := 1.0/float64(count+1), float64(count-1)
		if x <= left*h {
			return 0.0
		}
		return x/h - left
	default:
		xi, h := openNode(level, order)
		Δ := math.Abs(x - xi)
		if Δ >= h {
			return 0.0
		}
		return 1.0 - Δ/h
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

func openNode(level, order uint64) (x, h float64) {
	h = 1.0 / float64(uint64(2)<<level)
	x = float64(order+1) * h
	return
}
