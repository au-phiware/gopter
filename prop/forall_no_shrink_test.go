// +build !fail

package prop_test

import (
	"fmt"
	"regexp"
	"os/exec"
	"testing"
)

func GoTestOutput(t testing.TB, result string) []byte {
	cmd := exec.Command("go", "test", "-v", "-tags=fail", "-run", t.Name())
	out, err := cmd.CombinedOutput()
	if _, exit := err.(*exec.ExitError); !exit && result == "FAIL" {
		t.Fatal(err)
	}
	if err == nil && result == "FAIL" {
		t.Error("go test should have exited with status code 1")
	}
	if found, err := regexp.Match(fmt.Sprintf("\\b%s: %s\\b", result, t.Name()), out); !found || err != nil {
		t.Error("Test did not fail", err)
	}
	return out
}

func TestForAllNoShrinkInvalidParam(t *testing.T) {
	if found, err := regexp.Match("\\bFirst param of ForrAll has to be a func: int\\b", GoTestOutput(t, "FAIL")); !found || err != nil {
		t.Error("Failed to panic", err)
	}
}

func TestForAllNoShrinkUndecided(t *testing.T) {
	if found, err := regexp.Match("\\bforall_no_shrink.go\\b", GoTestOutput(t, "SKIP")); !found || err != nil {
		t.Error("Failed to panic", err)
	}
}
