package prop

import "testing"
import "github.com/leanovate/gopter"

// ErrorProp creates a property that will always fail with an error.
// Mostly used as a fallback when setup/initialization fails
func ErrorPropT(err error) gopter.PropT {
	return func(genParams *gopter.GenParameters) func(*testing.T) {
		return func(t *testing.T) {
			t.Fatal(err)
		}
	}
}
