package polynomial

import (
	"math"
)

// Closed is a basis in [0, 1]^n.
type Closed struct {
	nd uint
	np uint
}

// NewClosed creates a basis in [0, 1]^n.
func NewClosed(dimensions uint, order uint) *Closed {
	return &Closed{dimensions, order}
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
		xi, step := closedNode(level, order)

		x := point[i]
		delta := math.Abs(x - xi)
		if delta >= step {
			return 0.0 // value *= 0.0
		}

		if power == 1 {
			// Use two linear segments. The reason is that there are three
			// points (including the endpoints), but a first-order polynomial
			// can satisfy only two of them in general.
			value *= 1.0 - delta/step
			continue
		}

		// The left endpoint of the local support.
		xl := xi - step
		value *= (x - xl) / (xi - xl)
		power -= 1

		// The right endpoint of the local support.
		xr := xi + step
		value *= (x - xr) / (xi - xr)
		power -= 1

		// Find the rest of the needed ancestors.
		for power > 0 {
			level, order = closedParent(level, order)
			xj, _ := closedNode(level, order)
			if equal(xj, xl) || equal(xj, xr) {
				continue
			}
			value *= (x - xj) / (xi - xj)
			power -= 1
		}
	}

	return value
}

// Integrate computes the integral of a basis function.
func (self *Closed) Integrate(index []uint64) float64 {
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
		_, step := closedNode(level, order)

		if power == 1 {
			// Use two liner segments. See the corresponding comment in Compute.
			if level == 1 {
				value *= 0.25
			} else {
				value *= step
			}
			continue
		}

		panic("not implemented")
	}

	return value
}

func closedNode(level, order uint64) (x, step float64) {
	if level == 0 {
		x, step = 0.5, 1.0
	} else {
		step = 1.0 / float64(uint64(2)<<(level-1))
		x = float64(order) * step
	}
	return
}

func closedParent(level, order uint64) (uint64, uint64) {
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
