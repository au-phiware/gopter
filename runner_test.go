package gopter

import (
	"sync/atomic"
	"testing"
)

func TestRunnerSingleWorker(t *testing.T) {
	parameters := DefaultTestParameters()
	testRunner := &runner{
		parameters: parameters,
		worker: func(num int, shouldStop shouldStop) func(*testing.T) {
			return func(*testing.T) {}
		},
	}

	testRunner.runWorkersT()(t)
}

func TestRunnerParallelWorkers(t *testing.T) {
	parameters := DefaultTestParameters()
	specs := []struct {
		workers int
		res     []func(*testing.T)
		exp     func(t, u *testing.T)
	}{
		// Test all pass
		{
			workers: 50,
			res: []func(*testing.T) {
				func(*testing.T) {},
			},
			exp: func(t, u *testing.T) {
				if u.Failed() || u.Skipped() {
					t.Error("Test did not pass")
				}
			},
		},
		// Test all fail
		{
			workers: 50,
			res: []func(*testing.T) {
				func(t *testing.T) { t.Fail() },
			},
			exp: func(t, u *testing.T) {
				if !u.Failed() || u.Skipped() {
					t.Error("Test did not fail")
				}
			},
		},
		// a pass and failure
		{
			workers: 2,
			res: []func(*testing.T) {
				func(*testing.T) {},
				func(t *testing.T) { t.Fail() },
			},
			exp: func(t, u *testing.T) {
				if !u.Failed() || u.Skipped() {
					t.Error("Test did not fail")
				}
			},
		},
		// a pass and multiple failures (first failure returned)
		{
			workers: 3,
			res: []func(*testing.T) {
				func(*testing.T) {},
				func(t *testing.T) { t.Fail() },
				func(t *testing.T) { t.Fail() },
			},
			exp: func(t, u *testing.T) {
				if !u.Failed() || u.Skipped() {
					t.Error("Test did not fail")
				}
			},
		},
	}

	for _, spec := range specs {
		parameters.Workers = spec.workers

		var called int64
		testRunner := &runner{
			parameters: parameters,
			worker: func(num int, shouldStop shouldStop) func(*testing.T) {
				return func(t *testing.T) {
					atomic.AddInt64(&called, 1)

					if num < len(spec.res) {
						spec.res[num](t)
					} else {
						spec.res[0](t)
					}
				}
			},
		}

		fakeT := &testing.T{}
		testRunner.runWorkersT()(fakeT)

		if called != int64(spec.workers) {
			t.Errorf("Not enough calls; want %d; got %d", spec.workers, called)
		}
		spec.exp(t, fakeT)
	}
}
