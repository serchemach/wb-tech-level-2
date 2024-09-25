package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type IncorrectParametersError struct{}

func (e IncorrectParametersError) Error() string {
	return "Некорректные параметры для метода."
}

func ParseInt(vals url.Values, valName string) (int, error) {
	valUnparsed := vals.Get(valName)
	if valUnparsed == "" {
		return 0, IncorrectParametersError{}
	}

	val, err := strconv.Atoi(valUnparsed)
	if err != nil {
		return 0, IncorrectParametersError{}
	}
	return val, nil
}

func ParseDate(vals url.Values, valName string) (time.Time, error) {
	valUnparsed := vals.Get("date")
	if valUnparsed == "" {
		return time.Time{}, IncorrectParametersError{}
	}

	val, err := time.Parse("2006-01-02", valUnparsed)
	if err != nil {
		return time.Time{}, IncorrectParametersError{}
	}

	return val, nil
}

func ParseEvent(vals url.Values) (Event, error) {
	userId, err := ParseInt(vals, "user_id")
	if err != nil {
		return Event{}, err
	}

	date, err := ParseDate(vals, "date")
	if err != nil {
		return Event{}, IncorrectParametersError{}
	}

	title := vals.Get("title")
	if title == "" {
		return Event{}, IncorrectParametersError{}
	}

	description := vals.Get("description")
	if description == "" {
		return Event{}, IncorrectParametersError{}
	}

	return Event{
		userId:      userId,
		date:        date,
		title:       title,
		description: description,
	}, nil
}

func ParseUpdatingEvent(vals url.Values) (Event, error) {
	event, err := ParseEvent(vals)
	if err != nil {
		return event, err
	}

	eventId, err := ParseInt(vals, "event_id")
	if err != nil {
		return Event{}, err
	}
	event.eventId = eventId

	return event, nil
}

func ParseDeletingArgs(vals url.Values) (int, int, error) {
	userId, err := ParseInt(vals, "user_id")
	if err != nil {
		return 0, 0, err
	}
	eventId, err := ParseInt(vals, "event_id")
	if err != nil {
		return 0, 0, err
	}

	return userId, eventId, nil
}

func ParseFetchingArgs(vals url.Values) (int, time.Time, error) {
	userId, err := ParseInt(vals, "user_id")
	if err != nil {
		return 0, time.Time{}, err
	}

	date, err := ParseDate(vals, "date")
	if err != nil {
		return 0, time.Time{}, err
	}

	return userId, date, nil
}

func EventToJson(event Event) string {
	return fmt.Sprintf(`{"user_id": %d, "date": "%s", "title": "%s", "description": "%s", "event_id": %d}`, event.userId, event.date.Format("2006-01-02"), event.title, event.description, event.eventId)
}

func EventListToJson(events []Event) string {
	var b strings.Builder
	b.WriteString("[")
	for i, event := range events {
		if i != 0 {
			b.WriteString(", ")
		}
		b.WriteString(EventToJson(event))
	}
	b.WriteString("]")
	return b.String()
}

func FormError(s string) string {
	return fmt.Sprintf(`{"error": "%s"}`, s)
}

func FormSuccessfulResponse(s string) string {
	// Обычно это жсоновые объекты, поэтому можно не оборачивать в кавычки
	return fmt.Sprintf(`{"result": %s}`, s)
}

func CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.PostForm)

	event, err := ParseEvent(r.PostForm)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(FormError(err.Error())))
		return
	}

	event = eventStorage.AddNewEvent(event)
	w.WriteHeader(200)
	w.Write([]byte(FormSuccessfulResponse(EventToJson(event))))
}

func UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	event, err := ParseUpdatingEvent(r.PostForm)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(FormError(err.Error())))
		return
	}

	event, err = eventStorage.UpdateEvent(event)
	if err != nil {
		w.WriteHeader(503)
		w.Write([]byte(FormError(err.Error())))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(FormSuccessfulResponse(EventToJson(event))))
}

func DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, eventId, err := ParseDeletingArgs(r.PostForm)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(FormError(err.Error())))
		return
	}

	event, err := eventStorage.DeleteEvent(userId, eventId)
	if err != nil {
		w.WriteHeader(503)
		w.Write([]byte(FormError(err.Error())))
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(FormSuccessfulResponse(EventToJson(event))))
}

func EventsForDayHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, date, err := ParseFetchingArgs(r.Form)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(FormError(err.Error())))
		return
	}

	events := eventStorage.GetEvents(userId, date, date)
	w.WriteHeader(200)
	w.Write([]byte(FormSuccessfulResponse(EventListToJson(events))))
}

func EventsForWeekHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, date, err := ParseFetchingArgs(r.Form)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(FormError(err.Error())))
		return
	}

	weekStart := date.AddDate(0, 0, -int(date.Weekday())+1)
	weekEnd := date.AddDate(0, 0, 7-int(date.Weekday()))
	fmt.Println(date, weekStart, weekEnd)

	events := eventStorage.GetEvents(userId, weekStart, weekEnd)
	w.WriteHeader(200)
	w.Write([]byte(FormSuccessfulResponse(EventListToJson(events))))
}

func EventsForMonthHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId, date, err := ParseFetchingArgs(r.Form)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(FormError(err.Error())))
		return
	}

	monthStart := date.AddDate(0, 0, -date.Day()+1)
	monthEnd := date.AddDate(0, 1, -date.Day())
	fmt.Println(date, monthStart, monthEnd)

	events := eventStorage.GetEvents(userId, monthStart, monthEnd)
	w.WriteHeader(200)
	w.Write([]byte(FormSuccessfulResponse(EventListToJson(events))))
}
