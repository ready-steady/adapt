package hybrid

import (
	"testing"

	"github.com/ready-steady/adapt/algorithm/internal"
	"github.com/ready-steady/assert"
)

func TestBranin(t *testing.T) {
	fixture := &fixtureBranin
	algorithm, strategy := prepare(fixture)

	surrogate := algorithm.Compute(fixture.target, strategy)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)
	assert.Equal(internal.IsUnique(surrogate.Indices, surrogate.Inputs), true, t)
	assert.Equal(internal.IsAdmissible(surrogate.Indices,
		surrogate.Inputs, fixture.parent), true, t)

	values := algorithm.Evaluate(surrogate, fixture.points)
	assert.EqualWithin(values, fixture.values, 0.1, t)
}
