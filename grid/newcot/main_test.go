package newcot

import (
	"reflect"
	"testing"
)

func assertEqual(actual, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got '%v' instead of '%v'", actual, expected)
	}
}

func TestComputeNodes1D(t *testing.T) {
	grid := New(1)

	levels := []uint8{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 1, 3, 1, 3, 5, 7}
	nodes := []float64{0.5, 0, 1, 0.25, 0.75, 0.125, 0.375, 0.625, 0.875}

	assertEqual(grid.ComputeNodes(levels, orders), nodes, t)
}

func TestComputeNodes2D(t *testing.T) {
	grid := New(2)

	levels := []uint8{
		0, 0,
		0, 1,
		0, 1,
		1, 0,
		1, 0,
		0, 2,
		0, 2,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		2, 0,
		2, 0,
	}

	orders := []uint32{
		0, 0,
		0, 0,
		0, 2,
		0, 0,
		2, 0,
		0, 1,
		0, 3,
		0, 0,
		0, 2,
		2, 0,
		2, 2,
		1, 0,
		3, 0,
	}

	nodes := []float64{
		0.50, 0.50,
		0.50, 0.00,
		0.50, 1.00,
		0.00, 0.50,
		1.00, 0.50,
		0.50, 0.25,
		0.50, 0.75,
		0.00, 0.00,
		0.00, 1.00,
		1.00, 0.00,
		1.00, 1.00,
		0.25, 0.50,
		0.75, 0.50,
	}

	assertEqual(grid.ComputeNodes(levels, orders), nodes, t)
}

func TestComputeChildren1D(t *testing.T) {
	grid := New(1)

	levels := []uint8{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 1, 3, 1, 3, 5, 7}
	childLevels := []uint8{1, 1, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4}
	childOrders := []uint32{0, 2, 1, 3, 1, 3, 5, 7, 1, 3, 5, 7, 9, 11, 13, 15}

	levels, orders = grid.ComputeChildren(levels, orders)

	assertEqual(levels, childLevels, t)
	assertEqual(orders, childOrders, t)
}

func TestComputeChildren2D(t *testing.T) {
	grid := New(2)

	levels := []uint8{
		0, 0,
		0, 1,
		0, 1,
		1, 0,
		1, 0,
		0, 2,
		0, 2,
		1, 1,
		1, 1,
		1, 1,
		1, 1,
		2, 0,
		2, 0,
	}

	orders := []uint32{
		0, 0,
		0, 0,
		0, 2,
		0, 0,
		2, 0,
		0, 1,
		0, 3,
		0, 0,
		0, 2,
		2, 0,
		2, 2,
		1, 0,
		3, 0,
	}

	childLevels := []uint8{
		1, 0,
		1, 0,
		0, 1,
		0, 1,
		1, 1,
		1, 1,
		0, 2,
		1, 1,
		1, 1,
		0, 2,
		2, 0,
		2, 0,
		1, 2,
		1, 2,
		0, 3,
		0, 3,
		1, 2,
		1, 2,
		0, 3,
		0, 3,
		2, 1,
		2, 1,
		2, 1,
		2, 1,
		3, 0,
		3, 0,
		3, 0,
		3, 0,
	}

	childOrders := []uint32{
		0, 0,
		2, 0,
		0, 0,
		0, 2,
		0, 0,
		2, 0,
		0, 1,
		0, 2,
		2, 2,
		0, 3,
		1, 0,
		3, 0,
		0, 1,
		2, 1,
		0, 1,
		0, 3,
		0, 3,
		2, 3,
		0, 5,
		0, 7,
		1, 0,
		1, 2,
		3, 0,
		3, 2,
		1, 0,
		3, 0,
		5, 0,
		7, 0,
	}

	levels, orders = grid.ComputeChildren(levels, orders)

	assertEqual(levels, childLevels, t)
	assertEqual(orders, childOrders, t)
}
