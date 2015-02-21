package adhier

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/ready-steady/numeric/basis/linhat"
	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/support/assert"
)

func TestComputeStep(t *testing.T) {
	interpolator := prepare(&fixtureStep)
	surrogate := interpolator.Compute(step)
	assert.Equal(surrogate, fixtureStep.surrogate, t)
}

func TestEvaluateStep(t *testing.T) {
	interpolator := prepare(&fixtureStep)
	values := interpolator.Evaluate(fixtureStep.surrogate, fixtureStep.points)
	assert.Equal(values, fixtureStep.values, t)
}

func TestComputeHat(t *testing.T) {
	interpolator := prepare(&fixtureHat)
	surrogate := interpolator.Compute(hat)
	assert.Equal(surrogate, fixtureHat.surrogate, t)
}

func TestEvaluateHat(t *testing.T) {
	interpolator := prepare(&fixtureHat)
	values := interpolator.Evaluate(fixtureHat.surrogate, fixtureHat.points)
	assert.EqualWithin(values, fixtureHat.values, 1e-15, t)
}

func TestComputeCube(t *testing.T) {
	interpolator := prepare(&fixtureCube)
	surrogate := interpolator.Compute(cube)
	assert.Equal(surrogate, fixtureCube.surrogate, t)
}

func TestComputeBox(t *testing.T) {
	interpolator := prepare(&fixtureBox)
	surrogate := interpolator.Compute(box)
	assert.Equal(surrogate, fixtureBox.surrogate, t)
}

func TestEvaluateBox(t *testing.T) {
	interpolator := prepare(&fixtureBox)
	values := interpolator.Evaluate(fixtureBox.surrogate, fixtureBox.points)
	assert.EqualWithin(values, fixtureBox.values, 1e-15, t)
}

func BenchmarkComputeHat(b *testing.B) {
	interpolator := prepare(&fixtureHat)

	for i := 0; i < b.N; i++ {
		interpolator.Compute(hat)
	}
}

func BenchmarkComputeCube(b *testing.B) {
	interpolator := prepare(&fixtureCube, func(config *Config) {
		config.MaxLevel = 9
	})

	for i := 0; i < b.N; i++ {
		interpolator.Compute(cube)
	}
}

func BenchmarkComputeBox(b *testing.B) {
	interpolator := prepare(&fixtureBox, func(config *Config) {
		config.MaxLevel = 9
	})

	for i := 0; i < b.N; i++ {
		interpolator.Compute(box)
	}
}

func BenchmarkComputeMany(b *testing.B) {
	const (
		inputs  = 2
		outputs = 1000
	)

	interpolator := prepare(&fixture{
		surrogate: &Surrogate{
			Inputs:  inputs,
			Outputs: outputs,

			Level: 9,
		},
	})
	function := many(inputs, outputs)

	for i := 0; i < b.N; i++ {
		interpolator.Compute(function)
	}
}

func BenchmarkEvaluateHat(b *testing.B) {
	interpolator := prepare(&fixtureHat)
	surrogate := interpolator.Compute(hat)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		interpolator.Evaluate(surrogate, points)
	}
}

func BenchmarkEvaluateCube(b *testing.B) {
	interpolator := prepare(&fixtureCube, func(config *Config) {
		config.MaxLevel = 9
	})
	surrogate := interpolator.Compute(cube)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		interpolator.Evaluate(surrogate, points)
	}
}

func BenchmarkEvaluateBox(b *testing.B) {
	interpolator := prepare(&fixtureBox, func(config *Config) {
		config.MaxLevel = 9
	})
	surrogate := interpolator.Compute(box)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		interpolator.Evaluate(surrogate, points)
	}
}

func BenchmarkEvaluateMany(b *testing.B) {
	const (
		inputs  = 2
		outputs = 1000
	)

	interpolator := prepare(&fixture{
		surrogate: &Surrogate{
			Inputs:  inputs,
			Outputs: outputs,

			Level: 9,
		},
	})
	function := many(inputs, outputs)
	surrogate := interpolator.Compute(function)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		interpolator.Evaluate(surrogate, points)
	}
}

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

func prepare(fixture *fixture, arguments ...interface{}) *Interpolator {
	surrogate := fixture.surrogate

	ic, oc := surrogate.Inputs, surrogate.Outputs

	config := DefaultConfig(ic, oc)
	config.MaxLevel = surrogate.Level

	if len(arguments) > 0 {
		process, _ := arguments[0].(func(*Config))
		process(config)
	}

	interpolator, _ := New(newcot.NewClosed(ic), linhat.NewClosed(ic), config)

	return interpolator
}

func generate(surrogate *Surrogate) []float64 {
	const (
		count = 10000
	)

	generator := rand.New(rand.NewSource(0))
	points := make([]float64, count*surrogate.Inputs)
	for i := range points {
		points[i] = generator.Float64()
	}

	return points
}
