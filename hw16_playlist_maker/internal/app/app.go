package app

import (
	"context"
	"time"

	domain "github.com/elak/golang_home_work/hw16_playlist_maker/internal/storage"
)

type App struct { // TODO
	storage domain.Storage
}

func New(storage domain.Storage) *App {
	return &App{storage}
}

func (a *App) CreateGroup(ctx context.Context, item domain.Group) error {
	return a.storage.Groups().Create(ctx, item)
}

func (a *App) UpdateGroup(ctx context.Context, item domain.Group) error {
	return a.storage.Groups().Update(ctx, item)
}

func (a *App) ReadGroup(ctx context.Context, id domain.UUID) (*domain.Group, error) {
	return a.storage.Groups().Read(ctx, id)
}

func (a *App) DeleteGroup(ctx context.Context, id domain.UUID) error {
	return a.storage.Groups().Delete(ctx, id)
}

func (a *App) ListGroup(ctx context.Context) ([]*domain.Group, error) {
	return a.storage.Groups().List(ctx, nil)
}

func (a *App) CreateTemplate(ctx context.Context, item domain.Template) error {
	return a.storage.Templates().Create(ctx, item)
}

func (a *App) UpdateTemplate(ctx context.Context, item domain.Template) error {
	return a.storage.Templates().Update(ctx, item)
}

func (a *App) ReadTemplate(ctx context.Context, id domain.UUID) (*domain.Template, error) {
	return a.storage.Templates().Read(ctx, id)
}

func (a *App) DeleteTemplate(ctx context.Context, id domain.UUID) error {
	return a.storage.Templates().Delete(ctx, id)
}

func (a *App) ListTemplate(ctx context.Context) ([]*domain.Template, error) {
	return a.storage.Templates().List(ctx, nil)
}

func (a *App) CreateCategory(ctx context.Context, item domain.Category) error {
	return a.storage.Categories().Create(ctx, item)
}

func (a *App) UpdateCategory(ctx context.Context, item domain.Category) error {
	return a.storage.Categories().Update(ctx, item)
}

func (a *App) ReadCategory(ctx context.Context, id domain.UUID) (*domain.Category, error) {
	return a.storage.Categories().Read(ctx, id)
}

func (a *App) DeleteCategory(ctx context.Context, id domain.UUID) error {
	return a.storage.Categories().Delete(ctx, id)
}

func (a *App) ListCategory(ctx context.Context) ([]*domain.Category, error) {
	return a.storage.Categories().List(ctx, nil)
}

func (a *App) CreateVideo(ctx context.Context, item domain.Video) error {
	return a.storage.Videos().Create(ctx, item)
}

func (a *App) UpdateVideo(ctx context.Context, item domain.Video) error {
	return a.storage.Videos().Update(ctx, item)
}

func (a *App) ReadVideo(ctx context.Context, id domain.UUID) (*domain.Video, error) {
	return a.storage.Videos().Read(ctx, id)
}

func (a *App) DeleteVideo(ctx context.Context, id domain.UUID) error {
	return a.storage.Videos().Delete(ctx, id)
}

func (a *App) ListVideo(ctx context.Context) ([]*domain.Video, error) {
	return a.storage.Videos().List(ctx, nil)
}

func (a *App) CreateHistory(ctx context.Context, item domain.History) error {
	return a.storage.History().Create(ctx, item)
}

func (a *App) UpdateHistory(ctx context.Context, item domain.History) error {
	return a.storage.History().Update(ctx, item)
}

func (a *App) ReadHistory(ctx context.Context, id domain.UUID) (*domain.History, error) {
	return a.storage.History().Read(ctx, id)
}

func (a *App) DeleteHistory(ctx context.Context, id domain.UUID) error {
	return a.storage.History().Delete(ctx, id)
}

func (a *App) ListHistory(ctx context.Context) ([]*domain.History, error) {
	return a.storage.History().List(ctx, nil)
}

func (a *App) MakePlayList(ctx context.Context, duration int, tlpID domain.UUID) ([]*domain.Video, error) {
	videos, err := a.storage.Videos().List(ctx, nil)
	if err != nil {
		return nil, err
	}

	groups := make(map[domain.UUID]*domain.Group)

	mapGroups := func(Group *domain.Group) bool {
		groups[Group.ID] = Group
		return false
	}

	_, err = a.storage.Groups().List(ctx, mapGroups)
	if err != nil {
		return nil, err
	}

	var history map[domain.UUID]time.Time
	mapHistory := func(record *domain.History) bool {
		history[record.VideoID] = record.LastSeen
		return false
	}
	_, err = a.storage.History().List(ctx, mapHistory)
	if err != nil {
		return nil, err
	}

	tpl, err := a.storage.Templates().Read(ctx, tlpID)
	if err != nil {
		return nil, err
	}

	list := MakePlaylist(duration, videos, groups, history, *tpl)

	return list.Content, nil
}
