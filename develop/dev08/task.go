package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
=== Взаимодействие с ОС ===

Необходимо реализовать собственный шелл

встроенные команды: cd/pwd/echo/kill/ps
поддержать fork/exec команды
конвеер на пайпах

Реализовать утилиту netcat (nc) клиент
принимать данные из stdin и отправлять в соединение (tcp/udp)
Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type Command interface {
	Exec([]string, io.Reader, io.Writer) error
}

type EmptyCommand struct{}

func (c EmptyCommand) Exec(args []string, stdin io.Reader, stdout io.Writer) error {
	return nil
}

type NotEnoughArgumentsError struct{}

func (e NotEnoughArgumentsError) Error() string {
	return "Not enough arguments"
}

type TooManyArgumentsError struct{}

func (e TooManyArgumentsError) Error() string {
	return "Too many arguments"
}

type CdCommand struct{}

func (c CdCommand) Exec(args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 0 {
		return NotEnoughArgumentsError{}
	}

	if len(args) > 1 {
		return TooManyArgumentsError{}
	}

	return os.Chdir(args[0])
}

type PwdCommand struct{}

func (c PwdCommand) Exec(args []string, stdin io.Reader, stdout io.Writer) error {
	dir, err := os.Getwd()
	if err == nil {
		stdout.Write([]byte(dir))
		stdout.Write([]byte("\n"))
	}

	return err
}

type EchoCommand struct{}

func (c EchoCommand) Exec(args []string, stdin io.Reader, stdout io.Writer) error {
	for i, arg := range args {
		if i != 0 {
			stdout.Write([]byte(" "))
		}
		stdout.Write([]byte(arg))
	}
	stdout.Write([]byte("\n"))
	return nil
}

type KillCommand struct{}

func (c KillCommand) Exec(args []string, stdin io.Reader, stdout io.Writer) error {
	if len(args) == 0 {
		return NotEnoughArgumentsError{}
	}

	if len(args) > 1 {
		return TooManyArgumentsError{}
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return process.Kill()
}

type PsCommand struct{}

func (c PsCommand) Exec(args []string, stdin io.Reader, stdout io.Writer) error {
	return nil
}

type ExecutableCommand struct {
	name string
}

func (c ExecutableCommand) Exec(args []string, stdin io.Reader, stdout io.Writer) error {
	return nil
}

type Statement struct {
	pipeChain []Command
	args      [][]string
	detach    bool
}

func (s Statement) Run() error {
	var previousStdout io.Reader
	for i, command := range s.pipeChain {
		r, w := io.Pipe()
		var err error
		if i == len(s.pipeChain)-1 {
			err = command.Exec(s.args[i], previousStdout, os.Stdout)
		} else {
			err = command.Exec(s.args[i], previousStdout, w)
		}

		if err != nil {
			return err
		}
		previousStdout = r
	}

	return nil
}

type WrongForkUseError struct{}

func (err WrongForkUseError) Error() string {
	return "The fork should only be used in a last command when using pipes"
}

func ParseCommand(name string) Command {
	switch name {
	case "":
		return EmptyCommand{}

	case "cd":
		return CdCommand{}
	case "pwd":
		return PwdCommand{}
	case "echo":
		return EchoCommand{}
	case "kill":
		return KillCommand{}
	case "ps":
		return PsCommand{}
	}

	return ExecutableCommand{name}
}

func ParseStatement(s string) (Statement, error) {
	commands := strings.Split(s, "|")
	statement := Statement{
		pipeChain: make([]Command, len(commands)),
		args:      make([][]string, len(commands)),
	}
	for i, command := range commands {
		fmt.Println(command)
		tokens := strings.Split(strings.Trim(command, " "), " ")
		if i != len(commands)-1 && tokens[len(tokens)-1] == "&" {
			return Statement{}, WrongForkUseError{}
		}

		statement.pipeChain[i] = ParseCommand(tokens[0])
		statement.args[i] = tokens[1:]
	}

	return statement, nil
}

func main() {
	l := log.New(os.Stderr, "", 1)

	for {
		// var line string
		curDir, err := os.Getwd()
		if err != nil {
			l.Fatal(err)
		}
		fmt.Printf("%s>", curDir)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		line := scanner.Text()

		if scanner.Err() != nil {
			fmt.Println(scanner.Err())
			continue
		}
		// fmt.Scanln(&line)

		statement, err := ParseStatement(line)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(statement)

		err = statement.Run()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

}
