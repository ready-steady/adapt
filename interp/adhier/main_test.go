package adhier

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gomath/format/mat"
	"github.com/gomath/numan/basis/linhat"
	"github.com/gomath/numan/grid/newcot"
)

func assertEqual(actual, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got '%v' instead of '%v'", actual, expected)
	}
}

// TestConstructStep deals with a one-input-one-output scenario.
func TestConstructStep(t *testing.T) {
	algorithm := New(newcot.New(1), linhat.New(1), 1)
	algorithm.maxLevel = stepFixture.surrogate.level

	surrogate := algorithm.Construct(step)

	assertEqual(surrogate, stepFixture.surrogate, t)
}

// TestEvaluateStep deals with a one-input-one-output scenario.
func TestEvaluateStep(t *testing.T) {
	algorithm := New(newcot.New(1), linhat.New(1), 1)

	values := algorithm.Evaluate(stepFixture.surrogate, stepFixture.points)

	assertEqual(values, stepFixture.values, t)
}

// TestConstructCube deals with a multiple-input-one-output scenario.
func TestConstructCube(t *testing.T) {
	algorithm := New(newcot.New(2), linhat.New(2), 1)
	algorithm.maxLevel = cubeFixture.surrogate.level

	surrogate := algorithm.Construct(cube)

	assertEqual(surrogate, cubeFixture.surrogate, t)
}

// TestConstructBox deals with a multiple-input-multiple-output scenario.
func TestConstructBox(t *testing.T) {
	algorithm := New(newcot.New(2), linhat.New(2), 3)
	algorithm.maxLevel = boxFixture.surrogate.level

	surrogate := algorithm.Construct(box)

	assertEqual(surrogate, boxFixture.surrogate, t)
}

// ExampleStep demonstrates a one-input-one-output scenario with a smooth
// function.
func ExampleStep() {
	algorithm := New(newcot.New(1), linhat.New(1), 1)
	algorithm.maxLevel = 19
	surrogate := algorithm.Construct(step)

	fmt.Println(surrogate)

	// Output:
	// Surrogate{ inputs: 1, outputs: 1, levels: 19, nodes: 38 }
}

// ExampleHat demonstrates a one-input-one-output scenario with a non-smooth
// function.
func ExampleHat() {
	algorithm := New(newcot.New(1), linhat.New(1), 1)
	algorithm.maxLevel = 9
	surrogate := algorithm.Construct(hat)

	fmt.Println(surrogate)

	if !testing.Verbose() {
		return
	}

	points := makeGrid1D(101)
	values := algorithm.Evaluate(surrogate, points)

	file, _ := mat.Open("hat.mat", "w7.3")
	defer file.Close()

	file.PutMatrix("x", 1, 101, points)
	file.PutMatrix("y", 1, 101, values)

	// Output:
	// Surrogate{ inputs: 1, outputs: 1, levels: 9, nodes: 305 }
}

// ExampleCube demonstrates a multiple-input-one-output scenario with a
// non-smooth function.
func ExampleCube() {
	algorithm := New(newcot.New(2), linhat.New(2), 1)
	algorithm.maxLevel = 9
	surrogate := algorithm.Construct(cube)

	fmt.Println(surrogate)

	if !testing.Verbose() {
		return
	}

	points := makeGrid2D(21)
	values := algorithm.Evaluate(surrogate, points)

	file, _ := mat.Open("cube.mat", "w7.3")
	defer file.Close()

	file.PutMatrix("x", 2, 21*21, points)
	file.PutMatrix("y", 1, 21*21, values)

	// Output:
	// Surrogate{ inputs: 2, outputs: 1, levels: 9, nodes: 377 }
}

func BenchmarkConstructHat(b *testing.B) {
	algorithm := New(newcot.New(1), linhat.New(1), 1)

	for i := 0; i < b.N; i++ {
		_ = algorithm.Construct(hat)
	}
}

func BenchmarkConstructCube(b *testing.B) {
	algorithm := New(newcot.New(2), linhat.New(2), 1)

	for i := 0; i < b.N; i++ {
		_ = algorithm.Construct(cube)
	}
}

func BenchmarkConstructBox(b *testing.B) {
	algorithm := New(newcot.New(2), linhat.New(2), 3)

	for i := 0; i < b.N; i++ {
		_ = algorithm.Construct(box)
	}
}
