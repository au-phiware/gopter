package arbitrary_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/arbitrary"
)

func Example_parseint() {
	parameters := gopter.DefaultTestParametersWithSeed(1234) // Example should generate reproducable results, otherwise DefaultTestParameters() will suffice

	arbitraries := arbitrary.DefaultArbitraries()
	properties := gopter.NewProperties(parameters)

	properties.Property("printed integers can be parsed", arbitraries.ForAllT(
		func(a int64) bool {
			str := fmt.Sprintf("%d", a)
			parsed, err := strconv.ParseInt(str, 10, 64)
			return err == nil && parsed == a
		}))

	// When using testing.T you might just use: properties.RunT(t)
	testing.Main(
		func(a, b string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{
				Name: "Example_parseint",
				F:    properties.RunT,
			},
		}, nil, nil)
	// Output:
	//
	// PASS
}
