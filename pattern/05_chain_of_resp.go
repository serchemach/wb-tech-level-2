package pattern

/*
	Реализовать паттерн «цепочка вызовов».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Chain-of-responsibility_pattern
*/

/*
Ответ:
Цепь вызовов (ответственности) удобно использовать когда имеется несколько последовательно идущих друг за другом операций, которые могут решить, что запрос по той или иной причине не нужно дальше обрабатывать и закончить цепочку.

Плюсы:
 - Можно контролировать порядок обработки запросов
 - Можно отделить классы, которые вызывают операции, от классов, которые их выполняют
Минусы:
 - Некоторые запросы могут остаться необработанными

*/

type Handler interface {
	Handle(input string) (string, error)
}

type AuthHandler struct {
	next Handler
}

type AuthError struct{}

func (err AuthError) Error() string {
	return "No auth"
}

func (h AuthHandler) Handle(input string) (string, error) {
	if input != "secret!" {
		return "", AuthError{}
	}

	if h.next != nil {
		return h.next.Handle(input)
	}

	return "Cool", nil
}

type InitHandler struct {
	next Handler
}

type InitError struct{}

func (err InitError) Error() string {
	return "Can't init"
}

func (h InitHandler) Handle(input string) (string, error) {
	if input != "good" {
		return "", InitError{}
	}

	if h.next != nil {
		return h.next.Handle(input)
	}
	return "GOOOOD!", nil
}
