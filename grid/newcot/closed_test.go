package newcot

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestClosedCompute1D(t *testing.T) {
	grid := NewClosed(1)

	levels := []uint32{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 1, 3, 1, 3, 5, 7}
	nodes := []float64{0.5, 0, 1, 0.25, 0.75, 0.125, 0.375, 0.625, 0.875}

	assert.Equal(grid.Compute(compose(levels, orders)), nodes, t)
}

func TestClosedCompute2D(t *testing.T) {
	grid := NewClosed(2)

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

	assert.Equal(grid.Compute(compose(levels, orders)), nodes, t)
}

func TestClosedRefine1D(t *testing.T) {
	grid := NewClosed(1)

	levels := []uint32{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 1, 3, 1, 3, 5, 7}
	childLevels := []uint32{1, 1, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4}
	childOrders := []uint32{0, 2, 1, 3, 1, 3, 5, 7, 1, 3, 5, 7, 9, 11, 13, 15}

	indices := grid.Refine(compose(levels, orders))

	assert.Equal(indices, compose(childLevels, childOrders), t)
}

func TestClosedRefine2D(t *testing.T) {
	grid := NewClosed(2)

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
		1, 1,
		1, 1,
		2, 0,
		1, 1,
		1, 1,
		1, 2,
		1, 2,
		0, 3,
		0, 3,
		1, 2,
		1, 2,
		0, 3,
		0, 3,
		2, 1,
		1, 2,
		2, 1,
		1, 2,
		2, 1,
		1, 2,
		2, 1,
		1, 2,
		3, 0,
		3, 0,
		2, 1,
		2, 1,
		3, 0,
		3, 0,
		2, 1,
		2, 1,
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
		0, 0,
		0, 2,
		3, 0,
		2, 0,
		2, 2,
		0, 1,
		2, 1,
		0, 1,
		0, 3,
		0, 3,
		2, 3,
		0, 5,
		0, 7,
		1, 0,
		0, 1,
		1, 2,
		0, 3,
		3, 0,
		2, 1,
		3, 2,
		2, 3,
		1, 0,
		3, 0,
		1, 0,
		1, 2,
		5, 0,
		7, 0,
		3, 0,
		3, 2,
	}

	indices := grid.Refine(compose(levels, orders))

	assert.Equal(indices, compose(childLevels, childOrders), t)
}

func TestClosedParent(t *testing.T) {
	grid := NewClosed(1)

	children := compose(
		[]uint32{0, 1, 1, 2, 2, 3, 3, 3, 3},
		[]uint32{0, 0, 2, 1, 3, 1, 3, 5, 7},
	)

	parents := compose(
		[]uint32{0, 0, 0, 1, 1, 2, 2, 2, 2},
		[]uint32{0, 0, 0, 0, 2, 1, 1, 3, 3},
	)

	for i := range children {
		grid.Parent(children[i:i+1], 0)
	}

	assert.Equal(children, parents, t)
}

func BenchmarkClosedRefine(b *testing.B) {
	const (
		dimensions  = 20
		targetLevel = 3
	)

	grid := NewClosed(dimensions)

	// Level 0
	indices := make([]uint64, dimensions)

	// Level 1, 2, …, (targetLevel - 1)
	for i := 1; i < targetLevel; i++ {
		indices = grid.Refine(indices)
	}

	b.ResetTimer()

	// Level targetLevel
	for i := 0; i < b.N; i++ {
		grid.Refine(indices)
	}
}
