package gopter

import (
	"sync"
	"testing"
	"time"
)

type shouldStop func() bool

type worker func(int, shouldStop) func(*testing.T)

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
	if r.parameters.Workers < 2 {
		runner := r.worker(0, func() (stop bool) { return })
		return func(t *testing.T) {
			start := time.Now()
			runner(t)
			t.Logf("Elapsed time: %s", time.Since(start))
		}
	}
	var waitGroup sync.WaitGroup
	waitGroup.Add(r.parameters.Workers)

	return func(t *testing.T) {
		var stopFlag Flag
		defer stopFlag.Set()
		start := time.Now()
		for i := 0; i < r.parameters.Workers; i++ {
			go func(runner func(*testing.T)) {
				defer waitGroup.Done()
				runner(t)
			}(r.worker(i, stopFlag.Get))
		}
		waitGroup.Wait()
		t.Logf("Elapsed time: %s", time.Since(start))
	}
}
