package main

import (
	"bufio"
	"fmt"
	"time"

	// "syscall"

	// "io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
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

type RunningProcess struct {
	name  string
	start time.Time
}

var runningProcesses map[int]RunningProcess

type Command interface {
	Exec([]string, *os.File, *os.File, bool) error
}

type EmptyCommand struct{}

func (c EmptyCommand) Exec(args []string, stdin *os.File, stdout *os.File, detach bool) error {
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

func (c CdCommand) Exec(args []string, stdin *os.File, stdout *os.File, detach bool) error {
	if len(args) == 0 {
		return NotEnoughArgumentsError{}
	}

	if len(args) > 1 {
		return TooManyArgumentsError{}
	}

	return os.Chdir(args[0])
}

type PwdCommand struct{}

func (c PwdCommand) Exec(args []string, stdin *os.File, stdout *os.File, detach bool) error {
	dir, err := os.Getwd()
	if err == nil {
		stdout.Write([]byte(dir))
		stdout.Write([]byte("\n"))
	}

	return err
}

type EchoCommand struct{}

func (c EchoCommand) Exec(args []string, stdin *os.File, stdout *os.File, detach bool) error {
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

func (c KillCommand) Exec(args []string, stdin *os.File, stdout *os.File, detach bool) error {
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

func (c PsCommand) Exec(args []string, stdin *os.File, stdout *os.File, detach bool) error {
	fmt.Fprintln(stdout, "PID TIME NAME")
	for pid, process := range runningProcesses {
		fmt.Fprintf(stdout, "%d %s %s\n", pid, time.Since(process.start).String(), process.name)
	}
	return nil
}

type ExitCommand struct{}

func (c ExitCommand) Exec(args []string, stdin *os.File, stdout *os.File, detach bool) error {
	os.Exit(0)
	return nil
}

type ExecutableCommand struct {
	name string
}

func (c ExecutableCommand) Exec(args []string, stdin *os.File, stdout *os.File, detach bool) error {
	cmd := exec.Command(c.name, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout

	var err error
	if detach {
		go func() {
			err = cmd.Start()
			runningProcesses[cmd.Process.Pid] = RunningProcess{
				name:  c.name,
				start: time.Now(),
			}
			err = cmd.Wait()
			delete(runningProcesses, cmd.Process.Pid)
		}()
	} else {
		err = cmd.Run()
	}

	return err
}

type Statement struct {
	pipeChain []Command
	args      [][]string
	detach    bool
}

func (s Statement) Run() error {
	var previousStdout *os.File
	var wg sync.WaitGroup
	var err error
	for i, command := range s.pipeChain {
		wg.Add(1)
		r, w, err2 := os.Pipe()
		if err2 != nil {
			return err2
		}

		if i == len(s.pipeChain)-1 {
			err = command.Exec(s.args[i], previousStdout, os.Stdout, s.detach)
		} else {
			err = command.Exec(s.args[i], previousStdout, w, false)
			w.Close()
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
	case "exit":
		return ExitCommand{}
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
		tokens := strings.Split(strings.Trim(command, " "), " ")
		if i != len(commands)-1 && tokens[len(tokens)-1] == "&" {
			return Statement{}, WrongForkUseError{}
		}

		if tokens[len(tokens)-1] == "&" {
			statement.detach = true
			tokens = tokens[:len(tokens)-1]
		}

		statement.pipeChain[i] = ParseCommand(tokens[0])
		statement.args[i] = tokens[1:]
	}

	return statement, nil
}

func main() {
	l := log.New(os.Stderr, "", 1)
	runningProcesses = make(map[int]RunningProcess)

	for {
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

		err = statement.Run()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

}
