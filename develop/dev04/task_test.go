package main

import (
	"reflect"
	"testing"
)

type testCase[T any, Q any] struct {
	input  T
	output Q
	hint   string
}

func TestAnagramGrouping(t *testing.T) {
	testCases := []testCase[[]string, map[string][]string]{
		{
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик"},
			output: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
			hint: "Simple example",
		},
		{
			input: []string{"пятак", "ПЯтка", "тЯпка", "лиСТок", "слиТок", "столик"},
			output: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
			hint: "Example with uppercase letters",
		},
		{
			input: []string{"ПЯтка", "лиСТок", "слиТок", "столик"},
			output: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
			},
			hint: "Example with a group of len 1",
		},
		{
			input: []string{"пятак", "пятка", "тяпка", "листок", "слиток", "Столик", "замок", "комаз", "Овощи", "кИлотс"},
			output: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"килотс", "листок", "слиток", "столик"},
				"замок":  {"замок", "комаз"},
			},
			hint: "Big example",
		},
	}

	for _, test := range testCases {
		t.Run(test.hint, func(t *testing.T) {
			testedOutput := GroupByAnagrams(test.input)

			if !reflect.DeepEqual(testedOutput, test.output) {
				t.Fatalf("Wrong output:\nexpected %v,\nrecieved %v\n", test.output, testedOutput)
			}
		})
	}

}
