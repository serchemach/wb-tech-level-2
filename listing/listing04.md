Что выведет программа? Объяснить вывод программы.

```go
package main

func main() {
	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	for n := range ch {
		println(n)
	}
}
```

Ответ:

Поскольку данный код не закрывает канал, по которому мы итерируемся, данный код попадёт в дедлок, но перед этим он выведет следующее:

```
0
1
2
3
4
5
6
7
8
9
```

Порядок гарантирован, поскольку посылаем данные в канал мы только из одной горутины.

