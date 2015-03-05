package rk4

import (
	"testing"

	"github.com/ready-steady/support/assert"
)

func TestComputeSimple(t *testing.T) {
	fixture := &fixtureSimple

	solution := Compute(fixture.dydx, fixture.y0, fixture.x0, fixture.Δx, fixture.n)

	assert.Equal(solution, fixture.solution, t)
}
