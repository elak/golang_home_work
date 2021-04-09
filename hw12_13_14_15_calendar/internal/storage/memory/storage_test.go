package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	stor := New()

	day1, err := time.Parse("02.01.2006", "31.01.2021")
	require.NoError(t, err)

	day2, err := time.Parse("02.01.2006", "24.02.2021")
	require.NoError(t, err)

	dueDay, err := time.Parse("02.01.2006", "31.12.2021")
	require.NoError(t, err)

	e := storage.Event{
		ID:          "UUID1",
		Title:       "Event 1",
		Date:        day1,
		DueDate:     dueDay,
		Description: "",
		Owner:       "",
		NotifyDate:  dueDay,
	}

	e2 := storage.Event{
		ID:          "UUID2",
		Title:       "Event 2",
		Date:        day1.Add(time.Hour * 12),
		DueDate:     dueDay,
		Description: "",
		Owner:       "",
		NotifyDate:  dueDay,
	}

	// добавление
	err = stor.CreateEvent(ctx, e)
	require.NoError(t, err)

	err = stor.CreateEvent(ctx, e2)
	require.NoError(t, err)

	// поиск
	dayEvents, err := stor.ListDayEvents(ctx, day1)
	require.NoError(t, err)
	require.Equal(t, 2, len(dayEvents))

	monthEvents, err := stor.ListMonthEvents(ctx, day1)
	require.NoError(t, err)
	require.Equal(t, 2, len(monthEvents))

	require.Equal(t, dayEvents, monthEvents)

	// обновление записей и индекса
	e.Date = day2
	err = stor.UpdateEvent(ctx, e.ID, e)
	require.NoError(t, err)

	dayEvents, err = stor.ListDayEvents(ctx, day1)
	require.NoError(t, err)
	require.Equal(t, 1, len(dayEvents))

	e.Date = day1
	err = stor.UpdateEvent(ctx, e.ID, e)
	require.NoError(t, err)

	dayEvents, err = stor.ListDayEvents(ctx, day1)
	require.NoError(t, err)
	require.Equal(t, 2, len(dayEvents))

	// удаление
	err = stor.DeleteEvent(ctx, e2.ID)
	require.NoError(t, err)

	dayEvents, err = stor.ListDayEvents(ctx, day1)
	require.NoError(t, err)
	require.Equal(t, 1, len(dayEvents))

	// контроль ошибок
	err = stor.DeleteEvent(ctx, "UUID3")
	require.ErrorIs(t, storage.ErrEventMustExist, err)

	e2.Date = day1
	err = stor.CreateEvent(ctx, e2)
	require.ErrorIs(t, storage.ErrEventTimeOccupied, err)
}
