package linhat

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestClosedEvaluateComposite(t *testing.T) {
	basis := NewClosed(1, 1)

	evaluate := func(level, order uint32, point float64) float64 {
		pair := uint64(level) | uint64(order)<<32
		return basis.EvaluateComposite([]uint64{pair}, []float64{1}, []float64{point})[0]
	}

	points := []float64{-1, 0, 0.25, 0.5, 0.75, 1, 2}

	cases := []struct {
		level  uint32
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
			values[j] = evaluate(cases[i].level, cases[i].order, points[j])
		}
		assert.Equal(values, cases[i].values, t)
	}
}
