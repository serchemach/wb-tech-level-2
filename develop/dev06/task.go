package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

type WrongRangeFormatError struct{}

func (e WrongRangeFormatError) Error() string {
	return "Неправильный формат промежутка чисел"
}

type ColRange struct {
	start int
	end   int
}

func (r ColRange) IsInRange(i int) bool {
	return i >= r.start && i <= r.end || i >= r.start && r.end == -1
}

func NewColRange(s string) (ColRange, error) {
	numDashes := strings.Count(s, "-")
	if numDashes > 1 {
		return ColRange{}, WrongRangeFormatError{}
	}

	if numDashes == 0 {
		num, err := strconv.Atoi(s)
		if err != nil {
			return ColRange{}, err
		}

		return ColRange{num, num}, nil
	}

	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return ColRange{}, WrongRangeFormatError{}
	}

	if parts[0] == "" && parts[1] == "" {
		return ColRange{
			start: 0,
			end:   -1,
		}, nil
	}

	var err error
	result := ColRange{}
	if parts[0] == "" {
		result.start = 0
	} else {
		result.start, err = strconv.Atoi(parts[0])
		if err != nil {
			return result, err
		}
	}

	if parts[1] == "" {
		result.end = -1
	} else {
		result.end, err = strconv.Atoi(parts[1])
		if err != nil {
			return result, err
		}
	}

	return result, nil
}

func ParseFieldNums(nums string) ([]ColRange, error) {
	groups := strings.Split(nums, ",")
	result := make([]ColRange, len(groups))
	var err error

	for i, group := range groups {
		result[i], err = NewColRange(group)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func KeepRightColumns(cols []string, colRanges []ColRange, delim string, output io.StringWriter) {
	written := false
	for i, col := range cols {
		for _, colRange := range colRanges {
			if colRange.IsInRange(i + 1) {
				if written {
					output.WriteString(delim)
				}
				output.WriteString(col)
				written = true
			}
		}
	}
	output.WriteString("\n")
}

func SplitToColumns(input *bufio.Reader, colRanges []ColRange, delim string, keepOnlyDelim bool, output io.StringWriter) error {
	line, err := input.ReadString('\n')
	if err == nil && line[len(line)-1] == '\n' {
		line = line[:len(line)-1]
	}

	for err == nil {
		columns := strings.Split(line, delim)

		// Строки без разделители можно оставить если не стоит флаг
		if len(columns) == 1 && columns[0] == line && !keepOnlyDelim {
			output.WriteString(line)
			output.WriteString("\n")
		} else if len(columns) > 1 {
			KeepRightColumns(columns, colRanges, delim, output)
		}

		line, err = input.ReadString('\n')
		if err == nil && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}
	}

	if !errors.Is(err, io.EOF) {
		return err
	}

	return nil
}

func main() {
	l := log.New(os.Stderr, "", 1)
	var (
		fieldNums     string
		delim         string
		onlyWithDelim bool
	)
	flag.StringVar(&fieldNums, "f", "1", "Номера колонок (начиная с 1), которые нужно оставить.\nФормат: x-y, где y >= x >= 1. Если необходимы только некоторые колонки, то можно перечислить их номера через ','.")
	flag.StringVar(&delim, "d", "\t", "Разделитель колонок.")
	flag.BoolVar(&onlyWithDelim, "s", false, "Оставить только строки с разделителем.")
	flag.Parse()

	colRanges, err := ParseFieldNums(fieldNums)
	if err != nil {
		l.Fatal("Ошибка при обработке номеров колонок: ", err)
	}

	// Дисклеймер: для того, чтобы не вводить данные в программу не только через "<" из файла, но и просто в консоли, необходимо в конце ввода нажать Ctrl-D для того, чтобы написать символ конца файла
	reader := bufio.NewReader(os.Stdin)
	SplitToColumns(reader, colRanges, delim, onlyWithDelim, os.Stdout)
}
