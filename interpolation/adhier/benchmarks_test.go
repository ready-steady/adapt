package adhier

import (
	"testing"
)

func BenchmarkComputeHat(b *testing.B) {
	fixture := &fixtureHat
	interpolator := prepare(fixture)

	for i := 0; i < b.N; i++ {
		interpolator.Compute(fixture.target)
	}
}

func BenchmarkComputeCube(b *testing.B) {
	fixture := &fixtureCube
	interpolator := prepare(fixture, func(config *Config) {
		config.MaxLevel = 9
	})

	for i := 0; i < b.N; i++ {
		interpolator.Compute(fixture.target)
	}
}

func BenchmarkComputeBox(b *testing.B) {
	fixture := &fixtureBox
	interpolator := prepare(fixture, func(config *Config) {
		config.MaxLevel = 9
	})

	for i := 0; i < b.N; i++ {
		interpolator.Compute(fixture.target)
	}
}

func BenchmarkComputeMany(b *testing.B) {
	const (
		inputs  = 2
		outputs = 1000
	)

	interpolator := prepare(&fixture{
		surrogate: &Surrogate{
			Inputs:  inputs,
			Outputs: outputs,

			Level: 9,
		},
	})
	function := many(inputs, outputs)

	for i := 0; i < b.N; i++ {
		interpolator.Compute(function)
	}
}

func BenchmarkEvaluateHat(b *testing.B) {
	fixture := &fixtureHat
	interpolator := prepare(fixture)
	surrogate := interpolator.Compute(fixture.target)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		interpolator.Evaluate(surrogate, points)
	}
}

func BenchmarkEvaluateCube(b *testing.B) {
	fixture := &fixtureCube
	interpolator := prepare(fixture, func(config *Config) {
		config.MaxLevel = 9
	})
	surrogate := interpolator.Compute(fixture.target)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		interpolator.Evaluate(surrogate, points)
	}
}

func BenchmarkEvaluateBox(b *testing.B) {
	fixture := &fixtureBox
	interpolator := prepare(fixture, func(config *Config) {
		config.MaxLevel = 9
	})
	surrogate := interpolator.Compute(fixture.target)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		interpolator.Evaluate(surrogate, points)
	}
}

func BenchmarkEvaluateMany(b *testing.B) {
	const (
		inputs  = 2
		outputs = 1000
	)

	interpolator := prepare(&fixture{
		surrogate: &Surrogate{
			Inputs:  inputs,
			Outputs: outputs,

			Level: 9,
		},
	})
	function := many(inputs, outputs)
	surrogate := interpolator.Compute(function)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		interpolator.Evaluate(surrogate, points)
	}
}

func many(ni, no int) func([]float64, []float64, []uint64) {
	return func(x, y []float64, _ []uint64) {
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
