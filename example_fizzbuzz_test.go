package gopter_test

import (
	"testing"
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

// Fizzbuzz: See https://wikipedia.org/wiki/Fizz_buzz
func fizzbuzz(number int) (string, error) {
	if number <= 0 {
		return "", errors.New("Undefined")
	}
	switch {
	case number%15 == 0:
		return "FizzBuzz", nil
	case number%3 == 0:
		return "Fizz", nil
	case number%5 == 0:
		return "Buzz", nil
	}
	return strconv.Itoa(number), nil
}

func Example_fizzbuzz() {
	properties := gopter.NewProperties(nil)

	properties.Property("Undefined for all <= 0", prop.ForAllT(
		func(number int) bool {
			result, err := fizzbuzz(number)
			return err != nil && result == ""
		},
		gen.IntRange(math.MinInt32, 0),
	))

	properties.Property("Start with Fizz for all multiples of 3", prop.ForAllT(
		func(i int) bool {
			result, err := fizzbuzz(i * 3)
			return err == nil && strings.HasPrefix(result, "Fizz")
		},
		gen.IntRange(1, math.MaxInt32/3),
	))

	properties.Property("End with Buzz for all multiples of 5", prop.ForAllT(
		func(i int) bool {
			result, err := fizzbuzz(i * 5)
			return err == nil && strings.HasSuffix(result, "Buzz")
		},
		gen.IntRange(1, math.MaxInt32/5),
	))

	properties.Property("Int as string for all non-divisible by 3 or 5", prop.ForAllT(
		func(number int) bool {
			result, err := fizzbuzz(number)
			if err != nil {
				return false
			}
			parsed, err := strconv.ParseInt(result, 10, 64)
			return err == nil && parsed == int64(number)
		},
		gen.IntRange(1, math.MaxInt32).SuchThat(func(v interface{}) bool {
			return v.(int)%3 != 0 && v.(int)%5 != 0
		}),
	))

	testing.Main(
		func(a, b string) (bool, error) { return true, nil },
		[]testing.InternalTest{
			{
				Name: "Example_fizzbuzz",
				F:    properties.RunT,
			},
		}, nil, nil)
	// Output:
	//
	// PASS
}
