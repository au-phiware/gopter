package prop

import (
	"testing"
	"reflect"

	"github.com/leanovate/gopter"
)

func CheckForAllNoShrink(condition interface{}, gens ...gopter.Gen) func(*testing.T) {
	return ForAllNoShrinkT(condition, gens...).CheckT
}

func CheckForAllNoShrinkWithParameters(parameters *gopter.TestParameters, condition interface{}, gens ...gopter.Gen) func(*testing.T) {
	return ForAllNoShrinkT(condition, gens...).CheckWithParametersT(parameters)
}

/*
ForAllNoShrink creates a property that requires the check condition to be true for all values.
As the name suggests the generated values will not be shrinked if the condition falsiies.

"condition" has to be a function with the same number of parameters as the provided
generators "gens". The function may return a simple bool (true means that the
condition has passed), a string (empty string means that condition has passed),
a *PropResult, a func(*testing.T) or one of former combined with an error.
*/
func ForAllNoShrinkT(condition interface{}, gens ...gopter.Gen) gopter.PropT {
	callCheck, err := checkConditionFuncT(condition, len(gens))
	if err != nil {
		return func(_ *gopter.GenParameters) func(*testing.T) {
			return func(t *testing.T) {
				t.Fatal(err)
			}
		}
	}

	var pass int
	return gopter.SavePropT(func(genParams *gopter.GenParameters) func(*testing.T) {
		genResults := make([]*gopter.GenResult, len(gens))
		values := make([]reflect.Value, len(gens))
		var ok bool
		for i, gen := range gens {
			result := gen(genParams)
			genResults[i] = result
			values[i], ok = result.RetrieveAsValue()
			if !ok {
				return func(t *testing.T) {
					t.Skip()
				}
			}
		}
		runner := callCheck(values)
		return func(t *testing.T) {
			t.Helper()
			defer func(){
				if t.Failed() {
					t.Logf("Falsified after %d passed tests.", pass)
				} else {
					pass++
				}
			}()
			gopter.RunT(t, "#", runner, nil)
		}
	})
}
