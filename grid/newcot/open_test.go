package newcot

import (
	"testing"

	"github.com/go-math/support/assert"
)

func TestOpenComputeNodes1D(t *testing.T) {
	grid := NewOpen(1)

	levels := []uint32{0, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 0, 2, 4, 6, 0, 2, 4, 6, 8, 10, 12, 14}
	nodes := []float64{
		0.5000, 0.2500, 0.7500, 0.1250, 0.3750,
		0.6250, 0.8750, 0.0625, 0.1875, 0.3125,
		0.4375, 0.5625, 0.6875, 0.8125, 0.9375,
	}

	assert.Equal(grid.ComputeNodes(compose(levels, orders)), nodes, t)
}

func TestOpenComputeChildren1D(t *testing.T) {
	grid := NewOpen(1)

	levels := []uint32{0, 1, 1, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3}
	orders := []uint32{0, 0, 2, 0, 2, 4, 6, 0, 2, 4, 6, 8, 10, 12, 14}
	childLevels := []uint32{
		1, 1,
		2, 2, 2, 2,
		3, 3, 3, 3, 3, 3, 3, 3,
		4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4,
	}
	childOrders := []uint32{
		0, 2,
		0, 2, 4, 6,
		0, 2, 4, 6, 8, 10, 12, 14,
		0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30,
	}

	index := grid.ComputeChildren(compose(levels, orders))

	assert.Equal(index, compose(childLevels, childOrders), t)
}
