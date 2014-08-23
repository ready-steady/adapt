package newton_cotes

import (
	"reflect"
	"testing"
)

func assertEqual(expected, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Got", actual, "instead of", expected)
	}
}

func TestComputeOrders(t *testing.T) {
	expectedOrders := [][]uint16{{0}, {0, 2}, {1, 3}, {1, 3, 5, 7}}

	for level := range expectedOrders {
		actualOrders := ComputeOrders(uint8(level))

		assertEqual(expectedOrders[level], actualOrders, t)
	}
}

func TestComputeNodes(t *testing.T) {
	levels := []uint8{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint16{0, 0, 2, 1, 3, 1, 3, 5, 7}

	expectedNodes := []float64{0.5, 0, 1, 0.25, 0.75, 0.125, 0.375, 0.625, 0.875}
	actualNodes := ComputeNodes(levels, orders)

	assertEqual(expectedNodes, actualNodes, t)
}

func TestComputeChildren(t *testing.T) {
	levels := []uint8{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint16{0, 0, 2, 1, 3, 1, 3, 5, 7}

	expectedLevels := []uint8{1, 1, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4}
	expectedOrders := []uint16{0, 2, 1, 3, 1, 3, 5, 7, 1, 3, 5, 7, 9, 11, 13, 15}

	actualLevels, actualOrders := ComputeChildren(levels, orders)

	assertEqual(expectedLevels, actualLevels, t)
	assertEqual(expectedOrders, actualOrders, t)
}
