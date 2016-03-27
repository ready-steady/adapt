package local

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestFilter(t *testing.T) {
	indices := filter([]uint64{1, 2, 3, 4, 5, 6, 7, 8}, []float64{1.0, 2.0, 3.0, 4.0}, 1, 20, 2)
	assert.Equal(indices, []uint64{1, 2, 3, 4, 5, 6, 7, 8}, t)

	indices = filter([]uint64{1, 2, 3, 4, 5, 6, 7, 8}, []float64{0.0, 2.0, 3.0, 4.0}, 4, 20, 2)
	assert.Equal(indices, []uint64{1, 2, 3, 4, 5, 6, 7, 8}, t)

	indices = filter([]uint64{1, 2, 3, 4, 5, 6, 7, 8}, []float64{0.0, 2.0, 3.0, 4.0}, 1, 20, 2)
	assert.Equal(indices, []uint64{3, 4, 5, 6, 7, 8}, t)

	indices = filter([]uint64{1, 2, 3, 4, 5, 6, 7, 8}, []float64{1.0, 2.0, 3.0, 4.0}, 1, 10, 2)
	assert.Equal(indices, []uint64{1, 2, 3, 4}, t)
}
