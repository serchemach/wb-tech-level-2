package main

import (
	"fmt"
	"slices"
	"strings"
	"unicode"
)

/*
=== Поиск анаграмм по словарю ===

Напишите функцию поиска всех множеств анаграмм по словарю.
Например:
'пятак', 'пятка' и 'тяпка' - принадлежат одному множеству,
'листок', 'слиток' и 'столик' - другому.

Входные данные для функции: ссылка на массив - каждый элемент которого - слово на русском языке в кодировке utf8.
Выходные данные: Ссылка на мапу множеств анаграмм.
Ключ - первое встретившееся в словаре слово из множества
Значение - ссылка на массив, каждый элемент которого, слово из множества. Массив должен быть отсортирован по возрастанию.
Множества из одного элемента не должны попасть в результат.
Все слова должны быть приведены к нижнему регистру.
В результате каждое слово должно встречаться только один раз.

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func GetCharCounts(word string) string {
	counts := map[rune]int{}
	for _, char := range word {
		if _, ok := counts[char]; ok {
			counts[char]++
		} else {
			counts[char] = 1
		}
	}
	return fmt.Sprintf("%v", counts)
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func GroupByAnagrams(words []string) map[string][]string {
	anagrams := map[string][]string{}

	for _, word := range words {
		if !IsLower(word) {
			word = strings.ToLower(word)
		}

		signature := GetCharCounts(word)
		if slice, ok := anagrams[signature]; ok {
			// Попытка не супер много раз выделять память
			if cap(slice) < len(slice)+1 {
				anagrams[signature] = append(slice, make([]string, len(slice))...)[:len(slice)+1]
			} else {
				anagrams[signature] = slice[:len(slice)+1]
			}
			anagrams[signature][len(slice)] = word
		} else {
			anagrams[signature] = []string{word}
		}
	}

	result := map[string][]string{}
	for _, arr := range anagrams {
		if len(arr) <= 1 {
			continue
		}

		newKey := arr[0]
		slices.Sort(arr)
		result[newKey] = arr
	}
	return result
}

func main() {
}
