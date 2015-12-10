package local

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestTrackerFilter(t *testing.T) {
	test := func(tracker *tracker, indices []uint64, scores []float64, result []uint64) {
		assert.Equal(tracker.filter(indices, scores), result, t)
	}

	test(
		&tracker{ni: 2, lmin: 1, lmax: 20},
		[]uint64{1, 2, 3, 4, 5, 6, 7, 8},
		[]float64{1.0, 2.0, 3.0, 4.0},
		[]uint64{1, 2, 3, 4, 5, 6, 7, 8},
	)

	test(
		&tracker{ni: 2, lmin: 4, lmax: 20},
		[]uint64{1, 2, 3, 4, 5, 6, 7, 8},
		[]float64{0.0, 2.0, 3.0, 4.0},
		[]uint64{1, 2, 3, 4, 5, 6, 7, 8},
	)

	test(
		&tracker{ni: 2, lmin: 1, lmax: 20},
		[]uint64{1, 2, 3, 4, 5, 6, 7, 8},
		[]float64{0.0, 2.0, 3.0, 4.0},
		[]uint64{3, 4, 5, 6, 7, 8},
	)

	test(
		&tracker{ni: 2, lmin: 1, lmax: 10},
		[]uint64{1, 2, 3, 4, 5, 6, 7, 8},
		[]float64{1.0, 2.0, 3.0, 4.0},
		[]uint64{1, 2, 3, 4},
	)
}
