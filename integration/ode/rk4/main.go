// Package rk4 provides an integrator of system of ordinary differential
// equations based on the fourth-order Runge–Kutta method.
//
// https://en.wikipedia.org/wiki/Runge–Kutta_methods
package rk4

// Compute integrates the system of differential equations dy/dx = f(x, y) with
// the initial condition y(x₀) = y₀. The solution is returned at n equidistant
// points starting from and including x₀ with the step size Δx.
func Compute(dydx func(float64, []float64, []float64), y0 []float64,
	x0, Δx float64, n uint) []float64 {

	nd, ns := len(y0), int(n)

	solution := make([]float64, ns*nd)
	copy(solution, y0)

	z := make([]float64, nd)
	f1 := make([]float64, nd)
	f2 := make([]float64, nd)
	f3 := make([]float64, nd)
	f4 := make([]float64, nd)

	for k, x, y := 1, x0, y0; k < ns; k++ {
		// Step 1
		dydx(x, y, f1)

		// Step 2
		for i := 0; i < nd; i++ {
			z[i] = y[i] + Δx*f1[i]/2
		}
		dydx(x+Δx/2, z, f2)

		// Step 3
		for i := 0; i < nd; i++ {
			z[i] = y[i] + Δx*f2[i]/2
		}
		dydx(x+Δx/2, z, f3)

		// Step 4
		for i := 0; i < nd; i++ {
			z[i] = y[i] + Δx*f3[i]
		}
		dydx(x+Δx, z, f4)

		ynew := solution[k*nd:]
		for i := 0; i < nd; i++ {
			ynew[i] = y[i] + Δx*(f1[i]+2*f2[i]+2*f3[i]+f4[i])/6
		}

		x += Δx
		y = ynew
	}

	return solution
}
