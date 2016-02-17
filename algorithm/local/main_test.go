package local

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestStep(t *testing.T) {
	fixture := &fixtureStep
	interpolator, target := prepare(fixture)

	surrogate := interpolator.Compute(target)
	assert.Equal(surrogate, fixture.surrogate, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.Equal(values, fixture.values, t)

	integral := interpolator.Integrate(surrogate)
	assert.Equal(integral, fixture.integral, t)
}

func TestHat(t *testing.T) {
	fixture := &fixtureHat
	interpolator, target := prepare(fixture)

	surrogate := interpolator.Compute(target)
	assert.Equal(surrogate, fixture.surrogate, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 1e-15, t)

	integral := interpolator.Integrate(surrogate)
	assert.Equal(integral, fixture.integral, t)
}

func TestCube(t *testing.T) {
	fixture := &fixtureCube
	interpolator, target := prepare(fixture)

	surrogate := interpolator.Compute(target)
	assert.Equal(surrogate.Level, fixture.surrogate.Level, t)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 2e-15, t)

	integral := interpolator.Integrate(surrogate)
	assert.Equal(integral, fixture.integral, t)
}

func TestBox(t *testing.T) {
	fixture := &fixtureBox
	interpolator, target := prepare(fixture)

	surrogate := interpolator.Compute(target)
	assert.Equal(surrogate, fixture.surrogate, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 1e-15, t)

	integral := interpolator.Integrate(surrogate)
	assert.Equal(integral, fixture.integral, t)
}

func TestKraichnanOrszag(t *testing.T) {
	fixture := &fixtureKraichnanOrszag
	interpolator, target := prepare(fixture)

	surrogate := interpolator.Compute(target)
	assert.Equal(surrogate.Level, fixture.surrogate.Level, t)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 6e-14, t)

	integral := interpolator.Integrate(surrogate)
	assert.EqualWithin(integral, fixture.integral, 2e-14, t)
}

func TestParabola(t *testing.T) {
	fixture := &fixtureParabola
	interpolator, target := prepare(fixture)

	surrogate := interpolator.Compute(target)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 1e-6, t)

	integral := interpolator.Integrate(surrogate)
	assert.EqualWithin(integral, fixture.integral, 1e-6, t)
}
