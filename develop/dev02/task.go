package main

import (
	// "fmt"
	// "bytes"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

const numToken = 0
const charToken = 1
const escapeToken = 2

type EmptyTokenError struct{}

func (err EmptyTokenError) Error() string {
	return "No characters to parse"
}

type BadTokenError struct{}

func (err BadTokenError) Error() string {
	return "Couldn't decode a character of the token in utf8"
}

type BadNumTokenError struct{}

func (err BadNumTokenError) Error() string {
	return "Can't repeat an empty string"
}

func parseNum(s string) int {
	for i, char := range s {
		if !unicode.IsDigit(char) {
			return i
		}
	}
	return len(s)
}

func parseToken(s string) (int, string, error) {
	if len(s) == 0 {
		return 0, "", EmptyTokenError{}
	}

	firstChar, size := utf8.DecodeRuneInString(s)
	if firstChar == utf8.RuneError {
		return 0, "", BadTokenError{}
	}

	if !unicode.IsPrint(firstChar) {
		return 0, "", BadTokenError{}
	}

	switch {
	case firstChar == '\\':
		secondChar, secondSize := utf8.DecodeRuneInString(s[size:])
		if secondChar == utf8.RuneError {
			return 0, "", BadTokenError{}
		}
		return escapeToken, s[size : secondSize+size], nil
	case unicode.IsDigit(firstChar):
		numEnd := parseNum(s)
		return numToken, s[:numEnd], nil
	default:
		return charToken, s[:size], nil
	}
}

func Unpack(s string) (string, error) {
	tokenEnd := 0
	prevTokenStr := ""
	var b strings.Builder

	for tokenEnd <= len(s)-1 {
		tokenType, tokenStr, err := parseToken(s[tokenEnd:])
		tokenEnd += len(tokenStr)
		if tokenType == escapeToken {
			tokenEnd += len("\\")
		}
		// fmt.Println(tokenType, tokenStr, tokenEnd, len(s)-1, tokenEnd < len(s)-1)

		if err != nil {
			return "", err
		}

		if tokenType == numToken {
			numRepeat, err := strconv.Atoi(tokenStr)
			if err != nil {
				return "", err
			}

			if prevTokenStr == "" {
				return "", BadNumTokenError{}
			}

			for _ = range numRepeat - 1 {
				b.WriteString(prevTokenStr)
			}
		} else {
			b.WriteString(tokenStr)
			prevTokenStr = tokenStr
		}
	}

	return b.String(), nil
}

func main() {

}
