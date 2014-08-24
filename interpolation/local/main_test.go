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
			if x[i] > 0.5 {
				y[i] = 1
			}
		}
		return y
	}

	basis := newtoncotes.New()

	algorithm := New(basis)
	algorithm.maximalLevel = 3

	surrogate := algorithm.Construct(target)

	assertEqual(surrogate.nodeCount, uint32(9), t)
	assertEqual(surrogate.levels, []uint8{0, 1, 1, 2, 2, 3, 3, 3, 3}, t)
	assertEqual(surrogate.orders, []uint32{0, 0, 2, 1, 3, 1, 3, 5, 7}, t)
	assertEqual(surrogate.surpluses, []float64{0, 0, 1, 0, 1, 0, 0, 1, 1}, t)
}
