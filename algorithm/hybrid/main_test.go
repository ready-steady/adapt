package hybrid

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestBranin(t *testing.T) {
	fixture := &fixtureBranin
	interpolator, target := prepare(fixture)

	surrogate := interpolator.Compute(target)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 100.0, t)
}
