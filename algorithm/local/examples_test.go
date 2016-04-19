package local

import (
	"fmt"
	"math"

	"github.com/ready-steady/adapt/basis/polynomial"
	"github.com/ready-steady/adapt/grid/equidistant"
)

// Interpolation in one dimension.
func ExampleAlgorithm_step() {
	const (
		inputs          = 1
		outputs         = 1
		minLevel        = 1
		maxLevel        = 10
		localError      = 1e-4
		polynomialOrder = 1
	)

	grid := equidistant.NewClosed(inputs)
	basis := polynomial.NewClosed(inputs, polynomialOrder)
	algorithm := New(inputs, outputs, grid, basis)
	strategy := NewStrategy(inputs, outputs, grid, minLevel, maxLevel, localError)

	surrogate := algorithm.Compute(func(x, y []float64) {
		if x[0] <= 0.5 {
			y[0] = 1.0
		} else {
			y[0] = 0.0
		}
	}, strategy)

	fmt.Println(surrogate)

	// Output:
	// {inputs:1 outputs:1 nodes:20}
}

// Interpolation in two dimensions.
func ExampleAlgorithm_cube() {
	const (
		inputs          = 2
		outputs         = 1
		minLevel        = 1
		maxLevel        = 10
		localError      = 1e-4
		polynomialOrder = 1
	)

	grid := equidistant.NewClosed(inputs)
	basis := polynomial.NewClosed(inputs, polynomialOrder)
	algorithm := New(inputs, outputs, grid, basis)
	strategy := NewStrategy(inputs, outputs, grid, minLevel, maxLevel, localError)

	surrogate := algorithm.Compute(func(x, y []float64) {
		if math.Abs(2.0*x[0]-1.0) < 0.45 && math.Abs(2.0*x[1]-1.0) < 0.45 {
			y[0] = 1.0
		} else {
			y[0] = 0.0
		}
	}, strategy)

	fmt.Println(surrogate)

	// Output:
	// {inputs:2 outputs:1 nodes:477}
}
