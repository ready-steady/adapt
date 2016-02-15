package polynomial

import (
	"testing"

	"github.com/ready-steady/assert"

	grid "github.com/ready-steady/adapt/grid/equidistant"
)

func BenchmarkClosedCompute(b *testing.B) {
	const (
		nd = 10
		ns = 10000
	)

	basis := NewClosed(nd, 1)
	indices := generateIndices(nd, ns, grid.NewClosed(nd).Children)
	points := generatePoints(nd, ns, indices, closedNode)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < ns; j++ {
			basis.Compute(indices[j*nd:(j+1)*nd], points[j*nd:(j+1)*nd])
		}
	}
}

func TestClosedCompute(t *testing.T) {
	basis := NewClosed(1, 1)

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

func TestClosedIntegrate(t *testing.T) {
	basis := NewClosed(1, 1)

	levels := []uint64{0, 1, 2, 3}
	values := []float64{1.0, 0.25, 1.0 / 2 / 2, 1.0 / 2 / 2 / 2}

	for i := range levels {
		assert.Equal(basis.Integrate([]uint64{compose(levels[i], 0)}), values[i], t)
	}
}

func TestClosedParent(t *testing.T) {
	childLevels := []uint64{1, 1, 2, 2, 3, 3, 3, 3}
	childOrders := []uint64{0, 2, 1, 3, 1, 3, 5, 7}

	parentLevels := []uint64{0, 0, 1, 1, 2, 2, 2, 2}
	parentOrders := []uint64{0, 0, 0, 2, 1, 1, 3, 3}

	for i := range childLevels {
		level, order := closedParent(childLevels[i], childOrders[i])
		assert.Equal(level, parentLevels[i], t)
		assert.Equal(order, parentOrders[i], t)
	}
}
