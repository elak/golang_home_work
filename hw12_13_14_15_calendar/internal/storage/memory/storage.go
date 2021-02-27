package memorystorage

import (
	"github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"sync"
)

type Storage struct {
	// TODO
	mu sync.RWMutex
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) CreateEvent(event storage.Event) error {
	// TODO
	return nil
}

// TODO
