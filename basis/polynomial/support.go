package polynomial

import (
	"math"

	"github.com/ready-steady/adapt/internal"
)

func compute(index []uint64, point []float64, nd uint,
	compute func(uint64, uint64, float64) float64) float64 {

	value := 1.0
	for i := uint(0); i < nd && value != 0.0; i++ {
		value *= compute(index[i]&internal.LEVEL_MASK,
			index[i]>>internal.LEVEL_SIZE, point[i])
	}
	return value
}

func equal(one, two float64) bool {
	const ε = 1e-14 // ~= 2^(-46)
	return one == two || math.Abs(one-two) < ε
}

func integrate(index []uint64, nd uint,
	integrate func(uint64, uint64) float64) float64 {

	value := 1.0
	for i := uint(0); i < nd && value != 0.0; i++ {
		value *= integrate(index[i]&internal.LEVEL_MASK,
			index[i]>>internal.LEVEL_SIZE)
	}
	return value
}
