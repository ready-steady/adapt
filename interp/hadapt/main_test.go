package hadapt

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/gomath/format/mat"
	"github.com/gomath/numan/basis/linhat"
	"github.com/gomath/numan/grid/newcot"
)

func assertEqual(actual, expected interface{}, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got '%v' instead of '%v'", actual, expected)
	}
}

func TestConstruct1D(t *testing.T) {
	algorithm := New(newcot.New(1), linhat.New(1))
	algorithm.maxLevel = 4

	surrogate := algorithm.Construct(step)

	levels := []uint8{0, 1, 1, 2, 3, 3, 4, 4}
	orders := []uint32{0, 0, 2, 3, 5, 7, 9, 11}
	surpluses := []float64{1, 0, -1, -0.5, -0.5, 0, -0.5, 0}

	assertEqual(surrogate.level, uint8(4), t)
	assertEqual(surrogate.inCount, uint16(1), t)
	assertEqual(surrogate.nodeCount, uint32(8), t)

	assertEqual(surrogate.levels, levels, t)
	assertEqual(surrogate.orders, orders, t)
	assertEqual(surrogate.surpluses, surpluses, t)
}

func TestConstruct2D(t *testing.T) {
	algorithm := New(newcot.New(2), linhat.New(2))
	algorithm.maxLevel = 3

	surrogate := algorithm.Construct(cube)

	levels := []uint8{
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
	}

	orders := []uint32{
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
	}

	surpluses := []float64{
		1.0, -1.0, -1.0, -1.0, -1.0, -0.5, 1.0, 1.0, -0.5, 1.0,
		1.0, -0.5, -0.5,  0.0,  0.5,  0.5, 0.5, 0.5,  0.5, 0.5,
		0.0,  0.5,  0.5,  0.5,  0.5,  0.0, 0.5, 0.5,  0.0}

	assertEqual(surrogate.level, uint8(3), t)
	assertEqual(surrogate.inCount, uint16(2), t)
	assertEqual(surrogate.nodeCount, uint32(29), t)

	assertEqual(surrogate.levels, levels, t)
	assertEqual(surrogate.orders, orders, t)
	assertEqual(surrogate.surpluses, surpluses, t)
}

func TestEvaluate1D(t *testing.T) {
	algorithm := New(newcot.New(1), linhat.New(1))

	surrogate := &Surrogate{
		level:     4,
		inCount:   1,
		nodeCount: 8,
		levels:    []uint8{0, 1, 1, 2, 3, 3, 4, 4},
		orders:    []uint32{0, 0, 2, 3, 5, 7, 9, 11},
		surpluses: []float64{1, 0, -1, -0.5, -0.5, 0, -0.5, 0},
	}

	points := []float64{0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1}
	values := []float64{1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0}

	assertEqual(algorithm.Evaluate(surrogate, points), values, t)
}

func ExampleStep() {
	algorithm := New(newcot.New(1), linhat.New(1))
	algorithm.maxLevel = 19
	surrogate := algorithm.Construct(step)

	fmt.Println(surrogate)

	// Output:
	// Surrogate{ inputs: 1, levels: 19, nodes: 38 }
}

func ExampleHat() {
	algorithm := New(newcot.New(1), linhat.New(1))
	algorithm.maxLevel = 9
	surrogate := algorithm.Construct(hat)

	fmt.Println(surrogate)

	if !testing.Verbose() {
		return
	}

	points := make([]float64, 101)
	for i := range points {
		points[i] = 0.01 * float64(i)
	}
	values := algorithm.Evaluate(surrogate, points)

	file, _ := mat.Open("hat.mat", "w7.3")
	defer file.Close()

	file.PutMatrix("x", 101, 1, points)
	file.PutMatrix("y", 101, 1, values)

	// Output:
	// Surrogate{ inputs: 1, levels: 9, nodes: 305 }
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
