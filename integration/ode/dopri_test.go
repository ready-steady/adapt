package ode

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestDormandPrinceCompute1D(t *testing.T) {
	fixture := &fixture1D

	evaluations := fixture.evaluations

	derivative := func(x float64, y, f []float64) {
		assert.Equal(x, evaluations[0], t)
		evaluations = evaluations[1:]
		fixture.derivative(x, y, f)
	}

	integrator, err := NewDormandPrince(fixture.configure())
	assert.Success(err, t)

	values, err := integrator.Compute(derivative, fixture.points, fixture.initial)
	assert.Success(err, t)
	assert.EqualWithin(values, fixture.values, 1e-15, t)
}

func TestDormandPrinceCompute3D(t *testing.T) {
	fixture := &fixture3D

	integrator, err := NewDormandPrince(fixture.configure())
	assert.Success(err, t)

	values, err := integrator.Compute(fixture.derivative, fixture.points, fixture.initial)
	assert.Success(err, t)
	assert.EqualWithin(values, fixture.values, 1e-14, t)
}

func BenchmarkDormandPrinceCompute(b *testing.B) {
	fixture := &fixture3D
	integrator, _ := NewDormandPrince(fixture.configure())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		integrator.Compute(fixture.derivative, fixture.points, fixture.initial)
	}
}
