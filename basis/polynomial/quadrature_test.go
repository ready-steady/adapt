package polynomial

import (
	"math"
	"math/rand"
	"testing"

	"github.com/ready-steady/adapt/internal"
	"github.com/ready-steady/assert"
)

func f(x float64) float64 {
	return 4.0*x*x*x - 3.0*x*x + 1.0
}

func F(x float64) float64 {
	return x*x*x*x - x*x*x + x
}

func TestQuadrature(t *testing.T) {
	const (
		a = -6.0
		b = +6.0
	)

	nodes := uint(math.Ceil((float64(3) + 1.0) / 2.0))
	value := quadrature(a, b, nodes, f)

	assert.EqualWithin(value, F(b)-F(a), 1e-12, t)
}

func generateIndices(nd, ns uint, children func([]uint64) []uint64) []uint64 {
	parents := make([]uint64, nd)
	indices := append([]uint64{}, parents...)
	for uint(len(indices))/nd < ns {
		parents = children(parents)
		indices = append(indices, parents...)
	}
	return indices
}

func generatePoints(nd, ns uint, indices []uint64,
	locate func(level, order uint64) (float64, float64)) []float64 {

	levels, orders := internal.Decompose(indices)
	points := make([]float64, nd*ns)
	for i := range points {
		x, h := locate(levels[i], orders[i])
		a, b := math.Max(0.0, x-h), math.Min(1.0, x+h)
		points[i] = a + (b-a)*rand.Float64()
	}
	return points
}
