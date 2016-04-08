package local

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestFilter(t *testing.T) {
	const (
		εl = 0.0
		ni = 2
	)

	var indices []uint64

	indices = []uint64{1, 2, 3, 4, 5, 6, 7, 8}
	indices = filter(indices, []float64{1.0, 2.0, 3.0, 4.0}, 1, 20, εl, ni)
	assert.Equal(indices, []uint64{1, 2, 3, 4, 5, 6, 7, 8}, t)

	indices = []uint64{1, 2, 3, 4, 5, 6, 7, 8}
	indices = filter(indices, []float64{0.0, 2.0, 3.0, 4.0}, 4, 20, εl, ni)
	assert.Equal(indices, []uint64{1, 2, 3, 4, 5, 6, 7, 8}, t)

	indices = []uint64{1, 2, 3, 4, 5, 6, 7, 8}
	indices = filter(indices, []float64{0.0, 2.0, 3.0, 4.0}, 1, 20, εl, ni)
	assert.Equal(indices, []uint64{3, 4, 5, 6, 7, 8}, t)

	indices = []uint64{1, 2, 3, 4, 5, 6, 7, 8}
	indices = filter(indices, []float64{1.0, 2.0, 3.0, 4.0}, 1, 10, εl, ni)
	assert.Equal(indices, []uint64{1, 2, 3, 4}, t)
}
