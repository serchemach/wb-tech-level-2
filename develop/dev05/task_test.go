package main

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

type GrepInput struct {
	params   SearchFlags
	pattern  string
	filename string
}

type testCase[T any, Q any] struct {
	input  T
	output Q
	hint   string
}

func TestGrepLines(t *testing.T) {
	testCases := []testCase[GrepInput, string]{
		{
			input: GrepInput{
				params: SearchFlags{
					numAfter:        0,
					numBefore:       0,
					printCount:      false,
					ignoreCase:      false,
					exclude:         false,
					searchFixed:     true,
					printLineNumber: false,
				},
				pattern:  "1",
				filename: "test_files/input1.txt",
			},
			output: "test_files/output1.txt",
			hint:   "Simple example",
		},
		{
			input: GrepInput{
				params: SearchFlags{
					numAfter:        1,
					numBefore:       1,
					printCount:      false,
					ignoreCase:      false,
					exclude:         false,
					searchFixed:     true,
					printLineNumber: true,
				},
				pattern:  "1",
				filename: "test_files/input1.txt",
			},
			output: "test_files/output2.txt",
			hint:   "Example with context and line numbers",
		},
		{
			input: GrepInput{
				params: SearchFlags{
					numAfter:        1,
					numBefore:       1,
					printCount:      false,
					ignoreCase:      false,
					exclude:         true,
					searchFixed:     true,
					printLineNumber: true,
				},
				pattern:  "1",
				filename: "test_files/input1.txt",
			},
			output: "test_files/output3.txt",
			hint:   "Example with invert and with context and line numbers",
		},
	}

	for _, test := range testCases {
		t.Run(test.hint, func(t *testing.T) {
			file, err := os.Open(test.input.filename)
			if err != nil {
				t.Fatal("Error reading test input, ", err)
			}

			output, err := os.ReadFile(test.output)
			if err != nil {
				t.Fatal("Error reading test output, ", err)
			}

			matcher, err := AssembleMatcher(test.input.pattern, test.input.params)
			if err != nil {
				t.Fatal("Error while assembling the matcher, ", err)
			}

			reader := bufio.NewReader(file)
			var b strings.Builder
			err = FindMatching(reader, matcher, test.input.params, &b)

			inferredOutput := b.String()
			if inferredOutput != string(output) {
				t.Fatalf("Wrong output:\nexpected %v,\nrecieved %v\n", string(output), inferredOutput)
			}
		})
	}

}
