package hybrid

import (
	"testing"

	"github.com/ready-steady/adapt/algorithm/external"
	"github.com/ready-steady/assert"
)

func TestBranin(t *testing.T) {
	fixture := &fixtureBranin
	interpolator, target := prepare(fixture)

	progresses := make([]external.Progress, 0)
	target.DoneHandler = func(progress *external.Progress) bool {
		progresses = append(progresses, *progress)
		return false
	}

	surrogate := interpolator.Compute(target)

	assert.Equal(progresses[:1], fixture.progresses[:1], t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 1e+42, t)
}
