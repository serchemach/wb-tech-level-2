package main

import "testing"

type sortFileInput struct {
	file      string
	column    int
	sortFlags SortFlags
	separator string
}

type checkFileInput struct {
	file      string
	column    int
	sortFlags SortFlags
	separator string
}

type testCase[T any, Q any] struct {
	input  T
	output Q
	hint   string
}

// Дисклеймер: при некоторых локалях на линухе sort может быть case insensitive даже без выбора данной опции
// В данных тестах для корректных ответов использовался case sensitive sort
func TestSortFile(t *testing.T) {
	testCases := []testCase[sortFileInput, string]{
		{sortFileInput{
			file: `12399 323 
23 9923
90909

98 100
ss sadsf`,
			column:    1,
			sortFlags: SortFlags{},
			separator: " ",
		}, `
12399 323 
23 9923
90909
98 100
ss sadsf`, "Basic sorting with enough columns"},

		{sortFileInput{
			file: `12399 323 
23 9923
90909

98 100
ss sadsf`,
			column:    2,
			sortFlags: SortFlags{},
			separator: " ",
		}, `
90909
98 100
12399 323 
23 9923
ss sadsf`, "Basic sorting with less columns than in some lines"},

		{sortFileInput{
			file: `12399 323 
23 9923
90909

98 100
ss sadsf`,
			column:    1,
			sortFlags: SortFlags{sortReverse: true},
			separator: " ",
		}, `ss sadsf
98 100
90909
23 9923
12399 323 
`, "Basic reversed sorting"},

		{sortFileInput{
			file: `12399 323 
23 9923
90909
aa

98 100
ss sadsf`,
			column:    1,
			sortFlags: SortFlags{sortByNums: true},
			separator: " ",
		}, `
aa
ss sadsf
23 9923
98 100
12399 323 
90909`, "Basic number sorting"},

		{sortFileInput{
			file: `123 12 3 12 3123 

djfi aioj adfsja 23j4 jidsfa
123 12j3 i33 132ji 
22 3k21j3k21j 3lkj12io3 j
1K 234 21j3 k12j
2G 23 12 3 dfska j
2G 23 12 3 dfska j
2G 23 12 3 dfska j
4M cv ajkjad jfda
4G idj foiaj dfj
2K ijados fja f 
JAN 2 3klj 21j djsfoi
DEC j2i3 k12j 3ij
JAN 2j 3k1j2oi j
SEP j jdsij afj
`,
			column:    1,
			sortFlags: SortFlags{sortByNumsWithSuffix: true},
			separator: " ",
		}, `

DEC j2i3 k12j 3ij
JAN 2 3klj 21j djsfoi
JAN 2j 3k1j2oi j
SEP j jdsij afj
djfi aioj adfsja 23j4 jidsfa
22 3k21j3k21j 3lkj12io3 j
123 12 3 12 3123 
123 12j3 i33 132ji 
1K 234 21j3 k12j
2K ijados fja f 
4M cv ajkjad jfda
2G 23 12 3 dfska j
2G 23 12 3 dfska j
2G 23 12 3 dfska j
4G idj foiaj dfj`, "Basic suffix number sorting"},

		{sortFileInput{
			file: `123 12 3 12 3123 

djfi aioj adfsja 23j4 jidsfa
123 12j3 i33 132ji 
22 3k21j3k21j 3lkj12io3 j
1K 234 21j3 k12j
2G 23 12 3 dfska j
2G 23 12 3 dfska j
2G 23 12 3 dfska j
4M cv ajkjad jfda
4G idj foiaj dfj
2K ijados fja f 
JAN 2 3klj 21j djsfoi
DEC j2i3 k12j 3ij
JAN 2j 3k1j2oi j
SEP j jdsij afj
`,
			column:    1,
			sortFlags: SortFlags{sortByMonth: true},
			separator: " ",
		}, `

123 12 3 12 3123 
123 12j3 i33 132ji 
1K 234 21j3 k12j
22 3k21j3k21j 3lkj12io3 j
2G 23 12 3 dfska j
2G 23 12 3 dfska j
2G 23 12 3 dfska j
2K ijados fja f 
4G idj foiaj dfj
4M cv ajkjad jfda
djfi aioj adfsja 23j4 jidsfa
JAN 2 3klj 21j djsfoi
JAN 2j 3k1j2oi j
SEP j jdsij afj
DEC j2i3 k12j 3ij`, "Basic month sorting"},

		{sortFileInput{
			file: `123 12 3 12 3123 

djfi aioj adfsja 23j4 jidsfa
123 12j3 i33 132ji 
22 3k21j3k21j 3lkj12io3 j
1K 234 21j3 k12j
2G 23 12 3 dfska j
2G 23 12 3 dfska j
2G 23 12 3 dfska j
4M cv ajkjad jfda
4G idj foiaj dfj
2K ijados fja f 
JAN 2 3klj 21j djsfoi
DEC j2i3 k12j 3ij
JAN 2j 3k1j2oi j
SEP j jdsij afj
`,
			column:    1,
			sortFlags: SortFlags{sortByMonth: true, removeRepeating: true},
			separator: " ",
		}, `
123 12 3 12 3123 
123 12j3 i33 132ji 
1K 234 21j3 k12j
22 3k21j3k21j 3lkj12io3 j
2G 23 12 3 dfska j
2K ijados fja f 
4G idj foiaj dfj
4M cv ajkjad jfda
djfi aioj adfsja 23j4 jidsfa
JAN 2 3klj 21j djsfoi
JAN 2j 3k1j2oi j
SEP j jdsij afj
DEC j2i3 k12j 3ij`, "Month sorting without duplicated"},
	}

	for _, test := range testCases {
		t.Run(test.hint, func(t *testing.T) {
			testedOutput := SortFile(test.input.file, test.input.column, test.input.sortFlags, test.input.separator)
			// if err != nil {
			// 	t.Fatal("Error while unpacking: ", err)
			// }

			if testedOutput != test.output {

				t.Fatalf("Wrong output:\nexpected %v,\nrecieved %v\nstr1:%s\nstr2:%s", []rune(test.output), []rune(testedOutput), test.output, testedOutput)
			}
		})
	}
}

