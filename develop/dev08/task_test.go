package main

import (
	"slices"
	"testing"
)

type TestCase struct {
	input  string
	output Statement
	hint   string
}

func StatementEqual(a, b Statement) bool {
	if len(a.args) != len(b.args) {
		return false
	}

	i := 0
	for ; i < len(a.args) && slices.Equal(a.args[i], b.args[i]); i++ {
	}
	if i != len(a.args) {
		return false
	}

	return a.detach == b.detach && slices.Equal(a.pipeChain, b.pipeChain)
}

func TestParseStatement(t *testing.T) {
	testCases := []TestCase{
		{
			input: "echo 123",
			output: Statement{
				pipeChain: []Command{
					EchoCommand{},
				},
				args: [][]string{
					{"123"},
				},
			},
			hint: "Simple builtin command test",
		},
		{
			input: "echo 123 | echo",
			output: Statement{
				pipeChain: []Command{
					EchoCommand{},
					EchoCommand{},
				},
				args: [][]string{
					{"123"},
					{},
				},
			},
			hint: "Simple builtin pipe command test",
		},
		{
			input: "echo 123 | cat",
			output: Statement{
				pipeChain: []Command{
					EchoCommand{},
					ExecutableCommand{name: "cat"},
				},
				args: [][]string{
					{"123"},
					{},
				},
			},
			hint: "Simple external command pipe test",
		},
		{
			input: "echo 123 | cat &",
			output: Statement{
				pipeChain: []Command{
					EchoCommand{},
					ExecutableCommand{name: "cat"},
				},
				args: [][]string{
					{"123"},
					{},
				},
				detach: true,
			},
			hint: "Detach with pipe test",
		},
		{
			input: "cat &",
			output: Statement{
				pipeChain: []Command{
					ExecutableCommand{name: "cat"},
				},
				args: [][]string{
					{},
				},
				detach: true,
			},
			hint: "Detach without pipe test",
		},
	}

	for _, test := range testCases {
		t.Run(test.hint, func(t *testing.T) {
			output, err := ParseStatement(test.input)
			if err != nil {
				t.Fatal("Unexpected error", err)
			}

			if !StatementEqual(output, test.output) {
				t.Fatalf("Wrong output\nExpected: %v\nRecieved: %v\n", test.output, output)
			}
		})
	}
}
