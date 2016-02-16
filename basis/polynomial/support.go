package polynomial

import (
	"github.com/ready-steady/quadrature"
)

var ruleCache map[uint]*rule = make(map[uint]*rule)

type rule struct {
	x []float64
	w []float64
}

func getRule(order uint) *rule {
	if rule, ok := ruleCache[order]; ok {
		return rule
	}
	x, w := quadrature.Legendre(order, 0.0, 1.0)
	rule := &rule{x, w}
	ruleCache[order] = rule
	return rule
}

func integrate(a, b float64, order uint, target func(float64) float64) float64 {
	value, rule := 0.0, getRule(order)
	for i := uint(0); i < order; i++ {
		x := a + (b-a)*rule.x[i]
		w := (b - a) * rule.w[i]
		value += w * target(x)
	}
	return value
}
