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
func NewClosed(dimensions uint, power uint) *Closed {
	return &Closed{dimensions, power}
}

// Compute evaluates a basis function.
func (self *Closed) Compute(index []uint64, point []float64) float64 {
	nd, np, value := self.nd, self.np, 1.0
	for i := uint(0); i < nd && value != 0.0; i++ {
		value *= closedCompute(levelMask&index[i], index[i]>>levelSize, np, point[i])
	}
	return value
}

// Integrate computes the integral of a basis function.
func (self *Closed) Integrate(index []uint64) float64 {
	nd, np, value := self.nd, self.np, 1.0
	for i := uint(0); i < nd && value != 0.0; i++ {
		value *= closedIntegrate(levelMask&index[i], index[i]>>levelSize, np)
	}
	return value
}

func closedCompute(level, order uint64, power uint, x float64) float64 {
	if level < uint64(power) {
		power = uint(level)
	}
	if power == 0 {
		return 1.0
	}

	xi, step := closedNode(level, order)

	delta := math.Abs(x - xi)
	if delta >= step {
		return 0.0
	}

	if power == 1 {
		// Use two linear segments. The reason is that, taking into account the
		// endpoints, there are three points available in order to construct a
		// first-order polynomial; however, such a polynomial can satisfy only
		// any two of them.
		return 1.0 - delta/step
	}

	value := 1.0

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

	return value
}

func closedIntegrate(level, order uint64, power uint) float64 {
	if level < uint64(power) {
		power = uint(level)
	}
	if power == 0 {
		return 1.0
	}

	x, step := closedNode(level, order)

	if power == 1 {
		// Use two liner segments. See the corresponding comment in
		// closedCompute.
		if level == 1 {
			return 0.25
		} else {
			return step
		}
	}

	// Use a Gauss–Legendre quadrature rule to integrate. Such a rule with n
	// nodes integrates exactly polynomials up to order 2*n - 1.
	nodes := uint(math.Ceil((float64(power) + 1.0) / 2.0))
	return integrate(x-step, x+step, nodes, func(x float64) float64 {
		return closedCompute(level, order, power, x)
	})
}

func closedNode(level, order uint64) (x, step float64) {
	if level == 0 {
		x, step = 0.5, 0.5
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
