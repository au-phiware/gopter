/*
Package gopter contain the main interfaces of the GOlang Property TestER.

A simple property test might look like this:

    func TestSqrt(t *testing.T) {
	t.Run("greater one of all greater one", prop.CheckForAll(
    		func(v float64) bool {
    			return math.Sqrt(v) >= 1
    		},
    		gen.Float64Range(1, math.MaxFloat64),
    	))

    	t.Run("squared is equal to value", prop.CheckForAll(
    		func(v float64) bool {
    			r := math.Sqrt(v)
    			return math.Abs(r*r-v) < 1e-10*v
    		},
    		gen.Float64Range(0, math.MaxFloat64),
    	))
    }

Generally a property is just a function that takes GenParameters and produces
a testing func:

    type Prop func(*GenParameters) func(*testing.T)

but usually you will use prop.ForAll, prop.ForAllNoShrink or arbitrary.ForAll.
There is also the commands package, which can be helpful for stateful testing.
*/
package gopter
