package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

/*
=== HTTP server ===

Реализовать HTTP сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP библиотекой.
В рамках задания необходимо:
	1. Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
	2. Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
	3. Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
	4. Реализовать middleware для логирования запросов
Методы API: POST /create_event POST /update_event POST /delete_event GET /events_for_day GET /events_for_week GET /events_for_month
Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON документ содержащий либо {"result": "..."} в случае успешного выполнения метода,
либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
	1. Реализовать все методы.
	2. Бизнес логика НЕ должна зависеть от кода HTTP сервера.
	3. В случае ошибки бизнес-логики сервер должен возвращать HTTP 503. В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400. В случае остальных ошибок сервер должен возвращать HTTP 500. Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
	4. Код должен проходить проверки go vet и golint.
*/

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {

}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {

}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {

}

func EventsForDayHandler(w http.ResponseWriter, r *http.Request) {

}

func EventsForWeekHandler(w http.ResponseWriter, r *http.Request) {

}

func EventsForMonthHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	l := log.New(os.Stdout, "", 0)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /create_event", CreateEventHandler)
	mux.HandleFunc("POST /update_event", UpdateEventHandler)
	mux.HandleFunc("POST /delete_event", DeleteEventHandler)
	mux.HandleFunc("GET /events_for_day", EventsForDayHandler)
	mux.HandleFunc("GET /events_for_week", EventsForWeekHandler)
	mux.HandleFunc("GET /events_for_month", EventsForMonthHandler)
	wrappedMux := NewLoggerMiddleware(mux, l)

	log.Fatal(http.ListenAndServe(":8080", wrappedMux))
}
