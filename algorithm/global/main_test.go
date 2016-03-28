package global

import (
	"testing"

	"github.com/ready-steady/adapt/algorithm/internal"
	"github.com/ready-steady/assert"
)

func BenchmarkBranin(b *testing.B) {
	fixture := &fixtureBranin
	interpolator := prepare(fixture)
	for i := 0; i < b.N; i++ {
		interpolator.Compute(fixture.target)
	}
}

func TestBranin(t *testing.T) {
	fixture := &fixtureBranin
	interpolator := prepare(fixture)

	surrogate := interpolator.Compute(fixture.target)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)
	assert.Equal(internal.IsUnique(surrogate.Indices, surrogate.Inputs), true, t)
	assert.Equal(internal.IsAdmissible(surrogate.Indices, surrogate.Inputs), true, t)

	values := interpolator.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 0.1, t)
}
