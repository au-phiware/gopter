package prop_test

import (
	"strings"
	"testing"
	"unicode"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func MisimplementedConcat(a, b string) string {
	if strings.IndexFunc(a, unicode.IsDigit) > 5 {
		return b
	}
	return a + b
}

// Example_invalidconcat demonstrates shrinking of string
// Kudos to @exarkun and @itamarst for finding this issue
func Example_invalidconcat() {
	parameters := gopter.DefaultTestParametersWithSeed(1234) // Example should generate reproducable results, otherwise DefaultTestParameters() will suffice

	properties := gopter.NewProperties(parameters)

	properties.Property("length is sum of lengths", prop.ForAllT(
		func(a, b string) bool {
			return MisimplementedConcat(a, b) == a+b
		},
		gen.Identifier().WithLabel("a"),
		gen.Identifier().WithLabel("b"),
	))

	// When using testing.T you might just use: properties.RunT(t)
	testing.Main(
		func(a, b string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{
				Name: "Example_invalidconcat",
				F:    properties.RunT,
			},
		}, nil, nil)
	// Output:
	//
	// --- FAIL: Example_invalidconcat (0.00s)
	// 	--- FAIL: Example_invalidconcat/length_is_sum_of_lengths (0.00s)
	// 		--- FAIL: Example_invalidconcat/length_is_sum_of_lengths/shrink_a#02 (0.00s)
	// 			check_condition_func.go:39: ARG_0: pbahbxh6
	// 			check_condition_func.go:39: ARG_1: dl
	// 		--- FAIL: Example_invalidconcat/length_is_sum_of_lengths/shrink_a#09 (0.00s)
	// 			check_condition_func.go:39: ARG_0: bahbxh6
	// 			check_condition_func.go:39: ARG_1: dl
	// 		--- FAIL: Example_invalidconcat/length_is_sum_of_lengths/shrink_b (0.00s)
	// 			check_condition_func.go:39: ARG_0: bahbxh6
	// 			check_condition_func.go:39: ARG_1: l
	// 		--- FAIL: Example_invalidconcat/length_is_sum_of_lengths/original#17 (0.00s)
	// 			check_condition_func.go:39: ARG_0: pkpbahbxh6
	// 			check_condition_func.go:39: ARG_1: dl
	// 		forall.go:64: Falsified after 17 passed tests.
	// 		runner.go:72: Completed with seed: 1234
	// FAIL
}
