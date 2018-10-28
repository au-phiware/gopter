// +build fail

package prop_test

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)


func TestForAllInvalidParam(t *testing.T) {
	fail := prop.ForAllT(0)
	fail(gopter.DefaultGenParameters())(t)
}

func TestForAllUndecided(t *testing.T) {
	undecided := prop.ForAllT(func(a int) bool {
		return true
	}, gen.Int().SuchThat(func(interface{}) bool {
		return false
	}))
	undecided(gopter.DefaultGenParameters())(t)
}
