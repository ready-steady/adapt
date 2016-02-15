package polynomial

import (
	"math"
	"math/rand"
	"testing"

	"github.com/ready-steady/assert"
)

type childrener interface {
	Children([]uint64) []uint64
}

func f(x float64) float64 {
	return 4.0*x*x*x - 3.0*x*x + 1.0
}

func F(x float64) float64 {
	return x*x*x*x - x*x*x + x
}

func TestIntegrate(t *testing.T) {
	const (
		a = -6.0
		b = +6.0
	)

	nodes := uint(math.Ceil((float64(3) + 1.0) / 2.0))
	value := integrate(a, b, nodes, f)

	assert.EqualWithin(value, F(b)-F(a), 1e-12, t)
}

func compose(level, order uint64) uint64 {
	return level | order<<levelSize
}

func generateIndices(nd, ns uint, grid childrener) []uint64 {
	parents := make([]uint64, nd)
	indices := append([]uint64{}, parents...)
	for uint(len(indices))/nd < ns {
		parents = grid.Children(parents)
		indices = append(indices, parents...)
	}
	return indices
}

func generatePoints(nd, ns uint) []float64 {
	points := make([]float64, nd*ns)
	for i := range points {
		points[i] = rand.Float64()
	}
	return points
}
