package main

import (
	"testing"
)

type testCase struct {
	input  string
	output string
	hint   string
}

func TestUnpack(t *testing.T) {
	testCases := []testCase{
		{"abcd", "abcd", "Plain string"},
		{"a4bc2d5e", "aaaabccddddde", "Simple unpacking"},
		{"", "", "Empty string"},
		{`qwe\4\5`, "qwe45", "Simple escape sequence"},
		{`qwe\45`, "qwe44444", "Repeat an escape sequence"},
		{`qwe\\5`, `qwe\\\\\`, "Repeat an escaped backslash"},
	}

	for _, test := range testCases {
		t.Run(test.hint, func(t *testing.T) {
			testOutput, err := Unpack(test.input)
			if err != nil {
				t.Fatal("Error while unpacking: ", err)
			}

			if testOutput != test.output {
				t.Fatalf("Wrong output:\nexpected %s,\nrecieved %s", test.output, testOutput)
			}
		})
	}

	errCases := []testCase{
		{"45", "", "Bad number token"},
		{`45\`, "", "Escape sequence on nothing"},
	}

	for _, test := range errCases {
		t.Run(test.hint, func(t *testing.T) {
			testOutput, err := Unpack(test.input)
			if err == nil {
				t.Fatal("Expected an error")
			}

			if testOutput != test.output {
				t.Fatalf("Wrong output:\nexpected %s,\nrecieved %s", test.output, testOutput)
			}
		})
	}
}

func BenchmarkUnpack(b *testing.B) {
	input := `qwe\45a8sdjfniuaehi3hrna8h`
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		Unpack(input)
	}
}
