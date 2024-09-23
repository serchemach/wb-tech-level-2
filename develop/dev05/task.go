package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type PseudoQueue struct {
	storage    []string
	beginIndex int
	endIndex   int
	isFull     bool
}

func NewQueue(size int) PseudoQueue {
	return PseudoQueue{
		storage:    make([]string, size),
		beginIndex: 0,
		endIndex:   0,
		isFull:     false,
	}
}

func (q *PseudoQueue) Append(s string) {
	if len(q.storage) == 0 {
		return
	}

	if q.isFull {
		q.storage[q.beginIndex] = s
		q.beginIndex = (q.beginIndex + 1) % len(q.storage)
		q.endIndex = q.beginIndex
		return
	}

	q.storage[q.endIndex] = s
	q.endIndex = (q.endIndex + 1) % len(q.storage)
	if q.beginIndex == q.endIndex {
		q.isFull = true
	}
}

func (q PseudoQueue) Print(output io.Writer, printNumbers bool, firstLine int) {
	if !q.isFull && q.beginIndex == q.endIndex || len(q.storage) == 0 {
		return
	}

	curLine := firstLine
	if printNumbers {
		fmt.Fprintf(output, "%d:", curLine)
	}
	fmt.Fprintln(output, q.storage[0])
	curLine++

	for i := (q.beginIndex + 1) % len(q.storage); i != q.endIndex; i = (i + 1) % len(q.storage) {
		if printNumbers {
			fmt.Fprintf(output, "%d:", curLine)
		}
		fmt.Fprintln(output, q.storage[i])
		curLine++
	}
}

func FindMatching(reader *bufio.Reader, matcher func(string) bool, params SearchFlags, output io.Writer) error {
	queue := NewQueue(params.numBefore)
	leftToPrint := 0
	prevLeftToPrint := 0
	curLineNumber := 1

	line, err := reader.ReadString('\n')
	for err == nil {
		matches := matcher(line)
		prevLeftToPrint = leftToPrint
		if leftToPrint > 0 && !matches {
			if params.printLineNumber {
				fmt.Fprintf(output, "%d:", curLineNumber)
			}
			fmt.Fprintln(output, line)
			leftToPrint--
		}

		if matches {
			if prevLeftToPrint == 0 && leftToPrint == 0 {
				queue.Print(output, params.printLineNumber, curLineNumber-params.numBefore)
			}
			if params.printLineNumber {
				fmt.Fprintf(output, "%d-", curLineNumber)
			}
			fmt.Fprintln(output, line)
			leftToPrint = params.numAfter
		}

		queue.Append(line)

		curLineNumber++
		line, err = reader.ReadString('\n')
	}

	if !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func GetCount(reader *bufio.Reader, matcher func(string) bool) (int, error) {
	count := 0
	line, err := reader.ReadString('\n')
	for err == nil {
		if matcher(line) {
			count++
		}
		line, err = reader.ReadString('\n')
	}

	if !errors.Is(err, io.EOF) {
		return 0, err
	}

	return count, nil
}

func AssembleMatcher(pattern string, params SearchFlags) (func(string) bool, error) {
	if params.searchFixed {
		pattern = regexp.QuoteMeta(pattern)
	}

	if params.ignoreCase {
		pattern = "(?i)" + pattern
	}

	expr, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	// fmt.Println(pattern)

	if params.exclude {
		return func(line string) bool {
			return !expr.MatchString(line)
		}, nil
	}

	return func(line string) bool {
		return expr.MatchString(line)
	}, nil
}

type SearchFlags struct {
	numAfter        int
	numBefore       int
	printCount      bool
	ignoreCase      bool
	exclude         bool
	searchFixed     bool
	printLineNumber bool
}

func main() {
	l := log.New(os.Stderr, "", 1)

	var params SearchFlags
	contextNum := 0
	flag.IntVar(&params.numAfter, "A", 0, "Напечатать +N строк после совпадения")
	flag.IntVar(&params.numBefore, "B", 0, "Напечатать +N строк до совпадения")
	flag.IntVar(&contextNum, "C", 0, "Напечатать +-N строк вокруг совпадения")
	flag.BoolVar(&params.printCount, "c", false, "Вместо вывода строк, вывести их количество")
	flag.BoolVar(&params.ignoreCase, "i", false, "Игнорировать регистр")
	flag.BoolVar(&params.exclude, "v", false, "Вместо совпадения, исключать строки с паттерном")
	flag.BoolVar(&params.searchFixed, "F", false, "Точное совпадение со строкой, не паттерн")
	flag.BoolVar(&params.printLineNumber, "n", false, "Печатать номера строк")
	flag.Parse()

	if params.numAfter == 0 && params.numBefore == 0 {
		params.numAfter, params.numBefore = contextNum, contextNum
	}

	if len(flag.Args()) < 1 {
		l.Fatal("Не предоставлен паттерн")
	}

	matcher, err := AssembleMatcher(flag.Args()[0], params)
	if err != nil {
		l.Fatal(err)
	}

	// fmt.Println(params)
	for _, filename := range flag.Args()[1:] {
		file, err := os.Open(filename)
		// 	// fmt.Println(filename)
		if err != nil {
			l.Fatal(err)
		}

		reader := bufio.NewReader(file)
		if params.printCount {
			result, err := GetCount(reader, matcher)
			if err != nil {
				l.Fatal(err)
			}
			fmt.Println(result)
		} else {
			FindMatching(reader, matcher, params, os.Stdout)
		}

	}

}
