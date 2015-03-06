package adhier

import (
	"testing"
)

func BenchmarkComputeHat(b *testing.B) {
	interpolator := prepare(&fixtureHat)

	for i := 0; i < b.N; i++ {
		interpolator.Compute(hat)
	}
}

func BenchmarkComputeCube(b *testing.B) {
	interpolator := prepare(&fixtureCube, func(config *Config) {
		config.MaxLevel = 9
	})

	for i := 0; i < b.N; i++ {
		interpolator.Compute(cube)
	}
}

func BenchmarkComputeBox(b *testing.B) {
	interpolator := prepare(&fixtureBox, func(config *Config) {
		config.MaxLevel = 9
	})

	for i := 0; i < b.N; i++ {
		interpolator.Compute(box)
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
	interpolator := prepare(&fixtureHat)
	surrogate := interpolator.Compute(hat)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		interpolator.Evaluate(surrogate, points)
	}
}

func BenchmarkEvaluateCube(b *testing.B) {
	interpolator := prepare(&fixtureCube, func(config *Config) {
		config.MaxLevel = 9
	})
	surrogate := interpolator.Compute(cube)
	points := generate(surrogate)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		interpolator.Evaluate(surrogate, points)
	}
}

func BenchmarkEvaluateBox(b *testing.B) {
	interpolator := prepare(&fixtureBox, func(config *Config) {
		config.MaxLevel = 9
	})
	surrogate := interpolator.Compute(box)
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
