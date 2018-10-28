// +build fail

package prop_test

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)


func TestForAllNoShrinkInvalidParam(t *testing.T) {
	fail := prop.ForAllNoShrinkT(0)
	fail(gopter.DefaultGenParameters())(t)
}

func TestForAllNoShrinkUndecided(t *testing.T) {
	undecided := prop.ForAllNoShrinkT(func(a int) bool {
		return true
	}, gen.Int().SuchThat(func(interface{}) bool {
		return false
	}))
	undecided(gopter.DefaultGenParameters())(t)
}
