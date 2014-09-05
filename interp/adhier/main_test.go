package adhier

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/go-math/numan/basis/linhat"
	"github.com/go-math/numan/grid/newcot"
)

func assertEqual(actual, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got '%v' instead of '%v'", actual, expected)
	}
}

const epsilon = 1e-8

func assertAlmostEqual(actual, expected []float64, t *testing.T) {
	if len(actual) != len(expected) {
		goto error
	}

	for i := range actual {
		if math.Abs(actual[i]-expected[i]) > epsilon {
			goto error
		}
	}

	return

error:
	t.Fatalf("got '%v' instead of '%v'", actual, expected)
}

func TestConstructStep(t *testing.T) {
	algorithm := New(newcot.New(1), linhat.New(1), 1)
	algorithm.maxLevel = stepFixture.surrogate.level

	surrogate := algorithm.Construct(step)

	assertEqual(surrogate, stepFixture.surrogate, t)
}

func TestEvaluateStep(t *testing.T) {
	algorithm := New(newcot.New(1), linhat.New(1), 1)

	values := algorithm.Evaluate(stepFixture.surrogate, stepFixture.points)

	assertEqual(values, stepFixture.values, t)
}

func TestConstructCube(t *testing.T) {
	algorithm := New(newcot.New(2), linhat.New(2), 1)
	algorithm.maxLevel = cubeFixture.surrogate.level

	surrogate := algorithm.Construct(cube)

	assertEqual(surrogate, cubeFixture.surrogate, t)
}

func TestConstructBox(t *testing.T) {
	algorithm := New(newcot.New(2), linhat.New(2), 3)
	algorithm.maxLevel = boxFixture.surrogate.level

	surrogate := algorithm.Construct(box)

	assertEqual(surrogate, boxFixture.surrogate, t)
}

func TestEvaluateBox(t *testing.T) {
	algorithm := New(newcot.New(2), linhat.New(2), 3)

	values := algorithm.Evaluate(boxFixture.surrogate, boxFixture.points)

	assertAlmostEqual(values, boxFixture.values, t)
}

func BenchmarkHat(b *testing.B) {
	algorithm := New(newcot.New(1), linhat.New(1), 1)

	for i := 0; i < b.N; i++ {
		_ = algorithm.Construct(hat)
	}
}

func BenchmarkCube(b *testing.B) {
	algorithm := New(newcot.New(2), linhat.New(2), 1)

	for i := 0; i < b.N; i++ {
		_ = algorithm.Construct(cube)
	}
}

func BenchmarkBox(b *testing.B) {
	algorithm := New(newcot.New(2), linhat.New(2), 3)

	for i := 0; i < b.N; i++ {
		_ = algorithm.Construct(box)
	}
}

func BenchmarkMany(b *testing.B) {
	const (
		inCount  = 2
		outCount = 1000
	)

	algorithm := New(newcot.New(inCount), linhat.New(inCount), outCount)
	function := many(inCount, outCount)

	for i := 0; i < b.N; i++ {
		_ = algorithm.Construct(function)
	}
}

// A one-input-one-output scenario with a non-smooth function.
func ExampleSelf_step() {
	const (
		inCount  = 1
		outCount = 1
	)

	grid := newcot.New(inCount)
	basis := linhat.New(inCount)

	algorithm := New(grid, basis, outCount)
	algorithm.maxLevel = 19

	surrogate := algorithm.Construct(step)
	fmt.Println(surrogate)

	// Output:
	// Surrogate{ inputs: 1, outputs: 1, levels: 19, nodes: 38 }
}

// A one-input-one-output scenario with a smooth function.
func ExampleSelf_hat() {
	const (
		inCount  = 1
		outCount = 1
	)

	grid := newcot.New(inCount)
	basis := linhat.New(inCount)

	algorithm := New(grid, basis, outCount)
	algorithm.maxLevel = 9

	surrogate := algorithm.Construct(hat)
	fmt.Println(surrogate)

	// Output:
	// Surrogate{ inputs: 1, outputs: 1, levels: 9, nodes: 305 }
}

// A multiple-input-one-output scenario with a non-smooth function.
func ExampleSelf_cube() {
	const (
		inCount  = 2
		outCount = 1
	)

	grid := newcot.New(inCount)
	basis := linhat.New(inCount)

	algorithm := New(grid, basis, outCount)
	algorithm.maxLevel = 9

	surrogate := algorithm.Construct(cube)
	fmt.Println(surrogate)

	// Output:
	// Surrogate{ inputs: 2, outputs: 1, levels: 9, nodes: 377 }
}

// A multiple-input-many-output scenario with a non-smooth function.
func ExampleSelf_many() {
	const (
		inCount  = 2
		outCount = 1000
	)

	grid := newcot.New(inCount)
	basis := linhat.New(inCount)

	algorithm := New(grid, basis, outCount)

	surrogate := algorithm.Construct(many(inCount, outCount))
	fmt.Println(surrogate)

	// Output:
	// Surrogate{ inputs: 2, outputs: 1000, levels: 9, nodes: 362 }
}
