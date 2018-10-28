// +build !fail

package gopter

import (
	"fmt"
	"regexp"
	"os/exec"
	"sync/atomic"
	"testing"
)

func GoTestOutput(t testing.TB) []byte {
	cmd := exec.Command("go", "test", "-v", "-tags=fail", "-run", t.Name())
	out, err := cmd.CombinedOutput()
	if _, exit := err.(*exec.ExitError); !exit {
		t.Fatal(err)
	}
	if err == nil {
		t.Error("go test should have exited with status code 1")
	}
	if found, err := regexp.Match(fmt.Sprintf("\\bFAIL: %s\\b", t.Name()), out); !found || err != nil {
		t.Error("Test did not fail", err)
	}
	return out
}

func TestSaveProp(t *testing.T) {
	if found, err := regexp.Match("\\bCheck paniced: Ouchy\\b", GoTestOutput(t)); !found || err != nil {
		t.Error("Failed to panic", err)
	}
}

func TestSavePropNested(t *testing.T) {
	if found, err := regexp.Match("\\bCheck paniced: Ouchy\\b", GoTestOutput(t)); !found || err != nil {
		t.Error("Failed to panic", err)
	}
}

func TestPropUndecided(t *testing.T) {
	var called int64
	prop := PropT(func(genParams *GenParameters) func(*testing.T) {
		atomic.AddInt64(&called, 1)

		return func(t *testing.T) {
			t.SkipNow()
		}
	})

	parameters := DefaultTestParameters()
	t.Run("", func(s *testing.T) {
		defer func() {
			if s.Skipped() {
				t.Error("Test was skipped")
			}
		}()
		prop.CheckWithParametersT(parameters)(s)
	})

	if called != int64(parameters.MinSuccessfulTests)+1 {
		t.Errorf("Invalid number of calls: %d", called)
	}
}

func TestPropMaxDiscardRatio(t *testing.T) {
	var called int64
	prop := PropT(func(genParams *GenParameters) func(*testing.T) {
		atomic.AddInt64(&called, 1)

		if genParams.MaxSize > 21 {
			return func(t *testing.T) {}
		}
		return func(t *testing.T) {
			t.SkipNow()
		}
	})

	parameters := DefaultTestParameters()
	parameters.MaxDiscardRatio = 0.2
	t.Run("with T", prop.CheckWithParametersT(parameters))

	if called != int64(parameters.MinSuccessfulTests)+22 {
		t.Errorf("Invalid number of calls: %d", called)
	}
}

func TestPropPassed(t *testing.T) {
	var called int64
	prop := PropT(func(genParams *GenParameters) func(*testing.T) {
		atomic.AddInt64(&called, 1)

		return func(t *testing.T) {}
	})

	parameters := DefaultTestParameters()
	prop.CheckWithParametersT(parameters)(t)

	if called != int64(parameters.MinSuccessfulTests) {
		t.Errorf("Invalid number of calls: %d", called)
	}
}

func TestPropFalse(t *testing.T) {
	if found, err := regexp.Match("\\bnumber of calls: 1\\b", GoTestOutput(t)); !found || err != nil {
		t.Error("Wrong number of calls", err)
	}
}

func TestPropError(t *testing.T) {
	if found, err := regexp.Match("\\bnumber of calls: 1\\b", GoTestOutput(t)); !found || err != nil {
		t.Error("Wrong number of calls", err)
	}
}

func TestPropPassedMulti(t *testing.T) {
	var called int64
	prop := PropT(func(genParams *GenParameters) func(*testing.T) {
		atomic.AddInt64(&called, 1)

		return func(t *testing.T) {}
	})

	parameters := DefaultTestParameters()
	parameters.Workers = 10
	prop.CheckWithParametersT(parameters)(t)

	if called != int64(parameters.MinSuccessfulTests) {
		t.Errorf("Invalid number of calls: %d", called)
	}
}
