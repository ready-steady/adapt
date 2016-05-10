package equidistant

import (
	"testing"

	"github.com/ready-steady/adapt/internal"
	"github.com/ready-steady/assert"
)

func TestOpenCompute1D(t *testing.T) {
	grid := NewOpen(1)

	levels := []uint64{0, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3}
	orders := []uint64{0, 0, 2, 0, 2, 4, 6, 0, 2, 4, 6, 8, 10, 12, 14}
	nodes := []float64{
		0.5000, 0.2500, 0.7500, 0.1250, 0.3750,
		0.6250, 0.8750, 0.0625, 0.1875, 0.3125,
		0.4375, 0.5625, 0.6875, 0.8125, 0.9375,
	}

	assert.Equal(grid.Compute(internal.Compose(levels, orders)), nodes, t)
}

func TestOpenCompute2D(t *testing.T) {
	grid := NewOpen(2)

	levels := []uint64{
		0, 0,

		0, 1,
		0, 1,
		0, 2,
		0, 2,
		0, 2,
		0, 2,

		1, 0,
		1, 0,

		1, 1,
		1, 1,
		1, 1,
		1, 1,

		1, 2,
		1, 2,
		1, 2,
		1, 2,
		1, 2,
		1, 2,
		1, 2,
		1, 2,

		2, 0,
		2, 0,
		2, 0,
		2, 0,

		2, 1,
		2, 1,
		2, 1,
		2, 1,
		2, 1,
		2, 1,
		2, 1,
		2, 1,

		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
		2, 2,
	}

	orders := []uint64{
		0, 0,

		0, 0,
		0, 2,
		0, 0,
		0, 2,
		0, 4,
		0, 6,

		0, 0,
		2, 0,

		0, 0,
		0, 2,
		2, 0,
		2, 2,

		0, 0,
		0, 2,
		0, 4,
		0, 6,
		2, 0,
		2, 2,
		2, 4,
		2, 6,

		0, 0,
		2, 0,
		4, 0,
		6, 0,

		0, 0,
		0, 2,
		2, 0,
		2, 2,
		4, 0,
		4, 2,
		6, 0,
		6, 2,

		0, 0,
		0, 2,
		0, 4,
		0, 6,
		2, 0,
		2, 2,
		2, 4,
		2, 6,
		4, 0,
		4, 2,
		4, 4,
		4, 6,
		6, 0,
		6, 2,
		6, 4,
		6, 6,
	}

	nodes := []float64{
		0.500, 0.500,

		0.500, 0.250,
		0.500, 0.750,
		0.500, 0.125,
		0.500, 0.375,
		0.500, 0.625,
		0.500, 0.875,

		0.250, 0.500,
		0.750, 0.500,

		0.250, 0.250,
		0.250, 0.750,
		0.750, 0.250,
		0.750, 0.750,

		0.250, 0.125,
		0.250, 0.375,
		0.250, 0.625,
		0.250, 0.875,
		0.750, 0.125,
		0.750, 0.375,
		0.750, 0.625,
		0.750, 0.875,

		0.125, 0.500,
		0.375, 0.500,
		0.625, 0.500,
		0.875, 0.500,

		0.125, 0.250,
		0.125, 0.750,
		0.375, 0.250,
		0.375, 0.750,
		0.625, 0.250,
		0.625, 0.750,
		0.875, 0.250,
		0.875, 0.750,

		0.125, 0.125,
		0.125, 0.375,
		0.125, 0.625,
		0.125, 0.875,
		0.375, 0.125,
		0.375, 0.375,
		0.375, 0.625,
		0.375, 0.875,
		0.625, 0.125,
		0.625, 0.375,
		0.625, 0.625,
		0.625, 0.875,
		0.875, 0.125,
		0.875, 0.375,
		0.875, 0.625,
		0.875, 0.875,
	}

	assert.Equal(grid.Compute(internal.Compose(levels, orders)), nodes, t)
}

func TestOpenIndex1D(t *testing.T) {
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
			orders: []uint64{1, 3},
		},
		{
			level:  2,
			levels: []uint64{2, 2, 2, 2},
			orders: []uint64{1, 3, 5, 7},
		},
		{
			level:  3,
			levels: []uint64{3, 3, 3, 3, 3, 3, 3, 3},
			orders: []uint64{1, 3, 5, 7, 9, 11, 13, 15},
		},
	}

	for _, c := range cases {
		assert.Equal(openIndex(c.level), internal.Compose(c.levels, c.orders), t)
	}
}

func TestOpenParent(t *testing.T) {
	grid := NewOpen(1)

	childLevels := []uint64{1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3}
	childOrders := []uint64{0, 2, 0, 2, 4, 6, 0, 2, 4, 6, 8, 10, 12, 14}

	parentLevels := []uint64{0, 0, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2, 2}
	parentOrders := []uint64{0, 0, 0, 0, 2, 2, 0, 0, 2, 2, 4, 4, 6, 6}

	for i := range childLevels {
		level, order := grid.Parent(childLevels[i], childOrders[i])
		assert.Equal(level, parentLevels[i], t)
		assert.Equal(order, parentOrders[i], t)
	}
}

func TestOpenRefine1D(t *testing.T) {
	grid := NewOpen(1)

	levels := []uint64{0, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3}
	orders := []uint64{0, 0, 2, 0, 2, 4, 6, 0, 2, 4, 6, 8, 10, 12, 14}
	childLevels := []uint64{
		1, 1,
		2, 2, 2, 2,
		3, 3, 3, 3, 3, 3, 3, 3,
		4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	}
	childOrders := []uint64{
		0, 2,
		0, 2, 4, 6,
		0, 2, 4, 6, 8, 10, 12, 14,
		0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30,
	}

	indices := grid.Refine(internal.Compose(levels, orders))

	assert.Equal(indices, internal.Compose(childLevels, childOrders), t)
}
