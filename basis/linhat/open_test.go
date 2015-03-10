package linhat

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestOpenCompute(t *testing.T) {
	basis := NewOpen(1)

	compute := func(level, order uint32, point float64) float64 {
		return basis.Compute([]uint64{compose(level, order)}, []float64{point})
	}

	points := []float64{
		0.00, 0.04, 0.08, 0.12, 0.16, 0.20, 0.24, 0.28, 0.32, 0.36, 0.40, 0.44, 0.48,
		0.52, 0.56, 0.60, 0.64, 0.68, 0.72, 0.76, 0.80, 0.84, 0.88, 0.92, 0.96, 1.00,
	}

	cases := []struct {
		level  uint32
		order  uint32
		values []float64
	}{
		{0, 0, []float64{
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		}},
		{1, 0, []float64{
			2.00, 1.84, 1.68, 1.52, 1.36, 1.20, 1.04, 0.88, 0.72, 0.56, 0.40, 0.24, 0.08,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}},
		{1, 2, []float64{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0.08, 0.24, 0.40, 0.56, 0.72, 0.88, 1.04, 1.20, 1.36, 1.52, 1.68, 1.84, 2.00,
		}},
		{2, 0, []float64{
			2.00, 1.68, 1.36, 1.04, 0.72, 0.40, 0.08, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}},
		{2, 2, []float64{
			0, 0, 0, 0, 0, 0, 0, 0.24, 0.56, 0.88, 0.80, 0.48, 0.16,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}},
		{2, 4, []float64{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0.16, 0.48, 0.80, 0.88, 0.56, 0.24, 0, 0, 0, 0, 0, 0, 0,
		}},
		{2, 6, []float64{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0.08, 0.40, 0.72, 1.04, 1.36, 1.68, 2.00,
		}},
		{3, 0, []float64{
			2.00, 1.36, 0.72, 0.08, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}},
		{3, 2, []float64{
			0, 0, 0, 0, 0.56, 0.80, 0.16, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}},
		{3, 4, []float64{
			0, 0, 0, 0, 0, 0, 0, 0.48, 0.88, 0.24, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}},
		{3, 6, []float64{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0.40, 0.96, 0.32,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}},
		{3, 8, []float64{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0.32, 0.96, 0.40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		}},
		{3, 10, []float64{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0.24, 0.88, 0.48, 0, 0, 0, 0, 0, 0, 0,
		}},
		{3, 12, []float64{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0.16, 0.80, 0.56, 0, 0, 0, 0,
		}},
		{3, 14, []float64{
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0.08, 0.72, 1.36, 2.00,
		}},
	}

	values := make([]float64, len(points))

	for i := range cases {
		for j := range values {
			values[j] = compute(cases[i].level, cases[i].order, points[j])
		}
		assert.EqualWithin(values, cases[i].values, 1e-15, t)
	}
}

func TestOpenIntegrate(t *testing.T) {
	basis := NewOpen(1)

	levels := []uint32{0, 1, 1, 2, 2, 2, 2}
	orders := []uint32{0, 0, 2, 0, 2, 4, 6}
	values := []float64{1, 0.5, 0.5, 0.25, 0.125, 0.125, 0.25}

	for i := range levels {
		assert.Equal(basis.Integrate([]uint64{compose(levels[i], orders[i])}), values[i], t)
	}
}
