package polynomial

import (
	"math"
	"math/rand"

	"github.com/ready-steady/adapt/internal"
)

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
	locate func(level, order uint64) (float64, float64, uint64)) []float64 {

	levels, orders := internal.Decompose(indices)
	points := make([]float64, nd*ns)
	for i := range points {
		x, h, _ := locate(levels[i], orders[i])
		a, b := math.Max(0.0, x-h), math.Min(1.0, x+h)
		points[i] = a + (b-a)*rand.Float64()
	}
	return points
}
