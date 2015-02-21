package ode

import (
	"errors"
)

// Config contains the configuration of an integrator.
type Config struct {
	MaximalStep float64 // The maximal step size.
	InitialStep float64 // The initial step size.

	AbsoluteTolerance float64 // The absolute error tolerance.
	RelativeTolerance float64 // The relative error tolerance.
}

// DefaultConfig returns the default configuration of an integrator.
func DefaultConfig() *Config {
	return &Config{
		MaximalStep: 0,
		InitialStep: 0,

		AbsoluteTolerance: 1e-6,
		RelativeTolerance: 1e-3,
	}
}

func (c *Config) verify() error {
	if c.MaximalStep < 0 {
		return errors.New("the maximal step should be nonnegative")
	}
	if c.InitialStep < 0 {
		return errors.New("the initial step should be nonnegative")
	}

	if c.AbsoluteTolerance <= 0 {
		return errors.New("the absolute tolerance should be positive")
	}
	if c.RelativeTolerance <= 0 {
		return errors.New("the relative tolerance should be positive")
	}

	return nil
}
