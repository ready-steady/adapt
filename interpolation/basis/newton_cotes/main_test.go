package newton_cotes

import (
	"testing"
	"reflect"
)

func assert_equal(expected, actual interface{}, t *testing.T) {
	if !reflect.DeepEqual(expected, actual) {
		t.Error("Expected", expected, "instead of", actual)
	}
}

func TestComputeLevelOrders(t *testing.T) {
	expected := [][]uint16{{0}, {0, 2}, {1, 3}, {1, 3, 5, 7}}

	for level := range expected {
		assert_equal(expected[level], ComputeLevelOrders(uint8(level)), t)
	}
}

func TestComputeNodes(t *testing.T) {
	levels := []uint8{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint16{0, 0, 2, 1, 3, 1, 3, 5, 7}
	expected := []float64{0.5, 0, 1, 0.25, 0.75, 0.125, 0.375, 0.625, 0.875}

	assert_equal(expected, ComputeNodes(levels, orders), t)
}
