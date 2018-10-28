package arbitrary_test

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/arbitrary"
)

type MyStringType string
type MyInt8Type int8
type MyInt16Type int16
type MyInt32Type int32
type MyInt64Type int64
type MyUInt8Type uint8
type MyUInt16Type uint16
type MyUInt32Type uint32
type MyUInt64Type uint64

type Foo struct {
	Name MyStringType
	Id1  MyInt8Type
	Id2  MyInt16Type
	Id3  MyInt32Type
	Id4  MyInt64Type
	Id5  MyUInt8Type
	Id6  MyUInt16Type
	Id7  MyUInt32Type
	Id8  MyUInt64Type
}

func Example_arbitrary_structs() {
	parameters := gopter.DefaultTestParametersWithSeed(1234) // Example should generate reproducable results, otherwise DefaultTestParameters() will suffice

	arbitraries := arbitrary.DefaultArbitraries()

	properties := gopter.NewProperties(parameters)

	properties.Property("MyInt64", arbitraries.ForAllT(
		func(id MyInt64Type) bool {
			return id > -1000
		}))
	properties.Property("MyUInt32Type", arbitraries.ForAllT(
		func(id MyUInt32Type) bool {
			return id < 2000
		}))
	properties.Property("Foo", arbitraries.ForAllT(
		func(foo *Foo) bool {
			return true
		}))
	properties.Property("Foo2", arbitraries.ForAllT(
		func(foo Foo) bool {
			return true
		}))

	testing.Main(
		func(a, b string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{
				Name: "Example_arbitrary_structs",
				F:    properties.RunT,
			},
		}, nil, nil)
	// Output:
	//
	// --- FAIL: Example_arbitrary_structs (0.02s)
	// 	--- FAIL: Example_arbitrary_structs/MyInt64 (0.00s)
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#01 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -800533414872418627
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#03 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -400266707436209314
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#05 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -200133353718104657
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#07 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -100066676859052329
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#09 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -50033338429526165
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#11 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -25016669214763083
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#13 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -12508334607381542
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#15 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -6254167303690771
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#17 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -3127083651845386
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#19 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -1563541825922693
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#21 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -781770912961347
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#23 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -390885456480674
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#25 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -195442728240337
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#27 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -97721364120169
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#29 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -48860682060085
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#31 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -24430341030043
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#33 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -12215170515022
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#35 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -6107585257511
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#37 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -3053792628756
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#39 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -1526896314378
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#41 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -763448157189
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#43 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -381724078595
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#45 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -190862039298
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#47 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -95431019649
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#49 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -47715509825
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#51 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -23857754913
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#53 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -11928877457
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#55 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -5964438729
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#57 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -2982219365
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#59 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -1491109683
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#61 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -745554842
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#63 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -372777421
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#65 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -186388711
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#67 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -93194356
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#69 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -46597178
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#71 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -23298589
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#73 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -11649295
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#75 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -5824648
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#77 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -2912324
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#79 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -1456162
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#81 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -728081
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#83 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -364041
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#85 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -182021
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#87 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -91011
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#89 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -45506
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#91 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -22753
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#93 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -11377
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#95 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -5689
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#97 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -2845
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#99 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -1423
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#103 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -1068
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#111 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -1002
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#129 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -1001
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/shrink_0#147 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -1000
	// 		--- FAIL: Example_arbitrary_structs/MyInt64/original#06 (0.00s)
	// 			check_condition_func.go:39: ARG_0: -1601066829744837253
	// 		forall.go:64: Falsified after 6 passed tests.
	// 		runner.go:71: Elapsed time: 3.332632ms
	// 		runner.go:72: Completed with seed: 1234
	// 	--- FAIL: Example_arbitrary_structs/MyUInt32Type (0.00s)
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#01 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 1080961160
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#03 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 540480580
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#05 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 270240290
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#07 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 135120145
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#09 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 67560073
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#11 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 33780037
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#13 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 16890019
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#15 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 8445010
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#17 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 4222505
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#19 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 2111253
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#21 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 1055627
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#23 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 527814
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#25 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 263907
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#27 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 131954
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#29 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 65977
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#31 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 32989
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#33 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 16495
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#35 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 8248
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#37 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 4124
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#39 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 2062
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#46 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 2030
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#54 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 2015
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/shrink_0#62 (0.00s)
	// 			check_condition_func.go:39: ARG_0: 2000
	// 		--- FAIL: Example_arbitrary_structs/MyUInt32Type/original (0.00s)
	// 			check_condition_func.go:39: ARG_0: 2161922319
	// 		forall.go:64: Falsified after 0 passed tests.
	// 		runner.go:71: Elapsed time: 1.503797ms
	// 		runner.go:72: Completed with seed: 1234
	// FAIL
}
