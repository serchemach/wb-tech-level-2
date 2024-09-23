package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func lessStrings(a, b string) (bool, error, error) {
	return a < b, nil, nil
}

func lessNums(a, b string) (bool, error, error) {
	aFloat, err1 := strconv.ParseFloat(a, 64)
	bFloat, err2 := strconv.ParseFloat(b, 64)

	return aFloat < bFloat, err1, err2
}

var monthOrdering = map[string]int{
	"APR": 4,
	"AUG": 8,
	"DEC": 12,
	"FEB": 2,
	"JAN": 1,
	"JUL": 7,
	"JUN": 6,
	"MAR": 3,
	"MAY": 5,
	"NOV": 11,
	"OCT": 10,
	"SEP": 9,
}

func lessMonths(a, b string) (bool, error, error) {
	aNum, ok1 := monthOrdering[a]
	bNum, ok2 := monthOrdering[b]
	var err1, err2 error
	if !ok1 {
		err1 = CompError{}
	}
	if !ok2 {
		err2 = CompError{}
	}

	return aNum < bNum, err1, err2
}

var numSuffixes = map[string]float64{
	"K": math.Pow10(3),
	"k": math.Pow10(3),
	"M": math.Pow10(6),
	"G": math.Pow10(9),
	"T": math.Pow10(12),
	"P": math.Pow10(15),
	"E": math.Pow10(18),
	"Z": math.Pow10(21),
	"Y": math.Pow10(24),
}

type CompError struct {
}

func (err CompError) Error() string {
	return "Incomparable values"
}

func strToFloatSuffix(a string) (float64, error) {
	if len(a) < 1 {
		return 0, CompError{}
	}

	degree := 1.0

	if val, ok := numSuffixes[a[len(a)-1:]]; ok {
		degree = val
	}

	preffix, err := strconv.ParseFloat(a[:len(a)-1], 64)
	if err != nil {
		return 0, err
	}

	return preffix * degree, nil
}

func lessNumsSuffix(a, b string) (bool, error, error) {
	aFloat, err1 := strToFloatSuffix(a)
	bFloat, err2 := strToFloatSuffix(b)

	return aFloat < bFloat, err1, err2
}

func lessReversed[T any](f func(T, T) bool) func(T, T) bool {
	return func(a, b T) bool {
		return !f(a, b)
	}
}

func compareWithErr(b bool, err1, err2 error, s1, s2 string) bool {
	if err1 != nil && err2 == nil {
		return true
	}

	if err1 == nil && err2 != nil {
		return false
	}

	if err1 != nil && err2 != nil {
		res, _, _ := lessStrings(s1, s2)
		return res
	}

	return b
}

func LineLess(i, j int, lines []Line, column int, lessFunc func(string, string) (bool, error, error)) bool {
	lineI := lines[i].sortingContent
	lineJ := lines[j].sortingContent
	// Empty lines are equal, so they are never less
	if len(lineI) == len(lineJ) && len(lineI) == 0 {
		return false
	}

	if len(lineI) < column && len(lineJ) > len(lineI) {
		return true
	}

	if len(lineJ) < column && len(lineI) > len(lineJ) {
		return false
	}

	if len(lineI) < column && len(lineJ) == len(lineJ) {
		b, err1, err2 := lessFunc(lineI[0], lineJ[0])
		return compareWithErr(b, err1, err2, lineI[0], lineJ[0])
	}

	if len(lineJ) < column && len(lineI) == len(lineJ) {
		b, err1, err2 := lessFunc(lineJ[0], lineI[0])
		return compareWithErr(b, err1, err2, lineJ[0], lineI[0])
	}

	b, err1, err2 := lessFunc(lineI[column-1], lineJ[column-1])
	return compareWithErr(b, err1, err2, lineI[column-1], lineJ[column-1])
}

