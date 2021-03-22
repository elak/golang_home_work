package storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEventTimeOccupied = errors.New("event time occupied")
	ErrEventMustExist    = errors.New("event must exist")
	ErrEventAlreadyExist = errors.New("event already exist")
)

type Event struct {
	// ID - уникальный идентификатор события (можно воспользоваться UUID).
	ID string

	// Заголовок - короткий текст.
	Title string

	// Дата и время события.
	Date time.Time

	// Длительность события (или дата и время окончания).
	DueDate time.Time

	// Описание события - длинный текст, опционально.
	Description string

	// ID пользователя, владельца события.
	Owner string

	// За сколько времени высылать уведомление.
	NotifyDate time.Time
}

type Storage interface {
	// Создать (событие).
	CreateEvent(ctx context.Context, event Event) error

	// Обновить (ID события, событие).
	UpdateEvent(ctx context.Context, eventID string, eventData Event) error

	// Удалить (ID события).
	DeleteEvent(ctx context.Context, eventID string) error

	// СписокСобытийНаДень (дата).
	ListDayEvents(ctx context.Context, eventsDay time.Time) ([]*Event, error)

	// СписокСобытийНаНеделю (дата начала недели).
	ListWeekEvents(ctx context.Context, eventsWeekStart time.Time) ([]*Event, error)

	// СписокСобытийНaМесяц (дата начала месяца).
	ListMonthEvents(ctx context.Context, eventsMonthStart time.Time) ([]*Event, error)
}
