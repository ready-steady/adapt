package local

import (
	"fmt"
	"math"

	"github.com/ready-steady/adapt/basis/polynomial"
	"github.com/ready-steady/adapt/grid/equidistant"
)

// Interpolation in one dimension.
func ExampleInterpolator_step() {
	const (
		inputs   = 1
		outputs  = 1
		absolute = 1e-4
	)

	grid, basis := equidistant.NewClosed(inputs), polynomial.NewClosed(inputs, 1)
	interpolator := New(grid, basis, NewConfig())

	target := NewTarget(inputs, outputs, absolute, func(x, y []float64) {
		if x[0] <= 0.5 {
			y[0] = 1.0
		} else {
			y[0] = 0.0
		}
	})

	surrogate := interpolator.Compute(target)

	fmt.Println(surrogate)

	// Output:
	// {inputs:1 outputs:1 nodes:18}
}

// Interpolation in two dimensions.
func ExampleInterpolator_cube() {
	const (
		inputs   = 2
		outputs  = 1
		absolute = 1e-4
	)

	grid, basis := equidistant.NewClosed(inputs), polynomial.NewClosed(inputs, 1)
	interpolator := New(grid, basis, NewConfig())

	target := NewTarget(inputs, outputs, absolute, func(x, y []float64) {
		if math.Abs(2.0*x[0]-1.0) < 0.45 && math.Abs(2.0*x[1]-1.0) < 0.45 {
			y[0] = 1.0
		} else {
			y[0] = 0.0
		}
	})

	surrogate := interpolator.Compute(target)

	fmt.Println(surrogate)

	// Output:
	// {inputs:2 outputs:1 nodes:377}
}
