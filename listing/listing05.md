Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
} 

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error
```

По той же причине, что и в вопросе 3. Просто тут неявное преобразование из customError в error из-за того что он стоит как тип возвращаемой переменной.

По сути, код test() эквивалентен

``` go
func test() *customError {
	var cusErr *customError
	{
		// do something
	}
	cusErr = nil
	return cusErr
} 
```

Просто здесь преобразование в error происходит не на этапе возврата, а на этапе присваивания к переменной типа error.

