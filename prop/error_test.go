package prop_test

import (
	"regexp"
	"errors"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/prop"
)

func OutputMatches(t *testing.T, u *testing.T, pattern string) bool {
	output := reflect.ValueOf(*u).FieldByName("output").Bytes()
	found, err := regexp.Match(pattern, output)
	if err != nil {
		t.Fatal(err)
	}
	return found
}

func TestErrorProp(t *testing.T) {
	done := make(chan struct{})
	go func() {
		fakeT := &testing.T{}
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Invalid error result: %#v", err)
			}
			if !fakeT.Failed() || fakeT.Skipped() {
				t.Errorf("Invalid error result: %#v", fakeT)
			}
			if !OutputMatches(t, fakeT, "\\bBooom\\b") {
				t.Errorf("Result did not report message: %#v", fakeT)
			}
			done<-struct{}{}
		}()
		p := prop.ErrorPropT(errors.New("Booom"))
		p(gopter.DefaultGenParameters())(fakeT)

		t.Errorf("Test runner did not panic")
	}()
	<-done
}
