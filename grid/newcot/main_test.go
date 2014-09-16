package newcot

import (
	"testing"

	"github.com/go-math/support/assert"
)

func TestComputeNodes1D(t *testing.T) {
	grid := New(1)

	levels := []uint32{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 1, 3, 1, 3, 5, 7}
	nodes := []float64{0.5, 0, 1, 0.25, 0.75, 0.125, 0.375, 0.625, 0.875}

	assert.Equal(grid.ComputeNodes(compose(levels, orders)), nodes, t)
}

func TestComputeNodes2D(t *testing.T) {
	grid := New(2)

	levels := []uint32{
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

	assert.Equal(grid.ComputeNodes(compose(levels, orders)), nodes, t)
}

func TestComputeChildren1D(t *testing.T) {
	grid := New(1)

	levels := []uint32{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 1, 3, 1, 3, 5, 7}
	childLevels := []uint32{1, 1, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4}
	childOrders := []uint32{0, 2, 1, 3, 1, 3, 5, 7, 1, 3, 5, 7, 9, 11, 13, 15}

	index := grid.ComputeChildren(compose(levels, orders))

	assert.Equal(index, compose(childLevels, childOrders), t)
}

func TestComputeChildren2D(t *testing.T) {
	grid := New(2)

	levels := []uint32{
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

	childLevels := []uint32{
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

	index := grid.ComputeChildren(compose(levels, orders))

	assert.Equal(index, compose(childLevels, childOrders), t)
}

func BenchmarkComputeChildren(b *testing.B) {
	const (
		inputs      = 20
		targetLevel = 3
	)

	grid := New(inputs)

	// Level 0
	index := make([]uint64, inputs)

	// Level 1, 2, …, (targetLevel - 1)
	for i := 1; i < targetLevel; i++ {
		index = grid.ComputeChildren(index)
	}

	b.ResetTimer()

	// Level targetLevel
	for i := 0; i < b.N; i++ {
		grid.ComputeChildren(index)
	}
}

func compose(levels []uint32, orders []uint32) []uint64 {
	index := make([]uint64, len(levels))

	for i := range levels {
		index[i] = uint64(levels[i]) | uint64(orders[i])<<32
	}

	return index
}
