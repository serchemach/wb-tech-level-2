Что выведет программа? Объяснить вывод программы.

```go
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func asChan(vs ...int) <-chan int {
	c := make(chan int)

	go func() {
		for _, v := range vs {
			c <- v
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}

		close(c)
	}()
	return c
}

func merge(a, b <-chan int) <-chan int {
	c := make(chan int)
	go func() {
		for {
			select {
			case v := <-a:
				c <- v
			case v := <-b:
				c <- v
			}
		}
	}()
	return c
}

func main() {

	a := asChan(1, 3, 5, 7)
	b := asChan(2, 4 ,6, 8)
	c := merge(a, b )
	for v := range c {
		fmt.Println(v)
	}
}
```

Ответ:

Для начала стоит заметить, что поскольку функция merge не закрывает канал, который возвращает, после работы, цикл по этому каналу никогда не остановится.

Но до того момента, как горутины уйдут в дедлок, выведутся числа 
```
1, 3, 5, 7, 2, 4, 6, 8
```

В неопределённом порядке (не только потому что числа посылаются в канал со случайной задержкой, но и потому, что они посылаются одновременно из разных горутин)

