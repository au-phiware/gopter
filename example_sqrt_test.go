package gopter_test

import (
	"testing"
	"math"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func Example_sqrt() {
	parameters := gopter.DefaultTestParameters()
	parameters.Rng.Seed(1234) // Just for this example to generate reproducable results

	properties := gopter.NewProperties(parameters)

	properties.Property("greater one of all greater one", prop.ForAllT(
		func(v float64) bool {
			return math.Sqrt(v) >= 1
		},
		gen.Float64().SuchThat(func(x float64) bool { return x >= 1.0 }),
	))

	properties.Property("squared is equal to value", prop.ForAllT(
		func(v float64) bool {
			r := math.Sqrt(v)
			return math.Abs(r*r-v) < 1e-10*v
		},
		gen.Float64().SuchThat(func(x float64) bool { return x >= 0.0 }),
	))

	testing.Main(
		func(a, b string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{
				Name: "Example_sqrt",
				F:    properties.RunT,
			},
		}, nil, nil)
	// Output:
	//
	// PASS
}
