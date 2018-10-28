package prop

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func checkConditionFuncT(check interface{}, numArgs int) (func([]reflect.Value) func(*testing.T), error) {
	checkVal := reflect.ValueOf(check)
	checkType := checkVal.Type()

	if checkType.Kind() != reflect.Func {
		return nil, fmt.Errorf("First param of ForrAll has to be a func: %v", checkVal.Kind())
	}
	if checkType.NumIn() != numArgs {
		return nil, fmt.Errorf("Number of parameters does not match number of generators: %d != %d", checkType.NumIn(), numArgs)
	}
	if checkType.NumOut() == 0 {
		return nil, errors.New("At least one output parameters is required")
	} else if checkType.NumOut() > 2 {
		return nil, fmt.Errorf("No more than 2 output parameters are allowed: %d", checkType.NumOut())
	} else if checkType.NumOut() == 2 && !checkType.Out(1).Implements(typeOfError) {
		return nil, fmt.Errorf("No 2 output has to be error: %v", checkType.Out(1).Kind())
	}
	return func(values []reflect.Value) func(*testing.T) {
		var runner func(*testing.T)
		results := checkVal.Call(values)
		if checkType.NumOut() == 1 || results[1].IsNil() {
			runner = convertResultT(results[0].Interface(), nil)
		} else {
			runner = convertResultT(results[0].Interface(), results[1].Interface().(error))
		}
		return func(t *testing.T)  {
			t.Helper()
			defer func() {
				for i, arg := range values {
					t.Logf("ARG_%d: %v", i, arg)
				}
			}()
			runner(t)
		}
	}, nil
}
