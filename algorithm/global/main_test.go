package global

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestBranin(t *testing.T) {
	fixture := &fixtureBranin
	interpolator, target := prepare(fixture)

	surrogate := interpolator.Compute(target)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 1e-12, t)
}
