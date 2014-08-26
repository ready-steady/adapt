package newtoncotes

import (
	"reflect"
	"testing"
)

func assertEqual(actual, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got '%v' instead of '%v'", actual, expected)
	}
}

func TestComputeOrders(t *testing.T) {
	basis := New()

	orders := [][]uint32{{0}, {0, 2}, {1, 3}, {1, 3, 5, 7}}

	for level := range orders {
		assertEqual(basis.ComputeOrders(uint8(level)), orders[level], t)
	}
}

func TestComputeNodes(t *testing.T) {
	basis := New()

	levels := []uint8{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 1, 3, 1, 3, 5, 7}
	nodes := []float64{0.5, 0, 1, 0.25, 0.75, 0.125, 0.375, 0.625, 0.875}

	assertEqual(basis.ComputeNodes(levels, orders), nodes, t)
}

func TestComputeChildren(t *testing.T) {
	basis := New()

	levels := []uint8{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 1, 3, 1, 3, 5, 7}
	childLevels := []uint8{1, 1, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4}
	childOrders := []uint32{0, 2, 1, 3, 1, 3, 5, 7, 1, 3, 5, 7, 9, 11, 13, 15}

	levels, orders = basis.ComputeChildren(levels, orders)

	assertEqual(levels, childLevels, t)
	assertEqual(orders, childOrders, t)
}

func TestEvaluate(t *testing.T) {
	basis := New()

	points := []float64{-1, 0, 0.5, 1, 2}
	levels := []uint8{0, 1, 1, 2, 2}
	orders := []uint32{0, 0, 2, 1, 3}
	surpluses := []float64{1, 2, 3, 4, 5}
	values := []float64{0, 3, 1, 4, 0}

	assertEqual(basis.Evaluate(points, levels, orders, surpluses), values, t)
}
