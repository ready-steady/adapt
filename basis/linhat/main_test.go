package linhat

import (
	"testing"

	"github.com/go-math/support/assert"
)

func TestEvaluate(t *testing.T) {
	basis := New(1)

	points := []float64{-1, 0, 0.25, 0.5, 0.75, 1, 2}

	cases := []struct {
		level  uint8
		order  uint32
		values []float64
	}{
		{0, 0, []float64{0, 1, 1.0, 1, 1.0, 1, 0}},
		{1, 0, []float64{0, 1, 0.5, 0, 0.0, 0, 0}},
		{1, 2, []float64{0, 0, 0.0, 0, 0.5, 1, 0}},
		{2, 1, []float64{0, 0, 1.0, 0, 0.0, 0, 0}},
		{2, 3, []float64{0, 0, 0.0, 0, 1.0, 0, 0}},
	}

	values := make([]float64, len(points))

	for i := range cases {
		for j := range values {
			pair := uint64(cases[i].level)<<32 | uint64(cases[i].order)
			values[j] = basis.Evaluate([]uint64{pair}, []float64{points[j]})
		}
		assert.Equal(values, cases[i].values, t)
	}
}
