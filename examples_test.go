package adapt

import (
	"fmt"
	"math"

	"github.com/ready-steady/adapt/basis/linhat"
	"github.com/ready-steady/adapt/grid/newcot"
)

// Interpolation in one dimension.
func ExampleInterpolator_step() {
	const (
		inputs    = 1
		outputs   = 1
		tolerance = 1e-4
	)

	grid, basis := newcot.NewClosed(inputs), linhat.NewClosed(inputs)
	interpolator := New(grid, basis, NewConfig())

	target := NewTarget(inputs, outputs)
	target.ComputeHandler = func(x, y []float64) {
		if x[0] <= 0.5 {
			y[0] = 1
		} else {
			y[0] = 0
		}
	}
	target.RefineHandler = func(_, surplus []float64, _ float64) float64 {
		if math.Abs(surplus[0]) > tolerance {
			return 1
		} else {
			return 0
		}
	}

	surrogate := interpolator.Compute(target)

	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 1, outputs: 1, level: 9, nodes: 18}
}

// Interpolation in two dimensions.
func ExampleInterpolator_cube() {
	const (
		inputs    = 2
		outputs   = 1
		tolerance = 1e-4
	)

	grid, basis := newcot.NewClosed(inputs), linhat.NewClosed(inputs)
	interpolator := New(grid, basis, NewConfig())

	target := NewTarget(inputs, outputs)
	target.ComputeHandler = func(x, y []float64) {
		if math.Abs(2*x[0]-1) < 0.45 && math.Abs(2*x[1]-1) < 0.45 {
			y[0] = 1
		} else {
			y[0] = 0
		}
	}
	target.RefineHandler = func(_, surplus []float64, _ float64) float64 {
		if math.Abs(surplus[0]) > tolerance {
			return 1
		} else {
			return 0
		}
	}

	surrogate := interpolator.Compute(target)

	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 2, outputs: 1, level: 9, nodes: 377}
}