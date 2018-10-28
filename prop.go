package gopter

import (
	"math"
	"runtime/debug"
	"testing"
)

// Prop represent some kind of property that (drums please) can and should be checked
type PropT func(*GenParameters) func(*testing.T)

// SaveProp creates s save property by handling all panics from an inner property
func SavePropT(prop PropT) PropT {
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

func RunT(t *testing.T, runner func(t *testing.T), recovery func(t *testing.T)) {
	// TODO: guard against t.Parallel()
	t.Run("", func(t *testing.T) {
		defer recovery(t)
		runner(t)
	})
}

func (prop PropT) CheckT(t *testing.T) {
	prop.CheckWithParametersT(DefaultTestParameters())(t)
}

// Check the property using specific parameters
func (prop PropT) CheckWithParametersT(parameters *TestParameters) func(*testing.T) {
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
		worker: func(workerIdx int, shouldStop shouldStop) func(*testing.T) {
			var n int
			var d int
			var status testStatus

			isExhaused := func() bool {
				return n+d > parameters.MinSuccessfulTests &&
					1.0+float64(parameters.Workers*n)*parameters.MaxDiscardRatio < float64(d)
			}

			return func(t *testing.T) {
				for !shouldStop() &&
					n < int(iterations) &&
					status == TestPassed {
					size := float64(parameters.MinSize) + (sizeStep * float64(workerIdx+(parameters.Workers*(n+d))))
					runner := prop(genParameters.WithSize(int(size)))

					RunT(t, runner, func(t *testing.T) {
						if r := recover(); r != nil {
							t.Errorf("Check paniced: %v %s", r, debug.Stack())
						}
						if t.Failed() {
							status = TestFailed
						} else if t.Skipped() {
							d++
							if isExhaused() {
								status = TestExhausted
							}
						} else {
							n++
						}
						// TODO: how to establish Proof?
					})
				}

				if isExhaused() {
					status = TestExhausted
				}
			}
		},
	}
	return runner.runWorkersT()
}
