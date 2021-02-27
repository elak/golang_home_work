package sqlstorage

import (
	"context"
	"github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	// TODO
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(event storage.Event) error {
	// TODO
	return nil
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}
