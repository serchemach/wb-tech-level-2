package main

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
)

type EventTest struct {
	UserId      int    `json:"user_id"`
	Date        string `json:"date"`
	Title       string `json:"title"`
	Description string `json:"description"`
	EventId     int    `json:"event_id"`
}

type EventResponse struct {
	Result EventTest `json:"result"`
}

type EventListResponse struct {
	Result []EventTest `json:"result"`
}

func TestAPI(t *testing.T) {
	t.Run("Test Create Event", func(t *testing.T) {
		reqBody := "user_id=1&date=2024-09-10&title=hello&description=123"
		req := httptest.NewRequest("POST", "/create_event", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		CreateEventHandler(w, req)

		if w.Result().StatusCode != 200 {
			t.Fatal("Wrong status code", w.Result().StatusCode, ", expected 200")
		}

		decoder := json.NewDecoder(w.Result().Body)
		output := EventResponse{}
		err := decoder.Decode(&output)
		if err != nil {
			t.Fatal("Ошибка при чтении ответа:", err)
		}

		trueOutput := EventResponse{
			EventTest{
				UserId:      1,
				Date:        "2024-09-10",
				Title:       "hello",
				Description: "123",
				EventId:     1,
			},
		}

		if output != trueOutput {
			t.Fatalf("Wrong response\nExpected: %v\nRecieved: %v\n", trueOutput, output)
		}
	})

	t.Run("Test Update Event", func(t *testing.T) {
		reqBody := "user_id=1&date=2024-09-11&title=newhello&description=new&event_id=1"
		req := httptest.NewRequest("POST", "/update_event", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		UpdateEventHandler(w, req)

		if w.Result().StatusCode != 200 {
			t.Fatal("Wrong status code", w.Result().StatusCode, ", expected 200")
		}

		decoder := json.NewDecoder(w.Result().Body)
		output := EventResponse{}
		err := decoder.Decode(&output)
		if err != nil {
			t.Fatal("Ошибка при чтении ответа:", err)
		}

		trueOutput := EventResponse{
			EventTest{
				UserId:      1,
				Date:        "2024-09-11",
				Title:       "newhello",
				Description: "new",
				EventId:     1,
			},
		}

		if output != trueOutput {
			t.Fatalf("Wrong response\nExpected: %v\nRecieved: %v\n", trueOutput, output)
		}
	})

	t.Run("Test Delete Event", func(t *testing.T) {
		reqBody := "user_id=1&event_id=1"
		req := httptest.NewRequest("POST", "/create_event", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		w := httptest.NewRecorder()
		DeleteEventHandler(w, req)

		if w.Result().StatusCode != 200 {
			t.Fatal("Wrong status code", w.Result().StatusCode, ", expected 200")
		}

		decoder := json.NewDecoder(w.Result().Body)
		output := EventResponse{}
		err := decoder.Decode(&output)
		if err != nil {
			t.Fatal("Ошибка при чтении ответа:", err)
		}

		trueOutput := EventResponse{
			EventTest{
				UserId:      1,
				Date:        "2024-09-11",
				Title:       "newhello",
				Description: "new",
				EventId:     1,
			},
		}

		if output != trueOutput {
			t.Fatalf("Wrong response\nExpected: %v\nRecieved: %v\n", trueOutput, output)
		}
	})

	t.Run("Test Fetch Events For Month", func(t *testing.T) {
		// Create a bunch of events first
		trueEvents := [30]EventTest{}
		for i := range 30 {
			reqBody := fmt.Sprintf("user_id=1&date=2024-09-%02d&title=hello&description=123", i+1)
			req := httptest.NewRequest("POST", "/create_event", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			CreateEventHandler(w, req)

			trueEvents[i] = EventTest{
				UserId:      1,
				Date:        fmt.Sprintf("2024-09-%02d", i+1),
				Title:       "hello",
				Description: "123",
				EventId:     i + 2,
			}
		}

		reqBody := fmt.Sprintf("?user_id=1&date=2024-09-05")
		req := httptest.NewRequest("GET", "/events_for_month"+reqBody, nil)
		w := httptest.NewRecorder()
		EventsForMonthHandler(w, req)

		if w.Result().StatusCode != 200 {
			t.Fatal("Wrong status code", w.Result().StatusCode, ", expected 200")
		}

		decoder := json.NewDecoder(w.Result().Body)
		output := EventListResponse{}
		err := decoder.Decode(&output)
		if err != nil {
			t.Fatal("Ошибка при чтении ответа:", err)
		}

		slices.SortFunc(output.Result, func(a, b EventTest) int {
			return a.EventId - b.EventId
		})

		if !slices.Equal(output.Result, trueEvents[:]) {
			t.Fatalf("Wrong response\nExpected: %v\nRecieved: %v\n", trueEvents, output.Result)
		}
	})

	t.Run("Test Fetch Events For Day", func(t *testing.T) {
		// Create a bunch of events first
		trueEvents := [30]EventTest{}
		for i := range 30 {
			reqBody := fmt.Sprintf("user_id=1&date=2024-10-20&title=%d&description=123", i+1)
			req := httptest.NewRequest("POST", "/create_event", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			CreateEventHandler(w, req)

			trueEvents[i] = EventTest{
				UserId:      1,
				Date:        "2024-10-20",
				Title:       fmt.Sprintf("%d", i+1),
				Description: "123",
				EventId:     i + 32,
			}
		}

		reqBody := fmt.Sprintf("?user_id=1&date=2024-10-20")
		req := httptest.NewRequest("GET", "/events_for_day"+reqBody, nil)
		w := httptest.NewRecorder()
		EventsForDayHandler(w, req)

		if w.Result().StatusCode != 200 {
			t.Fatal("Wrong status code", w.Result().StatusCode, ", expected 200")
		}

		decoder := json.NewDecoder(w.Result().Body)
		output := EventListResponse{}
		err := decoder.Decode(&output)
		if err != nil {
			t.Fatal("Ошибка при чтении ответа:", err)
		}

		slices.SortFunc(output.Result, func(a, b EventTest) int {
			return a.EventId - b.EventId
		})

		if !slices.Equal(output.Result, trueEvents[:]) {
			t.Fatalf("Wrong response\nExpected: %v\nRecieved: %v\n", trueEvents, output.Result)
		}
	})

	t.Run("Test Fetch Events For Week", func(t *testing.T) {
		// Create a bunch of events first
		trueEvents := [7]EventTest{}
		for i := range 7 {
			reqBody := fmt.Sprintf("user_id=2&date=2024-11-%02d&title=new&description=123", i+3)
			req := httptest.NewRequest("POST", "/create_event", strings.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			w := httptest.NewRecorder()
			CreateEventHandler(w, req)

			trueEvents[i] = EventTest{
				UserId:      2,
				Date:        fmt.Sprintf("2024-11-%02d", i+3),
				Title:       "new",
				Description: "123",
				EventId:     i + 62,
			}
		}

		reqBody := fmt.Sprintf("?user_id=2&date=2024-11-05")
		req := httptest.NewRequest("GET", "/events_for_week"+reqBody, nil)
		w := httptest.NewRecorder()
		EventsForWeekHandler(w, req)

		if w.Result().StatusCode != 200 {
			t.Fatal("Wrong status code", w.Result().StatusCode, ", expected 200")
		}

		decoder := json.NewDecoder(w.Result().Body)
		output := EventListResponse{}
		err := decoder.Decode(&output)
		if err != nil {
			t.Fatal("Ошибка при чтении ответа:", err)
		}

		slices.SortFunc(output.Result, func(a, b EventTest) int {
			return a.EventId - b.EventId
		})

		if !slices.Equal(output.Result, trueEvents[:]) {
			t.Fatalf("Wrong response\nExpected: %v\nRecieved: %v\n", trueEvents, output.Result)
		}
	})
}
