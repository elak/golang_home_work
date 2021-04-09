package common

import (
	"context"

	"github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/elak/golang_home_work/hw12_13_14_15_calendar/internal/storage/sql"
)

type DisposableStorage interface {
	storage.Storage

	Prepare(ctx context.Context, storageURI string) error

	Connect(ctx context.Context, storageURI string) error

	Close(ctx context.Context) error
}

func New(storageType string) DisposableStorage {
	switch storageType {
	case "memory":
		return memorystorage.New()
	case "sql":
		return sqlstorage.New()
	default:
		return nil
	}
}
