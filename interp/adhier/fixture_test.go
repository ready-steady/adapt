package adhier

import (
	"math"
)

type fixture struct {
	surrogate *Surrogate
	points    []float64
	values    []float64
}

func step(x []float64) []float64 {
	y := make([]float64, len(x))
	for i := range x {
		if x[i] <= 0.5 {
			y[i] = 1
		}
	}
	return y
}

var stepFixture = fixture{
	surrogate: &Surrogate{
		level:     4,
		inCount:   1,
		outCount:  1,
		nodeCount: 8,

		levels:    []uint8{0, 1, 1, 2, 3, 3, 4, 4},
		orders:    []uint32{0, 0, 2, 3, 5, 7, 9, 11},
		surpluses: []float64{1, 0, -1, -0.5, -0.5, 0, -0.5, 0},
	},
	points: []float64{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1},
	values: []float64{1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0},
}

func cube(x []float64) []float64 {
	count := uint16(len(x)) / 2
	y := make([]float64, count)

	for i := uint16(0); i < count; i++ {
		if math.Abs(2*x[2*i]-1) < 0.45 && math.Abs(2*x[2*i+1]-1) < 0.45 {
			y[i] = 1
		}
	}

	return y
}

var cubeFixture = fixture{
	surrogate: &Surrogate{
		level:     3,
		inCount:   2,
		outCount:  1,
		nodeCount: 29,

		levels: []uint8{
			0, 0,
			1, 0,
			1, 0,
			0, 1,
			0, 1,
			2, 0,
			1, 1,
			1, 1,
			2, 0,
			1, 1,
			1, 1,
			0, 2,
			0, 2,
			3, 0,
			3, 0,
			2, 1,
			2, 1,
			1, 2,
			1, 2,
			3, 0,
			3, 0,
			2, 1,
			2, 1,
			1, 2,
			1, 2,
			0, 3,
			0, 3,
			0, 3,
			0, 3,
		},

		orders: []uint32{
			0, 0,
			0, 0,
			2, 0,
			0, 0,
			0, 2,
			1, 0,
			0, 0,
			0, 2,
			3, 0,
			2, 0,
			2, 2,
			0, 1,
			0, 3,
			1, 0,
			3, 0,
			1, 0,
			1, 2,
			0, 1,
			0, 3,
			5, 0,
			7, 0,
			3, 0,
			3, 2,
			2, 1,
			2, 3,
			0, 1,
			0, 3,
			0, 5,
			0, 7,
		},

		surpluses: []float64{
			1.0, -1.0, -1.0, -1.0, -1.0, -0.5, 1.0, 1.0, -0.5, 1.0,
			1.0, -0.5, -0.5, 0.0, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5,
			0.0, 0.5, 0.5, 0.5, 0.5, 0.0, 0.5, 0.5, 0.0,
		},
	},
}

func hat(x []float64) []float64 {
	y := make([]float64, len(x))
	for i, z := range x {
		z = 5*z - 1
		switch {
		case 0 <= z && z < 1:
			y[i] = 0.5 * z * z
		case 1 <= z && z < 2:
			y[i] = 0.5 * (-2*z*z + 6*z - 3)
		case 2 <= z && z < 3:
			y[i] = 0.5 * (3 - z) * (3 - z)
		}
	}
	return y
}

func box(x []float64) []float64 {
	count := len(x) / 2
	y := make([]float64, 3*count)

	for i := 0; i < count; i++ {
		x1, x2 := x[2*i+0], x[2*i+1]

		if x1 + x2 > 0.5 {
			y[3*i+0] = 1
		}

		if x1 - x2 > 0.5 {
			y[3*i+1] = 1
		}

		if x2 - x1 > 0.5 {
			y[3*i+2] = 1
		}
	}

	return y
}

var boxFixture = fixture{
	surrogate: &Surrogate{
		level:     3,
		inCount:   2,
		outCount:  3,
		nodeCount: 20,
	},
}
