package polynomial

import (
	"testing"

	"github.com/ready-steady/adapt/grid/equidistant"
	"github.com/ready-steady/adapt/internal"
	"github.com/ready-steady/assert"
)

func BenchmarkClosedCompute1(b *testing.B) {
	benchmarkClosedCompute(1, b)
}

func BenchmarkClosedCompute2(b *testing.B) {
	benchmarkClosedCompute(2, b)
}

func BenchmarkClosedCompute3(b *testing.B) {
	benchmarkClosedCompute(3, b)
}

func TestClosedCompute1D1P(t *testing.T) {
	basis := NewClosed(1, 1)

	compute := func(level, order uint64, point float64) float64 {
		return basis.Compute(internal.Compose([]uint64{level}, []uint64{order}), []float64{point})
	}

	points := []float64{0.0, 0.25, 0.5, 0.75, 1.0}

	cases := []struct {
		level  uint64
		order  uint64
		values []float64
	}{
		{0, 0, []float64{1.0, 1.0, 1.0, 1.0, 1.0}},
		{1, 0, []float64{1.0, 0.5, 0.0, 0.0, 0.0}},
		{1, 2, []float64{0.0, 0.0, 0.0, 0.5, 1.0}},
		{2, 1, []float64{0.0, 1.0, 0.0, 0.0, 0.0}},
		{2, 3, []float64{0.0, 0.0, 0.0, 1.0, 0.0}},
	}

	values := make([]float64, len(points))

	for i := range cases {
		for j := range values {
			values[j] = compute(cases[i].level, cases[i].order, points[j])
		}
		assert.Equal(values, cases[i].values, t)
	}
}

func TestClosedCompute1D3P(t *testing.T) {
	const (
		nd = 1
		np = 3
		nb = 4
		nn = 101
	)

	basis := NewClosed(nd, np)

	indices := internal.Compose(
		[]uint64{3, 3, 3, 3, 3, 3, 3, 3},
		[]uint64{1, 3, 5, 7, 9, 11, 13, 15},
	)

	points := make([]float64, nn)
	for i := range points {
		points[i] = float64(i) / (nn - 1)
	}

	values := make([]float64, nn)
	for i := range values {
		for j := 0; j < nb; j++ {
			values[i] += basis.Compute(indices[j:j+1], points[i:i+1])
		}
	}

	assert.Close(values, []float64{
		0.0000000000000000e+00, 2.0070399999999997e-01, 3.7683200000000000e-01,
		5.2940799999999999e-01, 6.5945600000000004e-01, 7.6800000000000013e-01,
		8.5606400000000005e-01, 9.2467200000000005e-01, 9.7484799999999983e-01,
		1.0076160000000001e+00, 1.0240000000000000e+00, 1.0250240000000002e+00,
		1.0117119999999999e+00, 9.8508799999999996e-01, 9.4617599999999991e-01,
		8.9599999999999991e-01, 8.3558399999999988e-01, 7.6595199999999986e-01,
		6.8812800000000007e-01, 6.0313600000000001e-01, 5.1199999999999990e-01,
		4.1574400000000011e-01, 3.1539200000000001e-01, 2.1196799999999993e-01,
		1.0649600000000009e-01, 0.0000000000000000e+00, 1.0649600000000009e-01,
		2.1196800000000018e-01, 3.1539200000000028e-01, 4.1574399999999978e-01,
		5.1199999999999990e-01, 6.0313600000000001e-01, 6.8812800000000007e-01,
		7.6595200000000008e-01, 8.3558400000000010e-01, 8.9599999999999991e-01,
		9.4617599999999991e-01, 9.8508799999999996e-01, 1.0117119999999999e+00,
		1.0250239999999999e+00, 1.0240000000000000e+00, 1.0076160000000001e+00,
		9.7484799999999994e-01, 9.2467200000000005e-01, 8.5606400000000005e-01,
		7.6799999999999990e-01, 6.5945599999999971e-01, 5.2940800000000032e-01,
		3.7683200000000028e-01, 2.0070400000000016e-01, 0.0000000000000000e+00,
		2.0070400000000016e-01, 3.7683200000000028e-01, 5.2940800000000032e-01,
		6.5945600000000049e-01, 7.6800000000000057e-01, 8.5606400000000038e-01,
		9.2467199999999972e-01, 9.7484799999999983e-01, 1.0076160000000001e+00,
		1.0240000000000000e+00, 1.0250239999999999e+00, 1.0117119999999999e+00,
		9.8508799999999996e-01, 9.4617599999999991e-01, 8.9599999999999991e-01,
		8.3558399999999988e-01, 7.6595199999999963e-01, 6.8812799999999963e-01,
		6.0313600000000045e-01, 5.1200000000000045e-01, 4.1574400000000039e-01,
		3.1539200000000028e-01, 2.1196800000000018e-01, 1.0649600000000009e-01,
		0.0000000000000000e+00, 1.0649600000000009e-01, 2.1196800000000018e-01,
		3.1539200000000028e-01, 4.1574400000000039e-01, 5.1200000000000045e-01,
		6.0313600000000045e-01, 6.8812799999999963e-01, 7.6595199999999963e-01,
		8.3558399999999988e-01, 8.9599999999999991e-01, 9.4617599999999991e-01,
		9.8508799999999996e-01, 1.0117119999999999e+00, 1.0250239999999999e+00,
		1.0240000000000000e+00, 1.0076160000000001e+00, 9.7484799999999983e-01,
		9.2467199999999972e-01, 8.5606400000000038e-01, 7.6800000000000057e-01,
		6.5945600000000049e-01, 5.2940800000000032e-01, 3.7683200000000028e-01,
		2.0070400000000016e-01, 0.0000000000000000e+00,
	}, 1e-15, t)
}

func TestClosedIntegrate(t *testing.T) {
	basis := NewClosed(1, 1)

	levels := []uint64{0, 1, 2, 3}
	values := []float64{1.0, 0.25, 1.0 / 2.0 / 2.0, 1.0 / 2.0 / 2.0 / 2.0}

	for i := range levels {
		indices := internal.Compose([]uint64{levels[i]}, []uint64{0})
		assert.Equal(basis.Integrate(indices), values[i], t)
	}
}

func benchmarkClosedCompute(power uint, b *testing.B) {
	const (
		nd = 10
		ns = 100000
	)

	basis := NewClosed(nd, power)
	indices := generateIndices(nd, ns, equidistant.NewClosed(nd).Refine)
	points := generatePoints(nd, ns, indices, basis.grid.Node)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < ns; j++ {
			basis.Compute(indices[j*nd:(j+1)*nd], points[j*nd:(j+1)*nd])
		}
	}
}
