package local

import (
	"reflect"
	"testing"

	"github.com/gomath/numerical/basis/newtoncotes"
)

func assertEqual(expected, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Got", actual, "instead of", expected)
	}
}

func TestCompute(t *testing.T) {
	target := func(x []float64) []float64 {
		y := make([]float64, len(x))
		for i := range x {
			if x[i] <= 0.5 {
				y[i] = 1
			}
		}
		return y
	}

	basis := newtoncotes.New()

	algorithm := New(basis)
	algorithm.maximalLevel = 4

	surrogate := algorithm.Construct(target)

	expectedLevels := []uint8{0, 1, 1, 2, 3, 3, 4, 4}
	expectedOrders := []uint32{0, 0, 2, 3, 5, 7, 9, 11}
	expectedSurpluses := []float64{1, 0, -1, -0.5, -0.5, 0, -0.5, 0}

	assertEqual(uint32(len(expectedLevels)), surrogate.nodeCount, t)
	assertEqual(expectedLevels, surrogate.levels, t)
	assertEqual(expectedOrders, surrogate.orders, t)
	assertEqual(expectedSurpluses, surrogate.surpluses, t)
}
