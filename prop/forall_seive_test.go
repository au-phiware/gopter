package prop_test

import (
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestSeive(t *testing.T) {
	type shrinkResult struct {
		value interface{}
		ok bool
	}
	shinkResults := []shrinkResult{
		shrinkResult{ shrinkResult{ 1, true }, true },
		shrinkResult{ nil, true },
	}
	propResults := []*gopter.PropResult{
		gopter.NewPropResult(false, ""),
		gopter.NewPropResult(false, ""),
		nil,
	}
	i := 0
	j := 0
	p := prop.ForAll(
		func(ptr *shrinkResult) *gopter.PropResult {
			if ptr == nil {
				t.Errorf("Nil passed seive")
			}
			t.Logf("condition: %#v", ptr)
			defer func() { j++ }()
			return propResults[j]
		},
		gen.PtrOf(
			gen.Struct(reflect.TypeOf(shrinkResult{}), nil).
			WithShrinker(
				func(value interface{}) gopter.Shrink {
					return func() (interface{}, bool) {
						defer func() { i++ }()
						if i < len(shinkResults) {
							t.Logf("shrink: %#v, %#v", shinkResults[i].value, shinkResults[i].ok)
							return shinkResults[i].value, shinkResults[i].ok
						}
						return nil, false
					}
				},
			),
		).
		SuchThat(
			func(ptr interface{}) bool {
				t.Logf("such that: %#v", ptr)
				return ptr != nil
			},
		),
	)

	p(gopter.DefaultGenParameters().CloneWithSeed(0xdeadbeef))
}
