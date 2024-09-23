package main

import (
	"bufio"
	"os"
	"slices"
	"strings"
	"testing"
)

type testCase[T any, Q any] struct {
	input  T
	output Q
	hint   string
}

func TestRangeParsing(t *testing.T) {
	testCases := []testCase[string, []ColRange]{
		{
			input:  "1",
			output: []ColRange{{1, 1}},
			hint:   "Simple test",
		},
		{
			input:  "1-",
			output: []ColRange{{1, -1}},
			hint:   "Simple infinite end range test",
		},
		{
			input:  "-2",
			output: []ColRange{{0, 2}},
			hint:   "Simple infinite beginning range test",
		},
		{
			input:  "1-2,3-,6,-10",
			output: []ColRange{{1, 2}, {3, -1}, {6, 6}, {0, 10}},
			hint:   "Multiple ranges test",
		},
	}

	for _, test := range testCases {
		t.Run(test.hint, func(t *testing.T) {
			predictedOutput, err := ParseFieldNums(test.input)

			if err != nil {
				t.Fatal("Encountered an error: ", err)
			}

			if !slices.Equal(test.output, predictedOutput) {
				t.Fatalf("Test output is wrong!\n Expected: %v\n Recieved:%v\n", test.output, predictedOutput)
			}
		})
	}
}

type ColumnFilteringInput struct {
	filename      string
	colRanges     []ColRange
	delim         string
	keepOnlyDelim bool
}

func TestColumnFiltering(t *testing.T) {
	testCases := []testCase[ColumnFilteringInput, string]{
		{
			input: ColumnFilteringInput{
				filename:      "test_files/ultimateTest.txt",
				colRanges:     []ColRange{{0, -1}},
				delim:         " ",
				keepOnlyDelim: false,
			},
			output: "test_files/output1.txt",
			hint:   "No filtering test",
		},
		{
			input: ColumnFilteringInput{
				filename:      "test_files/ultimateTest.txt",
				colRanges:     []ColRange{{1, 1}},
				delim:         " ",
				keepOnlyDelim: false,
			},
			output: "test_files/output2.txt",
			hint:   "Filter to one column",
		},
		{
			input: ColumnFilteringInput{
				filename:      "test_files/ultimateTest.txt",
				colRanges:     []ColRange{{2, 3}},
				delim:         " ",
				keepOnlyDelim: true,
			},
			output: "test_files/output3.txt",
			hint:   "Filter and discard non delim lines",
		},
		{
			input: ColumnFilteringInput{
				filename:      "test_files/ultimateTest.txt",
				colRanges:     []ColRange{{2, 2}, {4, 4}},
				delim:         " ",
				keepOnlyDelim: false,
			},
			output: "test_files/output4.txt",
			hint:   "Filter with multiple col ranges",
		},
	}

	for _, test := range testCases {
		t.Run(test.hint, func(t *testing.T) {
			inputFile, err := os.Open(test.input.filename)
			if err != nil {
				t.Fatal("Encountered an error: ", err)
			}

			reader := bufio.NewReader(inputFile)

			output, err := os.ReadFile(test.output)
			if err != nil {
				t.Fatal("Encountered an error: ", err)
			}

			var outputWriter strings.Builder
			err = SplitToColumns(reader, test.input.colRanges, test.input.delim, test.input.keepOnlyDelim, &outputWriter)
			if err != nil {
				t.Fatal("Encountered an error: ", err)
			}

			predictedOutput := outputWriter.String()
			if string(output) != predictedOutput {
				t.Fatalf("Test output is wrong!\n Expected: %v\n Recieved:%v\nExpected bytes:%v\nRecieved bytes:%v\n", string(output), predictedOutput, output, []byte(predictedOutput))
			}
		})
	}
}
