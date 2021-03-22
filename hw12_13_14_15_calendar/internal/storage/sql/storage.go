package sqlstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New() *Storage {
	return &Storage{}
}

// Обновить (ID события, событие).
func (s *Storage) UpdateEvent(ctx context.Context, eventID string, eventData storage.Event) error {
	qText := "INSERT INTO Events(ID, Title, Date, DueDate, Description, Owner, NotifyDate) VALUES(?, ?, ?, ?, ?, ?, ?) on duplicate KEY update Title=Title, Date=Date, DueDate=DueDate, Description=Description, Owner=Owner, NotifyDate=NotifyDate"
	_, err := s.db.ExecContext(ctx, qText, eventID, eventData.Title, eventData.Date, eventData.DueDate, eventData.Description, eventData.Owner, eventData.NotifyDate)

	return err
}

// Удалить (ID события).
func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	qText := "DELETE FROM Events WHERE ID=?"
	_, err := s.db.ExecContext(ctx, qText, eventID)
	return err
}

// Список событий в интервале дат.
func (s *Storage) doListRangeEvents(ctx context.Context, rangeStart time.Time, rangeEnd time.Time) (res []*storage.Event, err error) {
	qText := "SELECT ID, Title, Date, DueDate, Description, Owner, NotifyDate FROM Events WHERE Date BETWEEN ? and ?"

	rows, err := s.db.QueryContext(ctx, qText, rangeStart, rangeEnd)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		_ = rows.Err()
	}()

	res = make([]*storage.Event, 0)

	for rows.Next() {
		var e storage.Event
		if err := rows.Scan(&e.ID, &e.Title, &e.Date, &e.DueDate, &e.Description, &e.Owner, &e.NotifyDate); err != nil {
			return nil, err
		}

		res = append(res, &e)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// СписокСобытийНаДень (дата).
func (s *Storage) ListDayEvents(ctx context.Context, eventsDay time.Time) ([]*storage.Event, error) {
	return s.doListRangeEvents(ctx, eventsDay, eventsDay.AddDate(0, 0, 1))
}

// СписокСобытийНаНеделю (дата начала недели).
func (s *Storage) ListWeekEvents(ctx context.Context, eventsWeekStart time.Time) ([]*storage.Event, error) {
	return s.doListRangeEvents(ctx, eventsWeekStart, eventsWeekStart.AddDate(0, 0, 7))
}

// СписокСобытийНaМесяц (дата начала месяца).
func (s *Storage) ListMonthEvents(ctx context.Context, eventsMonthStart time.Time) ([]*storage.Event, error) {
	return s.doListRangeEvents(ctx, eventsMonthStart, eventsMonthStart.AddDate(0, 1, 0))
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	return s.UpdateEvent(ctx, event.ID, event)
}

func (s *Storage) Prepare(ctx context.Context, storageURI string) (err error) {
	err = s.Connect(ctx, storageURI)
	if err != nil {
		return
	}

	defer func() { err = s.Close(ctx) }()

	return
}

func (s *Storage) Connect(ctx context.Context, storageURI string) (err error) {
	s.db, err = sql.Open("mysql", storageURI)
	if err != nil {
		return err
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}
