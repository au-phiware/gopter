package prop_test

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func Example_shrink() {
	parameters := gopter.DefaultTestParametersWithSeed(0xdeadbeef) // Example should generate reproducable results, otherwise DefaultTestParameters() will suffice

	properties := gopter.NewProperties(parameters)

	properties.Property("fail", prop.ForAllT(
		func(arg int64) bool {
			return arg <= 0x7e7e67ef432d80
		},
		gen.Int64().WithLabel("arg"),
	))

	properties.Property("fail no shrink", prop.ForAllNoShrinkT(
		func(arg int64) bool {
			return arg <= 100
		},
		gen.Int64(),
	))

	// When using testing.T you might just use: properties.RunT(t)
	testing.Main(
		func(a, b string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{
				Name: "Example_shrink",
				F:    properties.RunT,
			},
		}, nil, nil)
	// Output:
	//
	// --- FAIL: Example_shrink (0.00s)
	// 	--- FAIL: Example_shrink/fail (0.00s)
	// 		--- FAIL: Example_shrink/fail/shrink_arg#01 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 142419327705724418
	// 		--- FAIL: Example_shrink/fail/shrink_arg#03 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 71209663852862209
	// 		--- FAIL: Example_shrink/fail/shrink_arg#05 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 35604831926431105
	// 		--- FAIL: Example_shrink/fail/original#01 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 284838655411448835
	// 		forall.go:64: Falsified after 1 passed tests.
	// 		runner.go:72: Completed with seed: 3735928559
	// 	--- FAIL: Example_shrink/fail_no_shrink (0.00s)
	// 		--- FAIL: Example_shrink/fail_no_shrink/#05 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 794074141116934704
	// 		forall_no_shrink.go:58: Falsified after 5 passed tests.
	// 		runner.go:72: Completed with seed: 3735928559
	// FAIL
}
