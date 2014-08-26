package newtoncotes

import (
	"reflect"
	"testing"
)

func assertEqual(expected, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("got '%v' instead of '%v'", actual, expected)
	}
}

func TestComputeOrders(t *testing.T) {
	expectedOrders := [][]uint32{{0}, {0, 2}, {1, 3}, {1, 3, 5, 7}}

	basis := New()

	for level := range expectedOrders {
		actualOrders := basis.ComputeOrders(uint8(level))

		assertEqual(expectedOrders[level], actualOrders, t)
	}
}

func TestComputeNodes(t *testing.T) {
	levels := []uint8{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 1, 3, 1, 3, 5, 7}

	expectedNodes := []float64{0.5, 0, 1, 0.25, 0.75, 0.125, 0.375, 0.625, 0.875}

	basis := New()
	actualNodes := basis.ComputeNodes(levels, orders)

	assertEqual(expectedNodes, actualNodes, t)
}

func TestComputeChildren(t *testing.T) {
	levels := []uint8{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 1, 3, 1, 3, 5, 7}

	expectedLevels := []uint8{1, 1, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4}
	expectedOrders := []uint32{0, 2, 1, 3, 1, 3, 5, 7, 1, 3, 5, 7, 9, 11, 13, 15}

	basis := New()
	actualLevels, actualOrders := basis.ComputeChildren(levels, orders)

	assertEqual(expectedLevels, actualLevels, t)
	assertEqual(expectedOrders, actualOrders, t)
}

func TestEvaluate(t *testing.T) {
	points := []float64{-1, 0, 0.5, 1, 2}
	levels := []uint8{0, 1, 1, 2, 2}
	orders := []uint32{0, 0, 2, 1, 3}
	surpluses := []float64{1, 2, 3, 4, 5}

	expectedValues := []float64{0, 3, 1, 4, 0}

	basis := New()

	for i := range points {
		actualValue := basis.Evaluate(points[i], levels, orders, surpluses)
		assertEqual(expectedValues[i], actualValue, t)
	}
}
