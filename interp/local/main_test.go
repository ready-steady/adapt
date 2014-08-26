package local

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gomath/format/mat"
	"github.com/gomath/numan/basis/newtoncotes"
)

func assertEqual(actual, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got '%v' instead of '%v'", actual, expected)
	}
}

func TestConstruct(t *testing.T) {
	algorithm := New(newtoncotes.New())
	algorithm.maxLevel = 4

	surrogate := algorithm.Construct(step)

	levels := []uint8{0, 1, 1, 2, 3, 3, 4, 4}
	orders := []uint32{0, 0, 2, 3, 5, 7, 9, 11}
	surpluses := []float64{1, 0, -1, -0.5, -0.5, 0, -0.5, 0}

	assertEqual(surrogate.level, uint8(4), t)
	assertEqual(surrogate.nodeCount, uint32(8), t)

	assertEqual(surrogate.levels, levels, t)
	assertEqual(surrogate.orders, orders, t)
	assertEqual(surrogate.surpluses, surpluses, t)
}

func TestEvaluate(t *testing.T) {
	algorithm := New(newtoncotes.New())

	surrogate := &Surrogate {
		level: 4,
		nodeCount: 8,
		levels: []uint8{0, 1, 1, 2, 3, 3, 4, 4},
		orders: []uint32{0, 0, 2, 3, 5, 7, 9, 11},
		surpluses: []float64{1, 0, -1, -0.5, -0.5, 0, -0.5, 0},
	}

	points := []float64{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1}
	values := []float64{1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0}

	assertEqual(algorithm.Evaluate(surrogate, points), values, t)
}

func ExampleStep() {
	algorithm := New(newtoncotes.New())
	algorithm.maxLevel = 20 - 1
	surrogate := algorithm.Construct(step)

	fmt.Println(surrogate)
	// Output:
	// Surrogate{ levels: 20, nodes: 38 }
}

func ExampleHat() {
	algorithm := New(newtoncotes.New())
	algorithm.maxLevel = 10 - 1
	surrogate := algorithm.Construct(hat)

	fmt.Println(surrogate)

	if !testing.Verbose() { return }

	points := make([]float64, 101)
	for i := range points { points[i] = 0.01 * float64(i) }
	values := algorithm.Evaluate(surrogate, points)

	file, _ := mat.Open("hat.mat", "w7.3")
	defer file.Close()

	file.PutMatrix("x", 101, 1, points)
	file.PutMatrix("y", 101, 1, values)

	// Output:
	// Surrogate{ levels: 10, nodes: 305 }
}

func step(x []float64) []float64 {
	y := make([]float64, len(x))
	for i := range x {
		if x[i] <= 0.5 {
			y[i] = 1
		}
	}
	return y
}

func hat(x []float64) []float64 {
	y := make([]float64, len(x))
	for i, z := range x {
		z = 5 * z - 1
		switch {
		case 0 <= z && z < 1: y[i] = 0.5 * z * z
		case 1 <= z && z < 2: y[i] = 0.5 * (-2 * z * z + 6 * z - 3)
		case 2 <= z && z < 3: y[i] = 0.5 * (3 - z) * (3 - z)
		}
	}
	return y
}