func SortLines(lines []Line, column int, lessFunc func(string, string) (bool, error, error), sortReverse bool) error {
	globalLess := func(i, j int) bool {
		return LineLess(i, j, lines, column, lessFunc)
	}
	if sortReverse {
		globalLess = lessReversed(globalLess)
	}
	sort.Slice(lines, globalLess)
	return nil
}

func CheckLinesForSorted(lines []Line, column int, lessFunc func(string, string) (bool, error, error), sortReverse bool) bool {
	globalLess := func(i, j int) bool {
		return LineLess(i, j, lines, column, lessFunc)
	}
	if sortReverse {
		globalLess = lessReversed(globalLess)
	}
	return sort.SliceIsSorted(lines, globalLess)
}

type SortFlags struct {
	sortByNums           bool
	sortReverse          bool
	removeRepeating      bool
	sortByMonth          bool
	ignoreTrailingSpaces bool
	sortByNumsWithSuffix bool
	checkSorted          bool
}

const DefaultSeparator = " "

type Line struct {
	sortingContent []string
	initialIndex   int
}

func SortFile(file string, column int, flags SortFlags, separator string) string {
	lines := strings.Split(file, "\n")
	splitLines := make([]Line, len(lines))
	for i, line := range lines {
		if flags.ignoreTrailingSpaces {
			line = strings.TrimSpace(line)
		}

		splitLines[i].sortingContent = strings.Split(line, separator)
		splitLines[i].initialIndex = i
	}

	lessFunc := lessStrings
	switch {
	case flags.sortByNums:
		lessFunc = lessNums
	case flags.sortByMonth:
		lessFunc = lessMonths
	case flags.sortByNumsWithSuffix:
		lessFunc = lessNumsSuffix
	}
	if flags.checkSorted {
		if CheckLinesForSorted(splitLines, column, lessFunc, flags.sortReverse) {
			return "Файл отсортирован"
		}
		return "Файл не отсортирован"
	}
	SortLines(splitLines, column, lessFunc, flags.sortReverse)

	var b strings.Builder
	for i, line := range splitLines {
		if flags.removeRepeating &&
			i != 0 && lines[splitLines[i-1].initialIndex] == lines[line.initialIndex] {
			continue
		}

		b.WriteString(lines[line.initialIndex])
		if i != len(splitLines)-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func main() {
	l := log.New(os.Stderr, "", 1)
	var (
		sortParams SortFlags
		column     int
	)
	flag.IntVar(&column, "k", 1, "Колонка для сортировки, разделитель колонок - пробел")
	flag.BoolVar(&sortParams.sortByNums, "n", false, "Сортировать по числовому значению")
	flag.BoolVar(&sortParams.sortReverse, "r", false, "Сортировать в обратном порядке")
	flag.BoolVar(&sortParams.removeRepeating, "u", false, "Не выводить повторяющиеся строки")
	flag.BoolVar(&sortParams.sortByMonth, "M", false, "Сортировать по названию месяца (в формате SEP, JAN, ...)")
	flag.BoolVar(&sortParams.ignoreTrailingSpaces, "b", false, "Игнорировать хвостовые пробелы")
	flag.BoolVar(&sortParams.checkSorted, "c", false, "Проверить, отсортированы ли данные. Сортировка в данном случае проводиться не будет")
	flag.BoolVar(&sortParams.sortByNumsWithSuffix, "h", false, "Сортировать по числовому значению с учётом суффиксов (2k, 2K, 2B, ...)")
	flag.Parse()

	if len(flag.Args()) == 0 {
		l.Fatal("На вход не было подано файлов")
	} else if len(flag.Args()) > 1 {
		l.Fatal("На вход было подано больше одного файла")
	}

	fileContents, err := os.ReadFile(flag.Args()[0])
	if err != nil {
		l.Fatalf("Ошибка во время чтения файла: %s\n", err)
	}

	if fileContents[len(fileContents)-1] == '\n' {
		fileContents = fileContents[:len(fileContents)-1]
	}

	result := SortFile(string(fileContents), column, sortParams, DefaultSeparator)
	fmt.Println(result)
}
