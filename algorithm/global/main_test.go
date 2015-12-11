package global

import (
	"testing"

	"github.com/ready-steady/assert"
)

func TestBranin(t *testing.T) {
	fixture := &fixtureBranin
	interpolator, target, metric := prepare(fixture)

	progresses := make([]Progress, 0)
	target.MonitorHandler = func(progress *Progress) {
		progresses = append(progresses, *progress)
	}

	surrogate := interpolator.Compute(target, metric)

	assert.Equal(progresses, fixture.progresses, t)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)
	assert.EqualWithin(surrogate.Surpluses, fixture.surrogate.Surpluses, 1e-12, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 1e-12, t)
}

func BenchmarkBranin(b *testing.B) {
	fixture := &fixtureBranin
	interpolator, target, metric := prepare(fixture)
	for i := 0; i < b.N; i++ {
		interpolator.Compute(target, metric)
	}
}
