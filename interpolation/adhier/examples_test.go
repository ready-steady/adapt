package adhier

import (
	"fmt"

	"github.com/ready-steady/numeric/basis/linhat"
	"github.com/ready-steady/numeric/grid/newcot"
)

// A one-input-one-output scenario with a non-smooth function.
func ExampleInterpolator_step() {
	const (
		inputs  = 1
		outputs = 1
	)

	grid := newcot.NewClosed(inputs)
	basis := linhat.NewClosed(inputs)

	config := DefaultConfig(inputs, outputs)
	config.MaxLevel = 19
	interpolator, _ := New(grid, basis, config)

	surrogate := interpolator.Compute(step)
	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 1, outputs: 1, level: 19, nodes: 38}
}

// A one-input-one-output scenario with a smooth function.
func ExampleInterpolator_hat() {
	const (
		inputs  = 1
		outputs = 1
	)

	grid := newcot.NewClosed(inputs)
	basis := linhat.NewClosed(inputs)

	config := DefaultConfig(inputs, outputs)
	config.MaxLevel = 9
	interpolator, _ := New(grid, basis, config)

	surrogate := interpolator.Compute(hat)
	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 1, outputs: 1, level: 9, nodes: 305}
}

// A multiple-input-one-output scenario with a non-smooth function.
func ExampleInterpolator_cube() {
	const (
		inputs  = 2
		outputs = 1
	)

	grid := newcot.NewClosed(inputs)
	basis := linhat.NewClosed(inputs)

	config := DefaultConfig(inputs, outputs)
	config.MaxLevel = 9
	interpolator, _ := New(grid, basis, config)

	surrogate := interpolator.Compute(cube)
	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 2, outputs: 1, level: 9, nodes: 377}
}

// A multiple-input-many-output scenario with a non-smooth function.
func ExampleInterpolator_many() {
	const (
		inputs  = 2
		outputs = 1000
	)

	grid := newcot.NewClosed(inputs)
	basis := linhat.NewClosed(inputs)
	config := DefaultConfig(inputs, outputs)
	config.MaxNodes = 300

	interpolator, _ := New(grid, basis, config)

	surrogate := interpolator.Compute(many(inputs, outputs))
	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 2, outputs: 1000, level: 9, nodes: 300}
}
