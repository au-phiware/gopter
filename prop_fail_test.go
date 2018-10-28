// +build fail

package gopter_test

import (
	"sync/atomic"
	"testing"

	"github.com/leanovate/gopter"
)

func TestSaveProp(t *testing.T) {
	prop := gopter.SavePropT(func(*gopter.GenParameters) func(*testing.T) {
		panic("Ouchy")
	})

	prop.CheckT(t)
}

func TestSavePropNested(t *testing.T) {
	prop := gopter.SavePropT(func(*gopter.GenParameters) func(*testing.T) {
		return func(*testing.T) {
			panic("Ouchy")
		}
	})

	prop.CheckT(t)
}

func TestPropFalse(t *testing.T) {
	var called int64
	prop := gopter.PropT(func(genParams *gopter.GenParameters) func(*testing.T) {
		atomic.AddInt64(&called, 1)

		return func(t *testing.T) {
			t.Fail()
		}
	})

	parameters := gopter.DefaultTestParameters()
	t.Run("", prop.CheckWithParametersT(parameters))

	t.Logf("number of calls: %d", called)
}

func TestPropError(t *testing.T) {
	var called int64
	prop := gopter.PropT(func(genParams *gopter.GenParameters) func(*testing.T) {
		atomic.AddInt64(&called, 1)

		return func(t *testing.T) {
			t.FailNow()
		}
	})

	parameters := gopter.DefaultTestParameters()
	t.Run("", prop.CheckWithParametersT(parameters))

	t.Logf("number of calls: %d", called)
}
