package adhier

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestComputeStep(t *testing.T) {
	fixture := &fixtureStep

	interpolator := prepare(fixture)
	surrogate := interpolator.Compute(step)

	assert.Equal(surrogate, fixture.surrogate, t)
}

func TestEvaluateStep(t *testing.T) {
	fixture := &fixtureStep

	interpolator := prepare(fixture)
	values := interpolator.Evaluate(fixture.surrogate, fixture.points)

	assert.Equal(values, fixture.values, t)
}

func TestComputeHat(t *testing.T) {
	fixture := &fixtureHat

	interpolator := prepare(fixture)
	surrogate := interpolator.Compute(hat)

	assert.Equal(surrogate, fixture.surrogate, t)
}

func TestEvaluateHat(t *testing.T) {
	fixture := &fixtureHat

	interpolator := prepare(fixture)
	values := interpolator.Evaluate(fixture.surrogate, fixture.points)

	assert.EqualWithin(values, fixture.values, 1e-15, t)
}

func TestComputeCube(t *testing.T) {
	fixture := &fixtureCube

	interpolator := prepare(fixture)
	surrogate := interpolator.Compute(cube)

	assert.Equal(surrogate, fixture.surrogate, t)
}

func TestComputeBox(t *testing.T) {
	fixture := &fixtureBox

	interpolator := prepare(fixture)
	surrogate := interpolator.Compute(box)

	assert.Equal(surrogate, fixture.surrogate, t)
}

func TestEvaluateBox(t *testing.T) {
	fixture := &fixtureBox

	interpolator := prepare(fixture)
	values := interpolator.Evaluate(fixture.surrogate, fixture.points)

	assert.EqualWithin(values, fixture.values, 1e-15, t)
}

func TestComputeEvaluateKraichnanOrszag(t *testing.T) {
	fixture := &fixtureKraichnanOrszag

	interpolator := prepare(fixture)
	surrogate := interpolator.Compute(kraichnanOrszag)

	assert.Equal(surrogate.Level, fixture.surrogate.Level, t)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)

	values := interpolator.Evaluate(surrogate, fixture.points)

	assert.EqualWithin(values, fixture.values, 6e-14, t)
}
