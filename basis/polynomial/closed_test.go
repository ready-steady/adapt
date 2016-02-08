package polynomial

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestCompute(t *testing.T) {
	basis := New(1, 1)

	compute := func(level, order uint64, point float64) float64 {
		return basis.Compute([]uint64{compose(level, order)}, []float64{point})
	}

	points := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

	cases := []struct {
		level  uint64
		order  uint64
		values []float64
	}{
		{0, 0, []float64{1.0, 1.0, 1.0, 1.0, 1.0}},
		{1, 0, []float64{1.0, 0.5, 0.0, 0.0, 0.0}},
		{1, 2, []float64{0.0, 0.0, 0.0, 0.5, 1.0}},
		{2, 1, []float64{0.0, 1.0, 0.0, 0.0, 0.0}},
		{2, 3, []float64{0.0, 0.0, 0.0, 1.0, 0.0}},
	}

	values := make([]float64, len(points))

	for i := range cases {
		for j := range values {
			values[j] = compute(cases[i].level, cases[i].order, points[j])
		}
		assert.Equal(values, cases[i].values, t)
	}
}
