package gopter

import (
	"fmt"
	"math"
	"runtime/debug"
	"testing"
)

// Prop represent some kind of property that (drums please) can and should be checked
type Prop func(*GenParameters) func(*testing.T)

// SaveProp creates s save property by handling all panics from an inner property
func SaveProp(prop Prop) Prop {
	return func(genParams *GenParameters) (result func(*testing.T)) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				result = func(t *testing.T) {
					t.Fatalf("Check paniced: %v %s", r, stack)
				}
			}
		}()

		return prop(genParams)
	}
}

// Check the property using specific parameters
func (prop Prop) Check(t *testing.T, parameters *TestParameters) *TestResult {
	iterations := math.Ceil(float64(parameters.MinSuccessfulTests) / float64(parameters.Workers))
	sizeStep := float64(parameters.MaxSize-parameters.MinSize) / (iterations * float64(parameters.Workers))

	genParameters := GenParameters{
		MinSize:        parameters.MinSize,
		MaxSize:        parameters.MaxSize,
		MaxShrinkCount: parameters.MaxShrinkCount,
		Rng:            parameters.Rng,
	}
	runner := &runner{
		parameters: parameters,
		worker: func(workerIdx int, shouldStop shouldStop) *TestResult {
			testResult := &TestResult{
				Status: TestPassed,
			}

			isExhaused := func() bool {
				return testResult.n+testResult.n > parameters.MinSuccessfulTests &&
					1.0+float64(parameters.Workers*testResult.n)*parameters.MaxDiscardRatio < float64(testResult.n)
			}

			for !shouldStop() &&
				testResult.n < int(iterations) &&
				testResult.Status == TestPassed {
				size := float64(parameters.MinSize) + (sizeStep * float64(workerIdx+(parameters.Workers*(testResult.n+testResult.n))))
				runner := prop(genParameters.WithSize(int(size)))

				// TODO: guard against t.Parallel()
				t.Run("", func(t *testing.T) {
					defer func() {
						if t.Failed() {
							testResult.Status = TestFailed
						}
						if t.Skipped() {
							testResult.d++
							if isExhaused() {
								testResult.Status = TestExhausted
							}
						} else {
							testResult.n++
						}
						// TODO: how to establish Proof?
					}()

					runner(t)
				})
			}

			if isExhaused() {
				testResult.Status = TestExhausted
			}
			return
		},
	}

	return runner.runWorkers()
}
