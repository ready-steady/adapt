package ode

import (
	"errors"
)

type Config struct {
	MaximalStep float64
	InitialStep float64

	AbsoluteTolerance float64
	RelativeTolerance float64
}

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
