package pattern

/*
	Реализовать паттерн «фабричный метод».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/Factory_method_pattern
*/

/*
Ответ:

Фабричный метод используется для того, чтобы вынести создание различных объектов с определённым интерфейсом в отдельное место.
Особенно это полезно когда спосов взаимодействия между продуктами не очень чёткий, поскольку данный паттерн может легко расширять свой функционал.

Плюсы:
 - Избегается связь между создателем и объектами
 - Создание продуктов выделяется в отдельное место
 - Легко вводить новые типы продуктов
Минусы:
 - Код разрастается, поскольку нужно много новых классов
*/

type ClothesCreator interface {
	MakeClothes(kind string) Clothes
}

type RealClothesCreator struct{}

type Clothes interface {
	Wear()
}

type Pants struct {
	size int
}

func (p Pants) Wear() {}

type Hat struct {
	size string
}

func (h Hat) Wear() {}

func (r RealClothesCreator) MakeClothes(kind string) Clothes {
	switch kind {
	case "hat":
		return Hat{"21"}
	case "pants":
		return Pants{21}
	default:
		return nil
	}
}
