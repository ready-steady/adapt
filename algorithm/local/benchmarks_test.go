package local

import (
	"math/rand"
	"testing"

	"github.com/ready-steady/adapt/algorithm/external"
)

func BenchmarkComputeHat(b *testing.B) {
	fixture := &fixtureHat
	algorithm, strategy := prepare(fixture)

	for i := 0; i < b.N; i++ {
		algorithm.Compute(fixture.target, strategy)
	}
}

func BenchmarkComputeCube(b *testing.B) {
	fixture := &fixtureCube
	algorithm, strategy := prepare(fixture)

	for i := 0; i < b.N; i++ {
		algorithm.Compute(fixture.target, strategy)
	}
}

func BenchmarkComputeBox(b *testing.B) {
	fixture := &fixtureBox
	algorithm, strategy := prepare(fixture)
	strategy.(*Strategy).lmax = 9

	for i := 0; i < b.N; i++ {
		algorithm.Compute(fixture.target, strategy)
	}
}

func BenchmarkComputeMany(b *testing.B) {
	const (
		inputs  = 2
		outputs = 1000
	)

	fixture := &fixture{
		rule:   "closed",
		target: many(inputs, outputs),
		surrogate: &external.Surrogate{
			Inputs:  inputs,
			Outputs: outputs,
		},
	}
	algorithm, strategy := prepare(fixture)

	for i := 0; i < b.N; i++ {
		algorithm.Compute(fixture.target, strategy)
	}
}

func BenchmarkEvaluateHat(b *testing.B) {
	fixture := &fixtureHat
	algorithm, strategy := prepare(fixture)
	surrogate := algorithm.Compute(fixture.target, strategy)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		algorithm.Evaluate(surrogate, points)
	}
}

func BenchmarkEvaluateCube(b *testing.B) {
	fixture := &fixtureCube
	algorithm, strategy := prepare(fixture)
	strategy.(*Strategy).lmax = 9
	surrogate := algorithm.Compute(fixture.target, strategy)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		algorithm.Evaluate(surrogate, points)
	}
}

func BenchmarkEvaluateBox(b *testing.B) {
	fixture := &fixtureBox
	algorithm, strategy := prepare(fixture)
	strategy.(*Strategy).lmax = 9
	surrogate := algorithm.Compute(fixture.target, strategy)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		algorithm.Evaluate(surrogate, points)
	}
}

func BenchmarkEvaluateMany(b *testing.B) {
	const (
		inputs  = 2
		outputs = 1000
	)

	fixture := &fixture{
		rule:   "closed",
		target: many(inputs, outputs),
		surrogate: &external.Surrogate{
			Inputs:  inputs,
			Outputs: outputs,
		},
	}

	algorithm, strategy := prepare(fixture)
	surrogate := algorithm.Compute(fixture.target, strategy)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		algorithm.Evaluate(surrogate, points)
	}
}

func many(ni, no int) func([]float64, []float64) {
	return func(x, y []float64) {
		sum, value := 0.0, 0.0

		for i := 0; i < ni; i++ {
			sum += x[i]
		}

		if sum > float64(ni)/4 {
			value = 1
		}

		for i := 0; i < no; i++ {
			y[i] = value
		}
	}
}

func generate(surrogate *external.Surrogate) []float64 {
	const (
		count = 10000
	)

	generator := rand.New(rand.NewSource(0))
	points := make([]float64, count*surrogate.Inputs)
	for i := range points {
		points[i] = generator.Float64()
	}

	return points
}
