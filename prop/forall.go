package prop

import (
	"fmt"
	"runtime/debug"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/leanovate/gopter"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

func CheckForAll(condition interface{}, gens ...gopter.Gen) func(*testing.T) {
	return ForAllT(condition, gens...).CheckT
}

func CheckForAllWithParameters(parameters *gopter.TestParameters, condition interface{}, gens ...gopter.Gen) func(*testing.T) {
	return ForAllT(condition, gens...).CheckWithParametersT(parameters)
}

/*
ForAll creates a property that requires the check condition to be true for all values, if the
condition falsiies the generated values will be shrinked.

"condition" has to be a function with the same number of parameters as the provided
generators "gens". The function may return a simple bool (true means that the
condition has passed), a string (empty string means that condition has passed),
a *PropResult, a func(*testing.T) or one of former combined with an error.
*/
func ForAllT(condition interface{}, gens ...gopter.Gen) gopter.PropT {
	callCheck, err := checkConditionFuncT(condition, len(gens))
	if err != nil {
		return func(genParams *gopter.GenParameters) func(*testing.T) {
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
			gopter.RunT(t, "original", runner, func(subT *testing.T) {
				subT.Helper()
				if r := recover(); r != nil {
					subT.Fatalf("Check paniced: %v %s", r, debug.Stack())
				}
				if subT.Failed() {
					for i, genResult := range genResults {
						name := nameForGenResult(i, genResult)
						nextValue := shrinkValueT(genParams.MaxShrinkCount, genResult, values[i].Interface(),
							func(v interface{}) bool {
								shrinkedOne := make([]reflect.Value, len(values))
								copy(shrinkedOne, values)
								if v == nil {
									shrinkedOne[i] = reflect.Zero(values[i].Type())
								} else {
									shrinkedOne[i] = reflect.ValueOf(v)
								}
								var success bool
								runner := callCheck(shrinkedOne)
								gopter.RunT(t, name, runner, func(t *testing.T) {
									success = !t.Failed()
								})
								return success
							})
						if nextValue == nil {
							values[i] = reflect.Zero(values[i].Type())
						} else {
							values[i] = reflect.ValueOf(nextValue)
						}
					}
				}
			})
		}
	})
}

func nameForGenResult(i int, genResult *gopter.GenResult) (name string) {
	if len(genResult.Labels) == 0 {
		name = strconv.Itoa(i)
	} else {
		name = strings.Join(genResult.Labels, "_")
	}
	name = fmt.Sprintf("shrink_%s", name)
	return
}

func shrinkValueT(maxShrinkCount int, genResult *gopter.GenResult, origValue interface{},
	check func(interface{}) bool) interface{} {
	lastValue := origValue

	shrinks := 0
	shrink := genResult.Shrinker(lastValue).Filter(genResult.Sieve)
	nextValue, ok := firstFailureT(shrink, check)
	for ok && shrinks < maxShrinkCount {
		shrinks++
		lastValue = nextValue

		shrink = genResult.Shrinker(lastValue).Filter(genResult.Sieve)
		nextValue, ok = firstFailureT(shrink, check)
	}

	return lastValue
}

func firstFailureT(shrink gopter.Shrink, check func(interface{}) bool) (value interface{}, ok bool) {
	value, ok = shrink()
	for ok {
		if !check(value) {
			return
		}
		value, ok = shrink()
	}
	return
}
