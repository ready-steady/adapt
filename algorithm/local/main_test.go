package local

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestStep(t *testing.T) {
	fixture := &fixtureStep
	interpolator := prepare(fixture)

	surrogate := interpolator.Compute(fixture.target)
	assert.Equal(surrogate, fixture.surrogate, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.Equal(values, fixture.values, t)
}

func TestHat(t *testing.T) {
	fixture := &fixtureHat
	interpolator := prepare(fixture)

	surrogate := interpolator.Compute(fixture.target)
	assert.Equal(surrogate, fixture.surrogate, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 1e-15, t)
}

func TestCube(t *testing.T) {
	fixture := &fixtureCube
	interpolator := prepare(fixture)

	surrogate := interpolator.Compute(fixture.target)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)
	assert.Equal(surrogate.Integral, fixture.surrogate.Integral, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 2e-15, t)
}

func TestBox(t *testing.T) {
	fixture := &fixtureBox
	interpolator := prepare(fixture)

	surrogate := interpolator.Compute(fixture.target)
	assert.Equal(surrogate, fixture.surrogate, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 1e-15, t)
}

func TestKraichnanOrszag(t *testing.T) {
	fixture := &fixtureKraichnanOrszag
	interpolator := prepare(fixture)

	surrogate := interpolator.Compute(fixture.target)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)
	assert.EqualWithin(surrogate.Integral, fixture.surrogate.Integral, 2e-14, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 6e-14, t)
}

func TestParabola(t *testing.T) {
	fixture := &fixtureParabola
	interpolator := prepare(fixture)

	surrogate := interpolator.Compute(fixture.target)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)
	assert.EqualWithin(surrogate.Integral, fixture.surrogate.Integral, 1e-6, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 1e-6, t)
}
