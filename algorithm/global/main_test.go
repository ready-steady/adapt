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
	assert.EqualWithin(values, fixture.values, 0.1, t)
}

func BenchmarkBranin(b *testing.B) {
	fixture := &fixtureBranin
	interpolator, target := prepare(fixture)
	for i := 0; i < b.N; i++ {
		interpolator.Compute(target)
	}
}
