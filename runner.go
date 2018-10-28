package gopter

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

type shouldStop func() bool

type worker func(string, int, shouldStop) func(*testing.T)

type runner struct {
	sync.RWMutex
	parameters *TestParameters
	worker     worker
}

func (r *runner) mergeCheckResults(r1, r2 *TestResult) *TestResult {
	var result TestResult

	switch {
	case r1 == nil:
		return r2
	case r1.Status != TestPassed && r1.Status != TestExhausted:
		result = *r1
	case r2.Status != TestPassed && r2.Status != TestExhausted:
		result = *r2
	default:
		result.Status = TestExhausted

		if r1.Succeeded+r2.Succeeded >= r.parameters.MinSuccessfulTests &&
			float64(r1.Discarded+r2.Discarded) <= float64(r1.Succeeded+r2.Succeeded)*r.parameters.MaxDiscardRatio {
			result.Status = TestPassed
		}
	}

	result.Succeeded = r1.Succeeded + r2.Succeeded
	result.Discarded = r1.Discarded + r2.Discarded

	return &result
}

func (r *runner) runWorkersT() func(*testing.T) {
	var runner func(*testing.T)
	if r.parameters.Workers < 2 {
		runner = r.worker("", 0, func() (stop bool) { return })
	} else {
		var waitGroup sync.WaitGroup
		waitGroup.Add(r.parameters.Workers)
		namef := fmt.Sprintf("worker%%0%dd", len(strconv.Itoa(r.parameters.Workers)))

		runner = func(t *testing.T) {
			t.Helper()
			var stopFlag Flag
			defer stopFlag.Set()
			for i := 0; i < r.parameters.Workers; i++ {
				go func(runner func(*testing.T)) {
					defer waitGroup.Done()
					runner(t)
				}(r.worker(fmt.Sprintf(namef, i), i, stopFlag.Get))
			}
			waitGroup.Wait()
		}
	}
	return func(t *testing.T) {
		t.Helper()
		defer func(start time.Time) {
			t.Logf("Elapsed time: %s", time.Since(start))
			t.Logf("Completed with seed: %d", r.parameters.Seed)
		}(time.Now())
		runner(t)
	}
}
