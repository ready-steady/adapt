package hybrid

import (
	"testing"

	"github.com/ready-steady/assert"

	interpolation "github.com/ready-steady/adapt/algorithm"
)

func TestBranin(t *testing.T) {
	fixture := &fixtureBranin
	algorithm, strategy := prepare(fixture)

	surrogate := algorithm.Compute(fixture.target, strategy)
	assert.Equal(surrogate.Nodes, fixture.surrogate.Nodes, t)
	assert.Equal(interpolation.Validate(surrogate.Indices, surrogate.Inputs,
		fixture.grid), true, t)

	values := algorithm.Evaluate(surrogate, fixture.points)
	assert.Close(values, fixture.values, 0.1, t)
}
