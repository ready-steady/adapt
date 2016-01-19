package local

import (
	"fmt"
	"math"

	"github.com/ready-steady/adapt/basis/linear"
	"github.com/ready-steady/adapt/grid/equidistant"
)

// Interpolation in one dimension.
func ExampleInterpolator_step() {
	const (
		inputs    = 1
		outputs   = 1
		tolerance = 1e-4
	)

	grid, basis := equidistant.NewClosed(inputs), linear.NewClosed(inputs)
	interpolator := New(grid, basis, NewConfig())

	target := NewTarget(inputs, outputs, func(x, y []float64) {
		if x[0] <= 0.5 {
			y[0] = 1.0
		} else {
			y[0] = 0.0
		}
	})

	metric := NewMetric(outputs, tolerance)

	surrogate := interpolator.Compute(target, metric)

	fmt.Println(surrogate)

	// Output:
	// {inputs:1 outputs:1 level:9 nodes:18}
}

// Interpolation in two dimensions.
func ExampleInterpolator_cube() {
	const (
		inputs    = 2
		outputs   = 1
		tolerance = 1e-4
	)

	grid, basis := equidistant.NewClosed(inputs), linear.NewClosed(inputs)
	interpolator := New(grid, basis, NewConfig())

	target := NewTarget(inputs, outputs, func(x, y []float64) {
		if math.Abs(2.0*x[0]-1.0) < 0.45 && math.Abs(2.0*x[1]-1.0) < 0.45 {
			y[0] = 1.0
		} else {
			y[0] = 0.0
		}
	})

	metric := NewMetric(outputs, tolerance)

	surrogate := interpolator.Compute(target, metric)

	fmt.Println(surrogate)

	// Output:
	// {inputs:2 outputs:1 level:9 nodes:377}
}
