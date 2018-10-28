package prop

import (
	"regexp"
	"errors"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
)

func OutputMatches(t *testing.T, u *testing.T, pattern string) bool {
	output := reflect.ValueOf(*u).FieldByName("output").Bytes()
	found, err := regexp.Match(pattern, output)
	if err != nil {
		t.Fatal(err)
	}
	return found
}

func TestConvertResult(t *testing.T) {
	fakeT := &testing.T{}
	convertResultT(true, nil)(fakeT)
	if fakeT.Failed() || fakeT.Skipped() {
		t.Errorf("Invalid true result: %#v", fakeT)
	}

	fakeT = &testing.T{}
	convertResultT(false, nil)(fakeT)
	if !fakeT.Failed() || fakeT.Skipped() {
		t.Errorf("Invalid false result: %#v", fakeT)
	}

	fakeT = &testing.T{}
	convertResultT("", nil)(fakeT)
	if fakeT.Failed() || fakeT.Skipped() {
		t.Errorf("Invalid string true result: %#v", fakeT)
	}

	fakeT = &testing.T{}
	convertResultT("Something is wrong", nil)(fakeT)
	if !fakeT.Failed() || fakeT.Skipped() {
		t.Errorf("Invalid string false result: %#v", fakeT)
	}
	if !OutputMatches(t, fakeT, "\\bSomething is wrong\\b") {
		t.Errorf("Result did not report message: %#v", fakeT)
	}

	done := make(chan struct{})
	go func() {
		fakeT = &testing.T{}
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
		convertResultT("Anthing", errors.New("Booom"))(fakeT)
		t.Errorf("Test runner did not panic")
	}()
	<-done

	fakeT = &testing.T{}
	convertResultT(&gopter.PropResult{
		Status: gopter.PropProof,
	}, nil)(fakeT)
	if fakeT.Failed() || fakeT.Skipped() {
		t.Errorf("Invalid true result: %#v", fakeT)
	}

	done = make(chan struct{})
	go func() {
		fakeT = &testing.T{}
		defer func() {
			if err := recover(); err != nil {
				t.Errorf("Invalid error result: %#v", err)
			}
			if !fakeT.Failed() || fakeT.Skipped() {
				t.Errorf("Invalid error result: %#v", fakeT)
			}
			if !OutputMatches(t, fakeT, "\\bInvalid check result: 0\\b") {
				t.Errorf("Result did not report message: %#v", fakeT)
			}
			done<-struct{}{}
		}()
		convertResultT(0, nil)(fakeT)
		t.Errorf("Test runner did not panic")
	}()
	<-done
}