func TestSortCheck(t *testing.T) {
	testCases :=
		[]testCase[checkFileInput, bool]{
			{checkFileInput{
				file: `
12399 323 
23 9923
90909
98 100
ss sadsf`,
				column:    1,
				sortFlags: SortFlags{checkSorted: true},
				separator: " ",
			}, true, "Check basically sorted file"},
		}
	for _, test := range testCases {
		t.Run(test.hint, func(t *testing.T) {
			testedOutput := SortFile(test.input.file, test.input.column, test.input.sortFlags, test.input.separator)
			// if err != nil {
			// 	t.Fatal("Error while unpacking: ", err)
			// }

			if testedOutput == "Файл не отсортирован" && test.output || testedOutput == "Файл отсортирован" && !test.output {

				t.Fatalf("Wrong output:\nexpected %v,\nrecieved %v\n", test.output, testedOutput)
			}
		})
	}
}

func BenchmarkSortFile(b *testing.B) {
	file := `123 12 3 12 3123 

djfi aioj adfsja 23j4 jidsfa
123 12j3 i33 132ji 
22 3k21j3k21j 3lkj12io3 j
1K 234 21j3 k12j
2G 23 12 3 dfska j
2G 23 12 3 dfska j
2G 23 12 3 dfska j
4M cv ajkjad jfda
4G idj foiaj dfj
2K ijados fja f 
JAN 2 3klj 21j djsfoi
DEC j2i3 k12j 3ij
JAN 2j 3k1j2oi j
SEP j jdsij afj
`
	column := 1
	sortFlags := SortFlags{sortByMonth: true, removeRepeating: true}
	separator := " "

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		SortFile(file, column, sortFlags, separator)
	}
}
