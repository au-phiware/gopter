package prop

import (
	"testing"

	"github.com/leanovate/gopter"
)

func convertResultT(result interface{}, err error) func(*testing.T) {
	if err != nil {
		return func(t *testing.T) {
			t.Fatal(err)
		}
	}
	switch r := result.(type) {
	case bool:
		if r {
			return func(t *testing.T) {}
		}
		return func(t *testing.T) { t.Fail() }
	case string:
		if r == "" {
			return func(t *testing.T) {}
		}
		return func(t *testing.T) {
			t.Error(r)
		}
	case *gopter.PropResult:
		return func(t *testing.T) {
			if len(r.Labels) > 0 {
				t.Log(r.Labels)
			}
			if r.Error != nil {
				t.Log(r.Error)
			}
			switch r.Status {
			case gopter.PropFalse:
				t.Fail()
			case gopter.PropError:
				t.FailNow()
			case gopter.PropUndecided:
				t.SkipNow()
			case gopter.PropProof:
				// TODO: t.SkipNow()
			case gopter.PropTrue:
			}
		}
	case func(*testing.T):
		return r
	}
	return func(t *testing.T) {
		t.Fatalf("Invalid check result: %#v", result)
	}
}
