package linhat

import (
	"math/rand"
	"testing"

	"github.com/ready-steady/numeric/grid/newcot"
	"github.com/ready-steady/support/assert"
)

func TestOpenEvaluate(t *testing.T) {
	basis := NewOpen(1)

	points := []float64{
		-1.00,
		0.00, 0.04, 0.08, 0.12, 0.16, 0.20, 0.24, 0.28, 0.32, 0.36, 0.40, 0.44, 0.48,
		0.52, 0.56, 0.60, 0.64, 0.68, 0.72, 0.76, 0.80, 0.84, 0.88, 0.92, 0.96, 1.00,
		2.00,
	}

	cases := []struct {
		level  uint32
		order  uint32
		values []float64
	}{
		{0, 0, []float64{
			0,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
			0,
		}},
		{1, 0, []float64{
			0,
			2.00, 1.84, 1.68, 1.52, 1.36, 1.20, 1.04, 0.88, 0.72, 0.56, 0.40, 0.24, 0.08,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0,
		}},
		{1, 2, []float64{
			0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0.08, 0.24, 0.40, 0.56, 0.72, 0.88, 1.04, 1.20, 1.36, 1.52, 1.68, 1.84, 2.00,
			0,
		}},
		{2, 0, []float64{
			0,
			2.00, 1.68, 1.36, 1.04, 0.72, 0.40, 0.08, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0,
		}},
		{2, 2, []float64{
			0,
			0, 0, 0, 0, 0, 0, 0, 0.24, 0.56, 0.88, 0.80, 0.48, 0.16,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0,
		}},
		{2, 4, []float64{
			0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0.16, 0.48, 0.80, 0.88, 0.56, 0.24, 0, 0, 0, 0, 0, 0, 0,
			0,
		}},
		{2, 6, []float64{
			0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0.08, 0.40, 0.72, 1.04, 1.36, 1.68, 2.00,
			0,
		}},
		{3, 0, []float64{
			0,
			2.00, 1.36, 0.72, 0.08, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0,
		}},
		{3, 2, []float64{
			0,
			0, 0, 0, 0, 0.56, 0.80, 0.16, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0,
		}},
		{3, 4, []float64{
			0,
			0, 0, 0, 0, 0, 0, 0, 0.48, 0.88, 0.24, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0,
		}},
		{3, 6, []float64{
			0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0.40, 0.96, 0.32,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0,
		}},
		{3, 8, []float64{
			0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0.32, 0.96, 0.40, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0,
		}},
		{3, 10, []float64{
			0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0.24, 0.88, 0.48, 0, 0, 0, 0, 0, 0, 0,
			0,
		}},
		{3, 12, []float64{
			0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0.16, 0.80, 0.56, 0, 0, 0, 0,
			0,
		}},
		{3, 14, []float64{
			0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0.08, 0.72, 1.36, 2.00,
			0,
		}},
	}

	values := make([]float64, len(points))

	for i := range cases {
		for j := range values {
			pair := uint64(cases[i].level) | uint64(cases[i].order)<<32
			values[j] = basis.Evaluate([]uint64{pair}, []float64{points[j]})
		}
		assert.AlmostEqual(values, cases[i].values, t)
	}
}

func BenchmarkOpenEvaluate(b *testing.B) {
	const (
		inputs = 20
		level  = 5
	)

	grid := newcot.NewOpen(inputs)
	basis := NewOpen(inputs)

	indices := make([]uint64, inputs)
	levelIndices := make([]uint64, inputs)
	for i := 1; i <= level; i++ {
		levelIndices = grid.ComputeChildren(levelIndices)
		indices = append(indices, levelIndices...)
	}

	count := len(indices) / inputs

	generator := rand.New(rand.NewSource(0))
	points := make([]float64, count*inputs)
	for i := range points {
		points[i] = generator.Float64()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < count; j++ {
			basis.Evaluate(indices[j*inputs:], points[j*inputs:])
		}
	}
}
