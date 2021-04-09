package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	data     map[string]*storage.Event
	idxByDay map[time.Time][]string
	mu       sync.RWMutex
}

func New() *Storage {
	res := Storage{}

	res.data = make(map[string]*storage.Event)
	res.idxByDay = make(map[time.Time][]string)

	return &res
}

func (s *Storage) doUpdateEvent(eventID string, eventData storage.Event, createNew bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldData, evenExists := s.data[eventID]

	if !evenExists && !createNew {
		return storage.ErrEventMustExist
	}

	if evenExists && createNew {
		return storage.ErrEventAlreadyExist
	}

	if !s.checkEventTimeUniq(eventID, eventData.Date) {
		return storage.ErrEventTimeOccupied
	}

	updateIdx := !evenExists

	if evenExists && eventData.Date != oldData.Date {
		eventDay := eventData.Date.Truncate(time.Hour * 24)
		oldDay := oldData.Date.Truncate(time.Hour * 24)
		if eventDay != oldDay {
			updateIdx = true
			s.deleteFromIdx(eventID, oldDay)
		}
	}

	s.data[eventID] = &eventData

	if updateIdx {
		eventDay := eventData.Date.Truncate(time.Hour * 24)
		s.idxByDay[eventDay] = append(s.idxByDay[eventDay], eventID)
	}

	return nil
}

func (s *Storage) checkEventTimeUniq(eventID string, eventDate time.Time) bool {
	bucket := s.idxByDay[eventDate.Truncate(time.Hour*24)]

	for _, id := range bucket {
		if id == eventID {
			continue
		} else if s.data[id].Date == eventDate {
			return false
		}
	}

	return true
}

func (s *Storage) deleteFromIdx(eventID string, oldDay time.Time) {
	bucket := s.idxByDay[oldDay]
	size := len(bucket)

	for i, id := range bucket {
		if id == eventID {
			bucket[i] = bucket[size-1]
			s.idxByDay[oldDay] = bucket[0 : size-1]
			return
		}
	}
}

// Обновить (ID события, событие).
func (s *Storage) UpdateEvent(ctx context.Context, eventID string, eventData storage.Event) error {
	return s.doUpdateEvent(eventID, eventData, false)
}

// Удалить (ID события).
func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	oldData, evenExists := s.data[eventID]
	if !evenExists {
		return storage.ErrEventMustExist
	}

	oldDay := oldData.Date.Truncate(time.Hour * 24)

	s.deleteFromIdx(eventID, oldDay)

	delete(s.data, eventID)

	return nil
}

// СписокСобытийНаДень (дата).
func (s *Storage) doListRangeEvents(rangeStart time.Time, rangeEnd time.Time) ([]*storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]*storage.Event, 0)
	rangeStart = rangeStart.Truncate(time.Hour * 24)
	for rangeStart.Before(rangeEnd) {
		bucket := s.idxByDay[rangeStart]

		for _, id := range bucket {
			res = append(res, s.data[id])
		}

		rangeStart = rangeStart.AddDate(0, 0, 1)
	}

	return res, nil
}

// СписокСобытийНаДень (дата).
func (s *Storage) ListDayEvents(ctx context.Context, eventsDay time.Time) ([]*storage.Event, error) {
	return s.doListRangeEvents(eventsDay, eventsDay.AddDate(0, 0, 1))
}

// СписокСобытийНаНеделю (дата начала недели).
func (s *Storage) ListWeekEvents(ctx context.Context, eventsWeekStart time.Time) ([]*storage.Event, error) {
	return s.doListRangeEvents(eventsWeekStart, eventsWeekStart.AddDate(0, 0, 7))
}

// СписокСобытийНaМесяц (дата начала месяца).
func (s *Storage) ListMonthEvents(ctx context.Context, eventsMonthStart time.Time) ([]*storage.Event, error) {
	return s.doListRangeEvents(eventsMonthStart, eventsMonthStart.AddDate(0, 1, 0))
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	return s.doUpdateEvent(event.ID, event, true)
}

func (s *Storage) Prepare(_ context.Context, _ string) error {
	return nil
}

func (s *Storage) Connect(_ context.Context, _ string) error {
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}
