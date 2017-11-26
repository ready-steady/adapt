package polynomial

import (
	"math"
	"testing"

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

	assert.Close(value, F(b)-F(a), 1e-12, t)
}
