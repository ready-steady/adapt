package equidistant

import (
	"testing"

	"github.com/ready-steady/adapt/internal"
	"github.com/ready-steady/assert"
)

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

func TestClosedCompute1D(t *testing.T) {
	grid := NewClosed(1)

	levels := []uint64{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint64{0, 0, 2, 1, 3, 1, 3, 5, 7}
	nodes := []float64{0.5, 0.0, 1.0, 0.25, 0.75, 0.125, 0.375, 0.625, 0.875}

	assert.Equal(grid.Compute(internal.Compose(levels, orders)), nodes, t)
}

func TestClosedCompute2D(t *testing.T) {
	grid := NewClosed(2)

	levels := []uint64{
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

	orders := []uint64{
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

	assert.Equal(grid.Compute(internal.Compose(levels, orders)), nodes, t)
}

func TestClosedIndex1D(t *testing.T) {
	cases := []struct {
		level  uint64
		levels []uint64
		orders []uint64
	}{
		{
			level:  0,
			levels: []uint64{0},
			orders: []uint64{0},
		},
		{
			level:  1,
			levels: []uint64{1, 1},
			orders: []uint64{0, 2},
		},
		{
			level:  2,
			levels: []uint64{2, 2},
			orders: []uint64{1, 3},
		},
		{
			level:  3,
			levels: []uint64{3, 3, 3, 3},
			orders: []uint64{1, 3, 5, 7},
		},
	}

	for _, c := range cases {
		assert.Equal(closedIndex(c.level), internal.Compose(c.levels, c.orders), t)
	}
}

func TestClosedIndex2D(t *testing.T) {
	grid := NewClosed(2)

	cases := []struct {
		level  []uint64
		levels []uint64
		orders []uint64
	}{
		{
			level: []uint64{0, 0},
			levels: []uint64{
				0, 0,
			},
			orders: []uint64{
				0, 0,
			},
		},
		{
			level: []uint64{0, 1},
			levels: []uint64{
				0, 1,
				0, 1,
			},
			orders: []uint64{
				0, 0,
				0, 2,
			},
		},
		{
			level: []uint64{1, 2},
			levels: []uint64{
				1, 2,
				1, 2,
				1, 2,
				1, 2,
			},
			orders: []uint64{
				0, 1,
				2, 1,
				0, 3,
				2, 3,
			},
		},
		{
			level: []uint64{2, 3},
			levels: []uint64{
				2, 3,
				2, 3,
				2, 3,
				2, 3,
				2, 3,
				2, 3,
				2, 3,
				2, 3,
			},
			orders: []uint64{
				1, 1,
				3, 1,
				1, 3,
				3, 3,
				1, 5,
				3, 5,
				1, 7,
				3, 7,
			},
		},
	}

	for _, c := range cases {
		assert.Equal(grid.Index(c.level), internal.Compose(c.levels, c.orders), t)
	}
}

func TestClosedParent(t *testing.T) {
	grid := NewClosed(1)

	childLevels := []uint64{1, 1, 2, 2, 3, 3, 3, 3}
	childOrders := []uint64{0, 2, 1, 3, 1, 3, 5, 7}

	parentLevels := []uint64{0, 0, 1, 1, 2, 2, 2, 2}
	parentOrders := []uint64{0, 0, 0, 2, 1, 1, 3, 3}

	for i := range childLevels {
		level, order := grid.Parent(childLevels[i], childOrders[i])
		assert.Equal(level, parentLevels[i], t)
		assert.Equal(order, parentOrders[i], t)
	}
}

func TestClosedRefine1D(t *testing.T) {
	grid := NewClosed(1)

	levels := []uint64{0, 1, 1, 2, 2, 3, 3, 3, 3}
	orders := []uint64{0, 0, 2, 1, 3, 1, 3, 5, 7}
	childLevels := []uint64{1, 1, 2, 2, 3, 3, 3, 3, 4, 4, 4, 4, 4, 4, 4, 4}
	childOrders := []uint64{0, 2, 1, 3, 1, 3, 5, 7, 1, 3, 5, 7, 9, 11, 13, 15}

	indices := grid.Refine(internal.Compose(levels, orders))

	assert.Equal(indices, internal.Compose(childLevels, childOrders), t)
}

func TestClosedRefine2D(t *testing.T) {
	grid := NewClosed(2)

	levels := []uint64{
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

	orders := []uint64{
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

	childLevels := []uint64{
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

	childOrders := []uint64{
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

	indices := grid.Refine(internal.Compose(levels, orders))

	assert.Equal(indices, internal.Compose(childLevels, childOrders), t)
}
