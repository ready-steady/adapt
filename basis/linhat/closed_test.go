package linhat

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestClosedCompute(t *testing.T) {
	basis := NewClosed(1)

	compute := func(level, order uint32, point float64) float64 {
		return basis.Compute([]uint64{compose(level, order)}, []float64{point})
	}

	points := []float64{0, 0.25, 0.5, 0.75, 1}

	cases := []struct {
		level  uint32
		order  uint32
		values []float64
	}{
		{0, 0, []float64{1, 1.0, 1, 1.0, 1}},
		{1, 0, []float64{1, 0.5, 0, 0.0, 0}},
		{1, 2, []float64{0, 0.0, 0, 0.5, 1}},
		{2, 1, []float64{0, 1.0, 0, 0.0, 0}},
		{2, 3, []float64{0, 0.0, 0, 1.0, 0}},
	}

	values := make([]float64, len(points))

	for i := range cases {
		for j := range values {
			values[j] = compute(cases[i].level, cases[i].order, points[j])
		}
		assert.Equal(values, cases[i].values, t)
	}
}

func TestClosedIntegrate(t *testing.T) {
	basis := NewClosed(1)

	levels := []uint32{0, 1, 2, 3}
	values := []float64{1, 0.25, 1.0 / 2 / 2, 1.0 / 2 / 2 / 2}

	for i := range levels {
		assert.Equal(basis.Integrate([]uint64{compose(levels[i], uint32(0))}), values[i], t)
	}
}
