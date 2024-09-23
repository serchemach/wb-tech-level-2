package pattern

/*
	Реализовать паттерн «строитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Builder_pattern
*/

/*
Пример:

Данный паттерн позволяет конструировать объекты шаг за шагом вместо того, чтобы создавать их одной огромной функцией с кучей параметров (чем часто занимаются мл либы на питоне кстати).
Самая простая аналогия - строительство машины или дома.
Также к нему могут добавить ещё один класс - Director (директор, руководитель), который управляет последовательностью и набором вызываемых методов строителем для достижения чего-то определённого.

Плюсы:
 - Легче создавать комплексные объекты
 - Изоляция сложной логики создания от бизнес логики
Минусы:
 - Количество классов разрастается, особенно если в случае когда создаются разные типы строителей (в том числе и для одного и того же класса)


*/

type MorningRoutineBuilder struct{}

func (m *MorningRoutineBuilder) MakeDrink(drink string) {}

func (m *MorningRoutineBuilder) WatchTV(channel string) {}

func (m *MorningRoutineBuilder) EatFood(food string) {}

type Director interface {
	Build()
}

type LongMorningDirector struct {
	m MorningRoutineBuilder
}

func (l *LongMorningDirector) Build() {
	l.m.MakeDrink("coffee")
	l.m.EatFood("steak")
	l.m.WatchTV("5")
}

type ShortMorningDirector struct {
	m MorningRoutineBuilder
}

func (l *ShortMorningDirector) Build() {
	l.m.EatFood("pasta")
}
