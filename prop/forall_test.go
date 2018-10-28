// +build !fail

package prop_test

import (
	"math"
	"regexp"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func TestSqrt(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("greater one of all greater one", prop.ForAllT(
		func(v float64) bool {
			return math.Sqrt(v) >= 1
		},
		gen.Float64Range(1, math.MaxFloat64),
	))

	properties.Property("squared is equal to value", prop.ForAllT(
		func(v float64) bool {
			r := math.Sqrt(v)
			return math.Abs(r*r-v) < 1e-10*v
		},
		gen.Float64Range(0, math.MaxFloat64),
	))

	properties.RunT(t)
}

func TestForAllInvalidParam(t *testing.T) {
	if found, err := regexp.Match("\\bFirst param of ForrAll has to be a func: int\\b", GoTestOutput(t, "FAIL")); !found || err != nil {
		t.Error("Failed to panic", err)
	}
}

func TestForAllUndecided(t *testing.T) {
	if found, err := regexp.Match("\\bforall.go\\b", GoTestOutput(t, "SKIP")); !found || err != nil {
		t.Error("Failed to panic", err)
	}
}
