package linhat

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestClosedEvaluateComposite(t *testing.T) {
	basis := NewClosed(1, 1)

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
			pair := uint64(cases[i].level) | uint64(cases[i].order)<<32
			basis.EvaluateComposite([]uint64{pair}, []float64{1},
				[]float64{points[j]}, values[j:j+1])
		}
		assert.Equal(values, cases[i].values, t)
	}
}
