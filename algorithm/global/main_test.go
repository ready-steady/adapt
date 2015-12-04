package global

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestBranin(t *testing.T) {
	fixture := &fixtureBranin
	interpolator, target := prepare(fixture)

	surrogate := interpolator.Compute(target)
	assert.Equal(surrogate, fixture.surrogate, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.Equal(values, fixture.values, t)
}
