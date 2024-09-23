package pattern

import "fmt"

/*
	Реализовать паттерн «комманда».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Command_pattern
*/

/*
Ответ:

Данный паттерн используется когда мы хотим ввести абстракцию для операций.
Самый наглядный пример - это кнопки в гуи. Вместо создания отдельного класса для каждой кнопки в интерфейсе, можно отделить от него бизнес-логику, вынеся её в отдельный класс комманд.

Можно ввести также так называемого Исполнителя (Invoker), который будет управлять исполнением команд.

Плюсы:
 - Разделение обязанностей
 - Удобнее работать с действиями (добавить историю, сделать отмену и т.д.)
 - Можно без проблем вводить новые операции
Минусы
 - Новый слой абстракции между интерфейсом и логикой программы
*/

type Command interface {
	Run()
}

type SubmitCommand struct {
	toSubmit string
}

func (c SubmitCommand) Run() {
	fmt.Println(c.toSubmit)
}

type ConvertCommand struct {
	toConvert string
	converted []rune
}

func (c *ConvertCommand) Run() {
	c.converted = []rune(c.toConvert)
}

type Button struct {
	action Command
}

func main() {
	submit := Button{
		action: SubmitCommand{"123"},
	}

	submit.action.Run()

	convert := Button{
		action: &ConvertCommand{"123", nil},
	}

	convert.action.Run()

}
