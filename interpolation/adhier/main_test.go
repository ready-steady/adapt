package adhier

import (
	"fmt"
	"testing"

	"github.com/ready-steady/numeric/basis/linhat"
	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/support/assert"
)

func TestConstructStep(t *testing.T) {
	fixtureStep.prepare()
	algorithm := makeInterpolator(1, 1, fixtureStep.surrogate.level)

	surrogate := algorithm.Compute(step)

	assert.Equal(surrogate, fixtureStep.surrogate, t)
}

func TestEvaluateStep(t *testing.T) {
	fixtureStep.prepare()
	algorithm := makeInterpolator(1, 1, 0)

	values := algorithm.Evaluate(fixtureStep.surrogate, fixtureStep.points)

	assert.Equal(values, fixtureStep.values, t)
}

func TestConstructCube(t *testing.T) {
	fixtureCube.prepare()
	algorithm := makeInterpolator(2, 1, fixtureCube.surrogate.level)

	surrogate := algorithm.Compute(cube)

	assert.Equal(surrogate, fixtureCube.surrogate, t)
}

func TestConstructBox(t *testing.T) {
	fixtureBox.prepare()
	algorithm := makeInterpolator(2, 3, fixtureBox.surrogate.level)

	surrogate := algorithm.Compute(box)

	assert.Equal(surrogate, fixtureBox.surrogate, t)
}

func TestEvaluateBox(t *testing.T) {
	fixtureBox.prepare()
	algorithm := makeInterpolator(2, 3, 0)

	values := algorithm.Evaluate(fixtureBox.surrogate, fixtureBox.points)

	assert.AlmostEqual(values, fixtureBox.values, t)
}

func BenchmarkHat(b *testing.B) {
	algorithm := makeInterpolator(1, 1, 0)

	for i := 0; i < b.N; i++ {
		_ = algorithm.Compute(hat)
	}
}

func BenchmarkCube(b *testing.B) {
	algorithm := makeInterpolator(2, 1, 0)

	for i := 0; i < b.N; i++ {
		_ = algorithm.Compute(cube)
	}
}

func BenchmarkBox(b *testing.B) {
	algorithm := makeInterpolator(2, 3, 0)

	for i := 0; i < b.N; i++ {
		_ = algorithm.Compute(box)
	}
}

func BenchmarkMany(b *testing.B) {
	algorithm := makeInterpolator(2, 1000, 0)
	function := many(2, 1000)

	for i := 0; i < b.N; i++ {
		_ = algorithm.Compute(function)
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

	config := DefaultConfig()
	config.MaxLevel = 19
	algorithm, _ := New(grid, basis, config, outputs)

	surrogate := algorithm.Compute(step)
	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 1, outputs: 1, levels: 19, nodes: 38}
}

// A one-input-one-output scenario with a smooth function.
func ExampleInterpolator_hat() {
	const (
		inputs  = 1
		outputs = 1
	)

	grid := newcot.NewClosed(inputs)
	basis := linhat.NewClosed(inputs)

	config := DefaultConfig()
	config.MaxLevel = 9
	algorithm, _ := New(grid, basis, config, outputs)

	surrogate := algorithm.Compute(hat)
	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 1, outputs: 1, levels: 9, nodes: 305}
}

// A multiple-input-one-output scenario with a non-smooth function.
func ExampleInterpolator_cube() {
	const (
		inputs  = 2
		outputs = 1
	)

	grid := newcot.NewClosed(inputs)
	basis := linhat.NewClosed(inputs)

	config := DefaultConfig()
	config.MaxLevel = 9
	algorithm, _ := New(grid, basis, config, outputs)

	surrogate := algorithm.Compute(cube)
	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 2, outputs: 1, levels: 9, nodes: 377}
}

// A multiple-input-many-output scenario with a non-smooth function.
func ExampleInterpolator_many() {
	const (
		inputs  = 2
		outputs = 1000
	)

	grid := newcot.NewClosed(inputs)
	basis := linhat.NewClosed(inputs)

	algorithm, _ := New(grid, basis, DefaultConfig(), outputs)

	surrogate := algorithm.Compute(many(inputs, outputs))
	fmt.Println(surrogate)

	// Output:
	// Surrogate{inputs: 2, outputs: 1000, levels: 9, nodes: 362}
}

func makeInterpolator(ic, oc uint16, ml uint8) *Interpolator {
	config := DefaultConfig()
	if ml > 0 {
		config.MaxLevel = ml
	}

	interpolator, _ := New(newcot.NewClosed(ic), linhat.NewClosed(ic), config, oc)

	return interpolator
}
