package prop_test

import (
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func Example_timeGen() {
	parameters := gopter.DefaultTestParametersWithSeed(1234) // Example should generate reproducable results, otherwise DefaultTestParameters() will suffice
	time.Local = time.UTC                                    // Just for this example to generate reproducable results

	properties := gopter.NewProperties(parameters)

	properties.Property("time in range format parsable", prop.ForAllT(
		func(actual time.Time) (bool, error) {
			str := actual.Format(time.RFC3339Nano)
			parsed, err := time.Parse(time.RFC3339Nano, str)
			return actual.Equal(parsed), err
		},
		gen.TimeRange(time.Now(), time.Duration(100*24*365)*time.Hour).WithLabel("actual"),
	))

	properties.Property("regular time format parsable", prop.ForAllT(
		func(actual time.Time) (bool, error) {
			str := actual.Format(time.RFC3339Nano)
			parsed, err := time.Parse(time.RFC3339Nano, str)
			return actual.Equal(parsed), err
		},
		gen.Time().WithLabel("actual"),
	))

	properties.Property("any time format parsable", prop.ForAllT(
		func(actual time.Time) (bool, error) {
			str := actual.Format(time.RFC3339Nano)
			parsed, err := time.Parse(time.RFC3339Nano, str)
			return actual.Equal(parsed), err
		},
		gen.AnyTime().WithLabel("actual"),
	))

	testing.Main(
		func(a, b string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{
				Name: "Example_timeGen",
				F:    properties.RunT,
			},
		}, nil, nil)
	// Output:
	//
	// --- FAIL: Example_timeGen (0.02s)
	// 	--- FAIL: Example_timeGen/any_time_format_parsable (0.01s)
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#01 (0.00s)
	// 			prop.go:37: parsing time "237903042092-02-10T19:15:18Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "03042092-02-10T19:15:18Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 237903042092-02-10 19:15:18 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#03 (0.00s)
	// 			prop.go:37: parsing time "118951522031-01-21T09:37:39Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "51522031-01-21T09:37:39Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 118951522031-01-21 09:37:39 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#05 (0.00s)
	// 			prop.go:37: parsing time "59475762000-07-12T04:48:50Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "5762000-07-12T04:48:50Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 59475762000-07-12 04:48:50 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#07 (0.00s)
	// 			prop.go:37: parsing time "29737881985-04-07T02:24:25Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "7881985-04-07T02:24:25Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 29737881985-04-07 02:24:25 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#09 (0.00s)
	// 			prop.go:37: parsing time "14868941977-08-19T13:12:13Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "8941977-08-19T13:12:13Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 14868941977-08-19 13:12:13 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#11 (0.00s)
	// 			prop.go:37: parsing time "7434471973-10-25T18:36:07Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "471973-10-25T18:36:07Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 7434471973-10-25 18:36:07 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#13 (0.00s)
	// 			prop.go:37: parsing time "3717236971-11-28T09:18:04Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "236971-11-28T09:18:04Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 3717236971-11-28 09:18:04 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#15 (0.00s)
	// 			prop.go:37: parsing time "1858619470-12-15T04:39:02Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "619470-12-15T04:39:02Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 1858619470-12-15 04:39:02 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#17 (0.00s)
	// 			prop.go:37: parsing time "929310720-06-24T02:19:31Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "10720-06-24T02:19:31Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 929310720-06-24 02:19:31 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#19 (0.00s)
	// 			prop.go:37: parsing time "464656345-03-29T01:09:46Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "56345-03-29T01:09:46Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 464656345-03-29 01:09:46 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#21 (0.00s)
	// 			prop.go:37: parsing time "232329157-08-15T00:34:53Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "29157-08-15T00:34:53Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 232329157-08-15 00:34:53 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#23 (0.00s)
	// 			prop.go:37: parsing time "116165563-10-24T00:17:27Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "65563-10-24T00:17:27Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 116165563-10-24 00:17:27 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#25 (0.00s)
	// 			prop.go:37: parsing time "58083766-11-27T00:08:44Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "3766-11-27T00:08:44Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 58083766-11-27 00:08:44 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#27 (0.00s)
	// 			prop.go:37: parsing time "29042868-06-14T00:04:22Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "2868-06-14T00:04:22Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 29042868-06-14 00:04:22 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#29 (0.00s)
	// 			prop.go:37: parsing time "14522419-03-24T12:02:11Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "2419-03-24T12:02:11Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 14522419-03-24 12:02:11 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#31 (0.00s)
	// 			prop.go:37: parsing time "7262194-08-12T06:01:06Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "194-08-12T06:01:06Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 7262194-08-12 06:01:06 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#33 (0.00s)
	// 			prop.go:37: parsing time "3632082-04-22T03:00:33Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "082-04-22T03:00:33Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 3632082-04-22 03:00:33 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#35 (0.00s)
	// 			prop.go:37: parsing time "1817026-02-26T01:30:17Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "026-02-26T01:30:17Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 1817026-02-26 01:30:17 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#37 (0.00s)
	// 			prop.go:37: parsing time "909498-01-28T12:45:09Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "98-01-28T12:45:09Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 909498-01-28 12:45:09 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#39 (0.00s)
	// 			prop.go:37: parsing time "455734-01-14T18:22:35Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "34-01-14T18:22:35Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 455734-01-14 18:22:35 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#41 (0.00s)
	// 			prop.go:37: parsing time "228852-01-07T21:11:18Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "52-01-07T21:11:18Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 228852-01-07 21:11:18 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#43 (0.00s)
	// 			prop.go:37: parsing time "115411-01-04T22:35:39Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "11-01-04T22:35:39Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 115411-01-04 22:35:39 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#45 (0.00s)
	// 			prop.go:37: parsing time "58690-07-03T23:17:50Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-07-03T23:17:50Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 58690-07-03 23:17:50 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#47 (0.00s)
	// 			prop.go:37: parsing time "30330-04-03T11:38:55Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-04-03T11:38:55Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 30330-04-03 11:38:55 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#49 (0.00s)
	// 			prop.go:37: parsing time "16150-02-15T17:49:28Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-02-15T17:49:28Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 16150-02-15 17:49:28 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#52 (0.00s)
	// 			prop.go:37: parsing time "12605-02-04T13:22:06Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "5-02-04T13:22:06Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 12605-02-04 13:22:06 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#56 (0.00s)
	// 			prop.go:37: parsing time "11275-09-15T23:41:51Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "5-09-15T23:41:51Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 11275-09-15 23:41:51 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#60 (0.00s)
	// 			prop.go:37: parsing time "10112-06-29T23:44:08Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "2-06-29T23:44:08Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10112-06-29 23:44:08 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#68 (0.00s)
	// 			prop.go:37: parsing time "10048-11-17T17:33:01Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "8-11-17T17:33:01Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10048-11-17 17:33:01 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#77 (0.00s)
	// 			prop.go:37: parsing time "10017-04-28T08:40:10Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "7-04-28T08:40:10Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10017-04-28 08:40:10 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#87 (0.00s)
	// 			prop.go:37: parsing time "10001-08-09T16:31:40Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "1-08-09T16:31:40Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10001-08-09 16:31:40 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#101 (0.00s)
	// 			prop.go:37: parsing time "10000-08-16T14:20:15Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-08-16T14:20:15Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-08-16 14:20:15 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#116 (0.00s)
	// 			prop.go:37: parsing time "10000-02-19T13:46:01Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-02-19T13:46:01Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-02-19 13:46:01 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#133 (0.00s)
	// 			prop.go:37: parsing time "10000-01-05T19:41:24Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-05T19:41:24Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-05 19:41:24 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#154 (0.00s)
	// 			prop.go:37: parsing time "10000-01-03T00:33:41Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-03T00:33:41Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-03 00:33:41 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#176 (0.00s)
	// 			prop.go:37: parsing time "10000-01-01T14:59:50Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-01T14:59:50Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-01 14:59:50 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#200 (0.00s)
	// 			prop.go:37: parsing time "10000-01-01T06:36:23Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-01T06:36:23Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-01 06:36:23 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#225 (0.00s)
	// 			prop.go:37: parsing time "10000-01-01T02:24:40Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-01T02:24:40Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-01 02:24:40 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#251 (0.00s)
	// 			prop.go:37: parsing time "10000-01-01T00:18:49Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-01T00:18:49Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-01 00:18:49 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#280 (0.00s)
	// 			prop.go:37: parsing time "10000-01-01T00:03:06Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-01T00:03:06Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-01 00:03:06 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#312 (0.00s)
	// 			prop.go:37: parsing time "10000-01-01T00:01:09Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-01T00:01:09Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-01 00:01:09 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#345 (0.00s)
	// 			prop.go:37: parsing time "10000-01-01T00:00:11Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-01T00:00:11Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-01 00:00:11 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#381 (0.00s)
	// 			prop.go:37: parsing time "10000-01-01T00:00:04Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-01T00:00:04Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-01 00:00:04 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#418 (0.00s)
	// 			prop.go:37: parsing time "10000-01-01T00:00:01Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-01T00:00:01Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-01 00:00:01 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/shrink_actual#456 (0.00s)
	// 			prop.go:37: parsing time "10000-01-01T00:00:00Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "0-01-01T00:00:00Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 10000-01-01 00:00:00 +0000 UTC
	// 		--- FAIL: Example_timeGen/any_time_format_parsable/original (0.01s)
	// 			prop.go:37: parsing time "237903042092-02-10T19:15:18.148265469Z" as "2006-01-02T15:04:05.999999999Z07:00": cannot parse "03042092-02-10T19:15:18.148265469Z" as "-"
	// 			check_condition_func.go:39: ARG_0: 237903042092-02-10 19:15:18.148265469 +0000 UTC
	// 		forall.go:64: Falsified after 0 passed tests.
	// 		runner.go:72: Completed with seed: 1234
	// FAIL
}
