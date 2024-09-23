package pattern

import "fmt"

/*
	Реализовать паттерн «посетитель».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Visitor_pattern
*/

/*
Ответ:

Паттерн посетитель используется для того, чтобы обходить набор объектов с разнородными интерфейсами и благодаря тому, что код обработки этих объектов лежит в посетителе, мы можем добавлять новый функционал, не изменяя код других классов.

Плюсы:
 - Можно добавлять новый функционал классам без их изменения
 - Можно группировать однотипные действия над разными классами вместе
Минусы:
 - В большинстве ООП языков у посетителя не будет доступа к приватным полям класса


*/

type Building interface {
	accept(v Visitor)
}

type School struct{}

func (s School) accept(v Visitor) {
	v.visitSchool(s)
}

type Shop struct{}

func (s Shop) accept(v Visitor) {
	v.visitShop(s)
}

type Visitor interface {
	visitSchool(s School)
	visitShop(s Shop)
}

type RecorderVisitor struct {
	visitedBuildings string
}

func (v *RecorderVisitor) visitSchool(s School) {
	if len(v.visitedBuildings) == 0 {
		v.visitedBuildings += "School"
		return
	}
	v.visitedBuildings += " School"
}

func (v *RecorderVisitor) visitShop(s Shop) {
	if len(v.visitedBuildings) == 0 {
		v.visitedBuildings += "Shop"
		return
	}
	v.visitedBuildings += " Shop"
}

type CounterVisitor struct {
	numOfSchools int
	numOfShops   int
}

func (v *CounterVisitor) visitSchool(s School) {
	v.numOfSchools++
}

func (v *CounterVisitor) visitShop(s Shop) {
	v.numOfSchools++
}

func main() {
	buildings := []Building{School{}, Shop{}, Shop{}, Shop{}}
	var (
		cv CounterVisitor
		rv RecorderVisitor
	)

	for _, building := range buildings {
		building.accept(&cv)
		building.accept(&rv)
	}

	fmt.Println(cv.numOfSchools, cv.numOfShops)
	fmt.Println(rv.visitedBuildings)
}
