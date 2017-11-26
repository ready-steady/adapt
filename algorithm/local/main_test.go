package local

import (
	"testing"

	"github.com/ready-steady/assert"

	interpolation "github.com/ready-steady/adapt/algorithm"
)

func TestBox(t *testing.T) {
	fixture := &fixtureBox
	algorithm, strategy := prepare(fixture)

	surrogate := algorithm.Compute(fixture.target, strategy)
	assert.Equal(surrogate, fixture.surrogate, t)

	values := algorithm.Evaluate(surrogate, fixture.points)
	assert.Close(values, fixture.values, 1e-15, t)
}

func TestCube(t *testing.T) {
	fixture := &fixtureCube
	algorithm, strategy := prepare(fixture)

	surrogate := algorithm.Compute(fixture.target, strategy)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)
	assert.Equal(surrogate.Integral, fixture.surrogate.Integral, t)
	assert.Equal(interpolation.Validate(surrogate.Indices, surrogate.Inputs,
		fixture.grid), true, t)

	values := algorithm.Evaluate(surrogate, fixture.points)
	assert.Close(values, fixture.values, 2e-15, t)
}

func TestHat(t *testing.T) {
	fixture := &fixtureHat
	algorithm, strategy := prepare(fixture)

	surrogate := algorithm.Compute(fixture.target, strategy)
	assert.Equal(surrogate, fixture.surrogate, t)

	values := algorithm.Evaluate(surrogate, fixture.points)
	assert.Close(values, fixture.values, 1e-15, t)
}

func TestKraichnanOrszag(t *testing.T) {
	fixture := &fixtureKraichnanOrszag
	algorithm, strategy := prepare(fixture)

	surrogate := algorithm.Compute(fixture.target, strategy)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)
	assert.Close(surrogate.Integral, fixture.surrogate.Integral, 2e-14, t)
	assert.Equal(interpolation.Validate(surrogate.Indices, surrogate.Inputs,
		fixture.grid), true, t)

	values := algorithm.Evaluate(surrogate, fixture.points)
	assert.Close(values, fixture.values, 6e-14, t)
}

func TestParabola(t *testing.T) {
	fixture := &fixtureParabola
	algorithm, strategy := prepare(fixture)

	surrogate := algorithm.Compute(fixture.target, strategy)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)
	assert.Close(surrogate.Integral, fixture.surrogate.Integral, 1e-6, t)
	assert.Equal(interpolation.Validate(surrogate.Indices, surrogate.Inputs,
		fixture.grid), true, t)

	values := algorithm.Evaluate(surrogate, fixture.points)
	assert.Close(values, fixture.values, 1e-6, t)
}

func TestStep(t *testing.T) {
	fixture := &fixtureStep
	algorithm, strategy := prepare(fixture)

	surrogate := algorithm.Compute(fixture.target, strategy)
	assert.Equal(surrogate, fixture.surrogate, t)

	values := algorithm.Evaluate(surrogate, fixture.points)
	assert.Equal(values, fixture.values, t)
}
