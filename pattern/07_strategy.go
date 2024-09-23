package pattern

/*
	Реализовать паттерн «стратегия».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Strategy_pattern
*/

/*
Ответ:

Стратения полезна когда у нас имеется набор алгоритмов, которые выполняют похожие действия (например, алгоритмы сортировки).
Из них выделяется общий интерфейс и они инкапсулируются в отдельные классы, чтобы отделить логику в классе, которые выбирает нужный алгоритм и сами алгоритмы.


Плюсы:
 - Можно менять алгоритмы в объекте во время работы программы
 - можно изолировать алгоритмы от их использования
Минусы:
 - При малом количестве неизменных алгоритмов нецелесообразно
 - Клиенты должны знать о существовании разных стратегий чтобы выбрать необходимую
 - Часто может быть заменено анонимными функциями
*/

type Place int

type Context struct {
	strategy RouteStrategy
}

func (c *Context) Route(a, b Place) {
	c.strategy.MakeRoute(a, b)
}

func (c *Context) SetStragegy(s RouteStrategy) {
	c.strategy = s
}

type RouteStrategy interface {
	MakeRoute(a, b Place)
}

type HungryRouteStrategy struct{}

func (h HungryRouteStrategy) MakeRoute(a, b Place) {}

type ExhaustiveRouteStrategy struct{}

func (h ExhaustiveRouteStrategy) MakeRoute(a, b Place) {}
