package main

import (
	"sync"
	"time"
)

type Event struct {
	userId      int
	date        time.Time
	title       string
	description string
	eventId     int
}

type UserStorage map[int]Event

type EventStorage struct {
	lastId  int
	storage map[int]UserStorage
	mu      sync.RWMutex
}

func NewEventStorage() *EventStorage {
	return &EventStorage{
		lastId:  0,
		storage: make(map[int]UserStorage),
	}
}

func (s *EventStorage) AddNewEvent(event Event) Event {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.lastId++
	event.eventId = s.lastId
	if val, ok := s.storage[event.userId]; ok {
		val[event.eventId] = event
	} else {
		s.storage[event.userId] = map[int]Event{
			event.eventId: event,
		}
	}

	return event
}

type NoEventError struct{}

func (e NoEventError) Error() string {
	return "События нет в базе."
}

func (s *EventStorage) UpdateEvent(event Event) (Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if userVal, ok := s.storage[event.userId]; ok {
		if _, fine := userVal[event.eventId]; fine {
			userVal[event.eventId] = event
		} else {
			return event, NoEventError{}
		}
	} else {
		return event, NoEventError{}
	}

	return event, nil
}

func pseudoVectorAppend(s []Event, el Event) []Event {
	if len(s) == cap(s) {
		s = append(s, make([]Event, len(s))...)[:len(s)]
	}
	return append(s, el)
}

func (s *EventStorage) GetEvents(userId int, periodStart time.Time, periodEnd time.Time) []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Event, 0)
	for _, event := range s.storage[userId] {
		if !periodStart.After(event.date) && !periodEnd.Before(event.date) {
			result = pseudoVectorAppend(result, event)
		}
	}

	return result
}

func (s *EventStorage) DeleteEvent(userId int, eventId int) (Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var event Event
	if userVal, ok := s.storage[userId]; ok {
		if val, fine := userVal[eventId]; fine {
			event = val
			delete(userVal, eventId)
		} else {
			return Event{}, NoEventError{}
		}
	} else {
		return Event{}, NoEventError{}
	}

	return event, nil
}
