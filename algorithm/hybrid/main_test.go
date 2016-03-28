package hybrid

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestBranin(t *testing.T) {
	fixture := &fixtureBranin
	interpolator := prepare(fixture)

	surrogate := interpolator.Compute(fixture.target)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 100.0, t)
}
