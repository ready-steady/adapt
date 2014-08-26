package local

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/gomath/numerical/basis/newtoncotes"
)

func assertEqual(expected, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Got", actual, "instead of", expected)
	}
}

func TestConstruct(t *testing.T) {
	algorithm := New(newtoncotes.New())
	algorithm.maximalLevel = 4

	surrogate := algorithm.Construct(step)

	expectedLevels := []uint8{0, 1, 1, 2, 3, 3, 4, 4}
	expectedOrders := []uint32{0, 0, 2, 3, 5, 7, 9, 11}
	expectedSurpluses := []float64{1, 0, -1, -0.5, -0.5, 0, -0.5, 0}

	assertEqual(uint8(4), surrogate.level, t)
	assertEqual(uint32(8), surrogate.nodeCount, t)

	assertEqual(expectedLevels, surrogate.levels, t)
	assertEqual(expectedOrders, surrogate.orders, t)
	assertEqual(expectedSurpluses, surrogate.surpluses, t)
}

func ExampleStep() {
	algorithm := New(newtoncotes.New())
	algorithm.maximalLevel = 20 - 1
	surrogate := algorithm.Construct(step)

	fmt.Println(surrogate)
	// Output:
	// Surrogate{ levels: 20, nodes: 38 }
}

func ExampleHat() {
	algorithm := New(newtoncotes.New())
	algorithm.maximalLevel = 10 - 1
	surrogate := algorithm.Construct(hat)

	fmt.Println(surrogate)
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
