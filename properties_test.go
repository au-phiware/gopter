// +build !fail

package gopter_test

import (
	"testing"
	"regexp"
	"os/exec"
	"fmt"
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

func TestProperties(t *testing.T) {
	out := GoTestOutput(t)
	if found, err := regexp.Match("\\bFAIL: TestProperties/always_fail/#00/#00\\b", out); !found || err != nil {
		t.Error("First test did not fail", err)
	}
	if found, err := regexp.Match("\\bFAIL: TestProperties/always_fail/#00/#01\\b", out); !found || err != nil {
		t.Error("Second test did not fail", err)
	}
}
