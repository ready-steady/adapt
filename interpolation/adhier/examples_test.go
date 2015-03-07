package adhier

import (
	"fmt"
	"math"

	"github.com/ready-steady/numeric/basis/linhat"
	"github.com/ready-steady/numeric/grid/newcot"
)

// Interpolation in one dimension.
func ExampleInterpolator_step() {
	const (
		inputs  = 1
		outputs = 1
	)

	target := func(x, y []float64, _ []uint64) {
		if x[0] <= 0.5 {
			y[0] = 1
		} else {
			y[0] = 0
		}
	}

	grid, basis := newcot.NewClosed(inputs), linhat.NewClosed(inputs)

	config := DefaultConfig(inputs, outputs)
	config.MaxLevel = 19

	interpolator, _ := New(grid, basis, config)
	surrogate := interpolator.Compute(target)

	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 1, outputs: 1, level: 19, nodes: 38}
}

// Interpolation in two dimensions.
func ExampleInterpolator_cube() {
	const (
		inputs  = 2
		outputs = 1
	)

	target := func(x, y []float64, _ []uint64) {
		if math.Abs(2*x[0]-1) < 0.45 && math.Abs(2*x[1]-1) < 0.45 {
			y[0] = 1
		} else {
			y[0] = 0
		}
	}

	grid, basis := newcot.NewClosed(inputs), linhat.NewClosed(inputs)

	config := DefaultConfig(inputs, outputs)
	config.MaxLevel = 9

	interpolator, _ := New(grid, basis, config)
	surrogate := interpolator.Compute(target)

	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 2, outputs: 1, level: 9, nodes: 377}
}
