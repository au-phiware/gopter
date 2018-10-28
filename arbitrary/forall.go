package arbitrary

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

func (a *Arbitraries) CheckForAll(condition interface{}) func(*testing.T) {
	return a.ForAllT(condition).CheckT
}

func (a *Arbitraries) CheckForAllWithParameters(parameters *gopter.TestParameters, condition interface{}) func(*testing.T) {
	return a.ForAllT(condition).CheckWithParametersT(parameters)
}

/*
ForAll creates a property that requires the check condition to be true for all
values, if the condition falsiies the generated values will be shrinked.

"condition" has to be a function with the any number of parameters that can
generated in context of the Arbitraries. The function may return a simple bool,
a *PropResult, a boolean with error or a *PropResult with error.
*/
func (a *Arbitraries) ForAllT(condition interface{}) gopter.PropT {
	conditionVal := reflect.ValueOf(condition)
	conditionType := conditionVal.Type()

	if conditionType.Kind() != reflect.Func {
		return prop.ErrorPropT(fmt.Errorf("Param of ForrAll has to be a func: %v", conditionType.Kind()))
	}

	gens := make([]gopter.Gen, conditionType.NumIn())
	for i := 0; i < conditionType.NumIn(); i++ {
		gens[i] = a.GenForType(conditionType.In(i))
	}

	return prop.ForAllT(condition, gens...)
}
