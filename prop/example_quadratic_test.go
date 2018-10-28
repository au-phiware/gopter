package prop_test

import (
	"errors"
	"math"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func solveQuadratic(a, b, c float64) (float64, float64, error) {
	if a == 0 {
		return 0, 0, errors.New("No solution")
	}
	v := b*b - 4*a*c
	if v < 0 {
		return 0, 0, errors.New("No solution")
	}
	v = math.Sqrt(v)
	return (-b + v) / 2 / a, (-b - v) / 2 / a, nil
}

func Example_quadratic() {
	parameters := gopter.DefaultTestParametersWithSeed(1234) // Example should generate reproducable results, otherwise DefaultTestParameters() will suffice

	properties := gopter.NewProperties(parameters)

	properties.Property("solve quadratic", prop.ForAllT(
		func(a, b, c float64) bool {
			x1, x2, err := solveQuadratic(a, b, c)
			if err != nil {
				return true
			}
			return math.Abs(a*x1*x1+b*x1+c) < 1e-5 && math.Abs(a*x2*x2+b*x2+c) < 1e-5
		},
		gen.Float64(),
		gen.Float64(),
		gen.Float64(),
	))

	properties.Property("solve quadratic with resonable ranges", prop.ForAllT(
		func(a, b, c float64) bool {
			x1, x2, err := solveQuadratic(a, b, c)
			if err != nil {
				return true
			}
			return math.Abs(a*x1*x1+b*x1+c) < 1e-5 && math.Abs(a*x2*x2+b*x2+c) < 1e-5
		},
		gen.Float64Range(-1e8, 1e8),
		gen.Float64Range(-1e8, 1e8),
		gen.Float64Range(-1e8, 1e8),
	))

	testing.Main(
		func(a, b string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{
				Name: "Example_sqrt",
				F:    properties.RunT,
			},
		}, nil, nil)
}
