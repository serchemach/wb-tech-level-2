package pattern

/*
	Реализовать паттерн «состояние».
Объяснить применимость паттерна, его плюсы и минусы, а также реальные примеры использования данного примера на практике.
	https://en.wikipedia.org/wiki/State_pattern
*/

/*
Ответ:

Данный паттерн удобно использовать, когда объект имеет ограниченное количество состояний, в каждом из которых он имеет разное поведение.
В таком случае можно вынести поведение в отдельный класс состояния и хранить в объекте именно его, вместо того, чтобы плодить огромную логику по выбору подходящего поведения.

Плюсы:
 - Код становится проще без огромных условий выбора поведения
 - Изолируется код, относящийся к поведению в определённом состоянии
Минусы:
 - Для объектов с маленьким количеством состояний может быть чересчур
*/

type TimeOfYear interface {
	GetPrecipitation() string
}

type Weather struct {
	state TimeOfYear
}

func (w *Weather) Precipitate() string {
	return w.state.GetPrecipitation()
}

func NewWeather(t TimeOfYear) Weather {
	return Weather{t}
}

type Winter struct{}

func (w *Winter) GetPrecipitation() string {
	return "snow"
}

type Summer struct{}

func (s *Summer) GetPrecipitation() string {
	return "rain"
}
