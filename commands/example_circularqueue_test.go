package commands_test

import (
	"fmt"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
)

// *****************************************
// Production code (i.e. the implementation)
// *****************************************

type Queue struct {
	inp  int
	outp int
	size int
	buf  []int
}

func New(n int) *Queue {
	return &Queue{
		inp:  0,
		outp: 0,
		size: n + 1,
		buf:  make([]int, n+1),
	}
}

func (q *Queue) Put(n int) int {
	if q.inp == 4 && n > 0 { // Intentional spooky bug
		q.buf[q.size-1] *= n
	}
	q.buf[q.inp] = n
	q.inp = (q.inp + 1) % q.size
	return n
}

func (q *Queue) Get() int {
	ans := q.buf[q.outp]
	q.outp = (q.outp + 1) % q.size
	return ans
}

func (q *Queue) Size() int {
	return (q.inp - q.outp + q.size) % q.size
}

func (q *Queue) Init() {
	q.inp = 0
	q.outp = 0
}

// *****************************************
//               Test code
// *****************************************

// cbState holds the expected state (i.e. its the commands.State)
type cbState struct {
	size         int
	elements     []int
	takenElement int
}

func (st *cbState) TakeFront() {
	st.takenElement = st.elements[0]
	st.elements = append(st.elements[:0], st.elements[1:]...)
}

func (st *cbState) PushBack(value int) {
	st.elements = append(st.elements, value)
}

func (st *cbState) String() string {
	return fmt.Sprintf("State(size=%d, elements=%v)", st.size, st.elements)
}

// Get command simply invokes the Get function on the queue and compares the
// result with the expected state.
var genGetCommand = gen.Const(&commands.ProtoCommand{
	Name: "Get",
	RunFunc: func(q commands.SystemUnderTest) commands.Result {
		return q.(*Queue).Get()
	},
	NextStateFunc: func(state commands.State) commands.State {
		state.(*cbState).TakeFront()
		return state
	},
	// The implementation implicitly assumes that Get is never called on an
	// empty Queue, therefore the command requires a corresponding pre-condition
	PreConditionFunc: func(state commands.State) bool {
		return len(state.(*cbState).elements) > 0
	},
	PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
		if result.(int) != state.(*cbState).takenElement {
			return &gopter.PropResult{Status: gopter.PropFalse}
		}
		return &gopter.PropResult{Status: gopter.PropTrue}
	},
})

// Put command puts a value into the queue by using the Put function. Since
// the Put function has an int argument the Put command should have a
// corresponding parameter.
type putCommand int

func (value putCommand) Run(q commands.SystemUnderTest) commands.Result {
	return q.(*Queue).Put(int(value))
}

func (value putCommand) NextState(state commands.State) commands.State {
	state.(*cbState).PushBack(int(value))
	return state
}

// The implementation implicitly assumes that that Put is never called if
// the capacity is exhausted, therefore the command requires a corresponding
// pre-condition.
func (putCommand) PreCondition(state commands.State) bool {
	s := state.(*cbState)
	return len(s.elements) < s.size
}

func (putCommand) PostCondition(state commands.State, result commands.Result) *gopter.PropResult {
	st := state.(*cbState)
	if result.(int) != st.elements[len(st.elements)-1] {
		return &gopter.PropResult{Status: gopter.PropFalse}
	}
	return &gopter.PropResult{Status: gopter.PropTrue}
}

func (value putCommand) String() string {
	return fmt.Sprintf("Put(%d)", value)
}

// We want to have a generator for put commands for arbitrary int values.
// In this case the command is actually shrinkable, e.g. if the property fails
// by putting a 1000, it might already fail for a 500 as well ...
var genPutCommand = gen.Int().Map(func(value int) commands.Command {
	return putCommand(value)
}).WithShrinker(func(v interface{}) gopter.Shrink {
	return gen.IntShrinker(int(v.(putCommand))).Map(func(value int) putCommand {
		return putCommand(value)
	})
})

// Size command is simple again, it just invokes the Size function and
// compares compares the result with the expected state.
// The Size function can be called any time, therefore this command does not
// require a pre-condition.
var genSizeCommand = gen.Const(&commands.ProtoCommand{
	Name: "Size",
	RunFunc: func(q commands.SystemUnderTest) commands.Result {
		return q.(*Queue).Size()
	},
	PostConditionFunc: func(state commands.State, result commands.Result) *gopter.PropResult {
		if result.(int) != len(state.(*cbState).elements) {
			return &gopter.PropResult{Status: gopter.PropFalse}
		}
		return &gopter.PropResult{Status: gopter.PropTrue}
	},
})

