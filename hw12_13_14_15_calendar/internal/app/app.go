package app

import (
	"context"

	"github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
	// TODO
}

type Logger interface { // TODO
}

type Storage interface {
	CreateEvent(event storage.Event) error
	// TODO
}

func New(logger Logger, storage Storage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(ctx context.Context, id string, title string) error {
	return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
