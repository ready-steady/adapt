package ode

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestDormandPrinceComputeToy(t *testing.T) {
	fixture := &fixtureToy

	evaluations := []float64{
		0.0000000000000000e+00, 2.0000000000000004e-02, 2.9999999999999999e-02,
		8.0000000000000016e-02, 8.8888888888888892e-02, 1.0000000000000001e-01,
		1.0000000000000001e-01, 1.2000000000000001e-01, 1.3000000000000000e-01,
		1.8000000000000002e-01, 1.8888888888888888e-01, 2.0000000000000001e-01,
		2.0000000000000001e-01, 2.2000000000000003e-01, 2.3000000000000001e-01,
		2.8000000000000003e-01, 2.8888888888888892e-01, 3.0000000000000004e-01,
		3.0000000000000004e-01, 3.2000000000000006e-01, 3.3000000000000007e-01,
		3.8000000000000006e-01, 3.8888888888888895e-01, 4.0000000000000002e-01,
		4.0000000000000002e-01, 4.2000000000000004e-01, 4.3000000000000005e-01,
		4.8000000000000004e-01, 4.8888888888888893e-01, 5.0000000000000000e-01,
		5.0000000000000000e-01, 5.2000000000000002e-01, 5.3000000000000003e-01,
		5.8000000000000007e-01, 5.8888888888888891e-01, 5.9999999999999998e-01,
		5.9999999999999998e-01, 6.2000000000000000e-01, 6.3000000000000000e-01,
		6.7999999999999994e-01, 6.8888888888888888e-01, 6.9999999999999996e-01,
		6.9999999999999996e-01, 7.1999999999999997e-01, 7.2999999999999998e-01,
		7.8000000000000003e-01, 7.8888888888888886e-01, 7.9999999999999993e-01,
		7.9999999999999993e-01, 8.1999999999999995e-01, 8.2999999999999996e-01,
		8.7999999999999989e-01, 8.8888888888888884e-01, 8.9999999999999991e-01,
		8.9999999999999991e-01, 9.1999999999999993e-01, 9.2999999999999994e-01,
		9.7999999999999998e-01, 9.8888888888888893e-01, 1.0000000000000000e+00,
		1.0000000000000000e+00,
	}

	derivative := func(x float64, y, f []float64) {
		assert.Equal(x, evaluations[0], t)
		evaluations = evaluations[1:]
		fixture.derivative(x, y, f)
	}

	integrator, err := NewDormandPrince(fixture.configure())
	assert.Success(err, t)

	values, stats, err := integrator.Compute(derivative, fixture.points, fixture.initial)
	assert.Success(err, t)
	assert.EqualWithin(values, fixture.values, 1e-15, t)
	assert.Equal(*stats, Stats{Evaluations: 61, Rejections: 0, Steps: 10}, t)
}

func TestDormandPrinceComputeNonstiff(t *testing.T) {
	fixture := &fixtureNonstiff

	integrator, err := NewDormandPrince(fixture.configure())
	assert.Success(err, t)

	values, stats, err := integrator.Compute(fixture.derivative, fixture.points, fixture.initial)
	assert.Success(err, t)
	assert.EqualWithin(values, fixture.values, 1e-14, t)
	assert.Equal(*stats, Stats{Evaluations: 151, Rejections: 3, Steps: 22}, t)
}

func TestDormandPrinceComputeStiff(t *testing.T) {
	fixture := &fixtureStiff

	integrator, err := NewDormandPrince(fixture.configure())
	assert.Success(err, t)

	values, stats, err := integrator.Compute(fixture.derivative, fixture.points, fixture.initial)
	assert.Success(err, t)
	assert.EqualWithin(values, fixture.values, 3e-13, t)
	assert.Equal(*stats, Stats{Evaluations: 20179, Rejections: 323, Steps: 3040}, t)
}

func BenchmarkDormandPrinceComputeNonstiff(b *testing.B) {
	fixture := &fixtureNonstiff
	integrator, _ := NewDormandPrince(fixture.configure())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		integrator.Compute(fixture.derivative, fixture.points, fixture.initial)
	}
}

func BenchmarkDormandPrinceComputeStiff(b *testing.B) {
	fixture := &fixtureStiff
	integrator, _ := NewDormandPrince(fixture.configure())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		integrator.Compute(fixture.derivative, fixture.points, fixture.initial)
	}
}
