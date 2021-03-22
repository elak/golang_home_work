package app

import (
	"context"

	"github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage storage.Storage
}

type Logger interface {
	LogMessage(msg string, msgLevel int8)
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)
}

func New(logger Logger, storage storage.Storage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(ctx context.Context, id string, title string) error {
	return a.storage.CreateEvent(ctx, storage.Event{ID: id, Title: title})
}