// cbCommands implements the command.Commands interface, i.e. is
// responsible for creating/destroying the system under test and generating
// commands and initial states (cbState)
var cbCommands = &commands.ProtoCommands{
	NewSystemUnderTestFunc: func(initialState commands.State) commands.SystemUnderTest {
		s := initialState.(*cbState)
		q := New(s.size)
		for e := range s.elements {
			q.Put(e)
		}
		return q
	},
	DestroySystemUnderTestFunc: func(sut commands.SystemUnderTest) {
		sut.(*Queue).Init()
	},
	InitialStateGen: gen.IntRange(1, 30).Map(func(size int) *cbState {
		return &cbState{
			size:     size,
			elements: make([]int, 0, size),
		}
	}),
	InitialPreConditionFunc: func(state commands.State) bool {
		s := state.(*cbState)
		return len(s.elements) >= 0 && len(s.elements) <= s.size
	},
	GenCommandFunc: func(state commands.State) gopter.Gen {
		return gen.OneGenOf(genGetCommand, genPutCommand, genSizeCommand)
	},
}

// Kudos to @jamesd for providing this real world example.
// ... of course he did not implemented the bug, that was evil me
//
// The bug only occures on the following conditions:
//  - the queue size has to be greater than 4
//  - the queue has to be filled entirely once
//  - Get operations have to be at least 5 elements behind put
//  - The Put at the end of the queue and 5 elements later have to be non-zero
//
// Lets see what gopter has to say:
//
// The output of this example will be
//  ! circular buffer: Falsified after 96 passed tests.
//  ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0)
//     Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0)
//     Put(0) Put(0) Put(0) Get Get Put(2) Get]
//  ARG_0_ORIGINAL (85 shrinks): initialState=State(size=7, elements=[])
//     sequential=[Put(-1855365712) Put(-1591723498) Get Size Size
//     Put(-1015561691) Get Put(397128011) Size Get Put(1943174048) Size
//     Put(1309500770) Size Get Put(-879438231) Size Get Put(-1644094687) Get
//     Put(-1818606323) Size Put(488620313) Size Put(-1219794505)
//     Put(1166147059) Get Put(11390361) Get Size Put(-1407993944) Get Get Size
//     Put(1393923085) Get Put(1222853245) Size Put(2070918543) Put(1741323168)
//     Size Get Get Size Put(2019939681) Get Put(-170089451) Size Get Get Size
//     Size Put(-49249034) Put(1229062846) Put(642598551) Get Put(1183453167)
//     Size Get Get Get Put(1010460728) Put(6828709) Put(-185198587) Size Size
//     Get Put(586459644) Get Size Put(-1802196502) Get Size Put(2097590857) Get
//     Get Get Get Size Put(-474576011) Size Get Size Size Put(771190414) Size
//     Put(-1509199920) Get Put(967212411) Size Get Put(578995532) Size Get Size
//     Get]
//
// Though this is not the minimal possible combination of command, its already
// pretty close.
func Example_circularqueue() {
	parameters := gopter.DefaultTestParametersWithSeed(1234) // Example should generate reproducable results, otherwise DefaultTestParameters() will suffice

	properties := gopter.NewProperties(parameters)

	properties.Property("circular buffer", commands.Prop(cbCommands))

	// When using testing.T you might just use: properties.RunT(t)
	testing.Main(
		func(a, b string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{
				Name: "Example_circularqueue",
				F:    properties.RunT,
			},
		}, nil, nil)
	// Output:
	//
	// --- FAIL: Example_circularqueue (0.37s)
	// 	--- FAIL: Example_circularqueue/circular_buffer (0.37s)
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Size Size Put(-1015561691) Get Put(397128011) Size Get Put(1943174048) Size Put(1309500770) Size Get Put(-879438231) Size Get Put(-1644094687) Get Put(-1818606323) Size Put(488620313) Size Put(-1219794505) Put(1166147059) Get Put(11390361) Get Size Put(-1407993944) Get Get Size Put(1393923085) Get Put(1222853245) Size Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get Put(-170089451) Size]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#05 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Size Size Put(-1015561691) Get Put(397128011) Size Get Put(1943174048) Size Put(1309500770) Size Get Put(-879438231) Size Get Put(-1644094687) Get Put(-1818606323) Size Put(488620313) Size Put(1222853245) Size Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get Put(-170089451) Size]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#13 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Size Size Put(-1015561691) Get Put(397128011) Put(1309500770) Size Get Put(-879438231) Size Get Put(-1644094687) Get Put(-1818606323) Size Put(488620313) Size Put(1222853245) Size Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get Put(-170089451) Size]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#38 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Size Size Put(-1015561691) Get Put(397128011) Put(1309500770) Size Get Put(-879438231) Size Get Put(-1644094687) Get Put(-1818606323) Size Put(488620313) Size Put(1222853245) Size Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#56 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Size Put(-1015561691) Get Put(397128011) Put(1309500770) Size Get Put(-879438231) Size Get Put(-1644094687) Get Put(-1818606323) Size Put(488620313) Size Put(1222853245) Size Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#74 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Put(-1015561691) Get Put(397128011) Put(1309500770) Size Get Put(-879438231) Size Get Put(-1644094687) Get Put(-1818606323) Size Put(488620313) Size Put(1222853245) Size Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#94 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Put(-1015561691) Get Put(397128011) Put(1309500770) Get Put(-879438231) Size Get Put(-1644094687) Get Put(-1818606323) Size Put(488620313) Size Put(1222853245) Size Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#117 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Put(-1015561691) Get Put(397128011) Put(1309500770) Get Put(-879438231) Get Put(-1644094687) Get Put(-1818606323) Size Put(488620313) Size Put(1222853245) Size Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#140 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Put(-1015561691) Get Put(397128011) Put(1309500770) Get Put(-879438231) Get Put(-1644094687) Get Put(-1818606323) Put(488620313) Size Put(1222853245) Size Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#166 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Put(-1015561691) Get Put(397128011) Put(1309500770) Get Put(-879438231) Get Put(-1644094687) Get Put(-1818606323) Put(488620313) Put(1222853245) Size Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#189 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Put(-1015561691) Get Put(397128011) Put(1309500770) Get Put(-879438231) Get Put(-1644094687) Get Put(-1818606323) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#219 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Put(-1015561691) Get Put(397128011) Put(1309500770) Get Put(-879438231) Get Put(-1644094687) Get Put(-1818606323) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Size Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#250 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Put(-1015561691) Get Put(397128011) Put(1309500770) Get Put(-879438231) Get Put(-1644094687) Get Put(-1818606323) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#285 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(-1591723498) Get Put(-1015561691) Get Put(397128011) Put(1309500770) Get Put(-879438231) Get Put(-1644094687) Get Put(-1818606323) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#320 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(-1015561691) Get Put(397128011) Put(1309500770) Get Put(-879438231) Get Put(-1644094687) Get Put(-1818606323) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#355 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(397128011) Put(1309500770) Get Put(-879438231) Get Put(-1644094687) Get Put(-1818606323) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#390 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(1309500770) Get Put(-879438231) Get Put(-1644094687) Get Put(-1818606323) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#425 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(-879438231) Get Put(-1644094687) Get Put(-1818606323) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#460 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(-1644094687) Get Put(-1818606323) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#495 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1818606323) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#531 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-909303162) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#567 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-454651581) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#603 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-227325791) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#639 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-113662896) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#675 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-56831448) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#711 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-28415724) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#747 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-14207862) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#783 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-7103931) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#819 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-3551966) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#855 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1775983) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#891 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-887992) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#927 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-443996) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#963 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-221998) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#999 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-110999) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1035 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-55500) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1071 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-27750) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1107 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-13875) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1143 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-6938) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1179 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-3469) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1215 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1735) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1251 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-868) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1287 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-434) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1323 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-217) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1359 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-109) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1395 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-55) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1431 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-28) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1467 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-14) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1503 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-7) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1539 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-4) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1575 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-2) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1611 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(488620313) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1647 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(1222853245) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1683 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(2070918543) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1719 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(1741323168) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1755 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(2019939681) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1792 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(1009969841) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1829 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(504984921) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1866 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(252492461) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1903 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(126246231) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1940 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(63123116) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#1977 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(31561558) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2014 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(15780779) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2051 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(7890390) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2088 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(3945195) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2125 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(1972598) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2162 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(986299) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2199 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(493150) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2236 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(246575) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2273 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(123288) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2310 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(61644) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2347 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(30822) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2384 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(15411) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2421 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(7706) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2458 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(3853) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2495 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(1927) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2532 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(964) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2569 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(482) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2606 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(241) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2643 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(121) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2680 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(61) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2717 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(31) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2754 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(16) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2791 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(8) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2828 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(4) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/shrink_0#2865 (0.00s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(0) Put(0) Get Put(0) Get Put(0) Put(0) Get Put(0) Get Put(0) Get Put(-1) Put(0) Put(0) Put(0) Put(0) Get Get Put(2) Get]
	// 		--- FAIL: Example_circularqueue/circular_buffer/original#96 (0.27s)
	// 			check_condition_func.go:39: ARG_0: initialState=State(size=7, elements=[]) sequential=[Put(-1855365712) Put(-1591723498) Get Size Size Put(-1015561691) Get Put(397128011) Size Get Put(1943174048) Size Put(1309500770) Size Get Put(-879438231) Size Get Put(-1644094687) Get Put(-1818606323) Size Put(488620313) Size Put(-1219794505) Put(1166147059) Get Put(11390361) Get Size Put(-1407993944) Get Get Size Put(1393923085) Get Put(1222853245) Size Put(2070918543) Put(1741323168) Size Get Get Size Put(2019939681) Get Put(-170089451) Size Get Get Size Size Put(-49249034) Put(1229062846) Put(642598551) Get Put(1183453167) Size Get Get Get Put(1010460728) Put(6828709) Put(-185198587) Size Size Get Put(586459644) Get Size Put(-1802196502) Get Size Put(2097590857) Get Get Get Get Size Put(-474576011) Size Get Size Size Put(771190414) Size Put(-1509199920) Get Put(967212411) Size Get Put(578995532) Size Get Size Get]
	// 		forall.go:64: Falsified after 96 passed tests.
	// 		runner.go:71: Elapsed time: 372.877117ms
	// 		runner.go:72: Completed with seed: 1234
	// FAIL
}
