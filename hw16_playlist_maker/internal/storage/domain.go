package storage

import (
	"context"
	"time"
)

const (
	ScopeUndefined RestrictionScope = iota
	ScopeGroup
	ScopeCategory
)

const (
	PriorityUndefined FillerPriority = iota
	PriorityAmount
	PriorityDuration
	PriorityOrder
	PriorityRandom
)

type UUID string

// Категория - метка ролика или группы роликов. Используется в фильтрах и ограничения стратегий заполнения.
type Category struct {
	ID    UUID   `json:"id"`
	Title string `json:"title"`
}

// Группа роликов - иерархическая структура роликов.
type Group struct {
	ID         UUID   `json:"id"`
	Title      string `json:"title"`
	ParentID   UUID   `json:"parent_id"`
	Order      int    `json:"order"`
	CategoryID UUID   `json:"category_id"`
}

// Ролики - элемент плейлиста.
type Video struct {
	ID         UUID   `json:"id"`
	Title      string `json:"title"`
	ParentID   UUID   `json:"parent_id"`
	Order      int    `json:"order"`
	CategoryID UUID   `json:"category_id"`
	Duration   int    `json:"duration"`
}

// Стратегия заполнения плейлиста.
type Template struct {
	ID           UUID                  `json:"id"`
	Title        string                `json:"title"`
	StartItems   []TemplateItem        `json:"start_items"`
	Items        []TemplateItem        `json:"items"`
	EndItems     []TemplateItem        `json:"end_items"`
	Restrictions []TemplateRestriction `json:"restrictions"`
}

// Элемент(шаг) стратегии заполнения плейлиста.
type TemplateItem struct {
	ID           UUID                  `json:"id"`
	Title        string                `json:"title"`
	Order        int                   `json:"order"`
	Duration     int                   `json:"duration"`
	Fillers      []TemplateFiller      `json:"fillers"`
	Restrictions []TemplateRestriction `json:"restrictions"`
}

type FillerPriority int

// Настройки подбора для шага заполнения плейлиста.
type TemplateFiller struct {
	ID             UUID           `json:"id"`
	Order          int            `json:"order"`
	CategoryID     UUID           `json:"category_id"`
	AllowRepeat    bool           `json:"allow_repeat"`
	GroupsPriority FillerPriority `json:"groups_priority"`
	VideosPriority FillerPriority `json:"videos_priority"`
}

type RestrictionScope int // на группу/ на категорию

// Ограничения стратегии(шага) заполнения плейлиста.
type TemplateRestriction struct {
	ID         UUID             `json:"id"`
	Title      string           `json:"title"`
	Scope      RestrictionScope `json:"scope"`
	CategoryID UUID             `json:"category_id"`
	GroupID    UUID             `json:"group_id"`
	Duration   int              `json:"duration"`
	Amount     int              `json:"amount"`
}

type History struct {
	VideoID  UUID      `json:"video_id"`
	LastSeen time.Time `json:"last_seen"`
}

type GroupsManager interface {
	Create(ctx context.Context, item Group) error
	Update(ctx context.Context, item Group) error
	Read(ctx context.Context, ID UUID) (*Group, error)
	Delete(ctx context.Context, ID UUID) error
	List(ctx context.Context, filter func(*Group) bool) ([]*Group, error)
}

type VideosManager interface {
	Create(ctx context.Context, item Video) error
	Update(ctx context.Context, item Video) error
	Read(ctx context.Context, ID UUID) (*Video, error)
	Delete(ctx context.Context, ID UUID) error
	List(ctx context.Context, filter func(*Video) bool) ([]*Video, error)
}

type CategoriesManager interface {
	Create(ctx context.Context, item Category) error
	Update(ctx context.Context, item Category) error
	Read(ctx context.Context, ID UUID) (*Category, error)
	Delete(ctx context.Context, ID UUID) error
	List(ctx context.Context, filter func(*Category) bool) ([]*Category, error)
}

type TemplatesManager interface {
	Create(ctx context.Context, item Template) error
	Update(ctx context.Context, item Template) error
	Read(ctx context.Context, ID UUID) (*Template, error)
	Delete(ctx context.Context, ID UUID) error
	List(ctx context.Context, filter func(*Template) bool) ([]*Template, error)
}

type HistoryManager interface {
	Create(ctx context.Context, item History) error
	Update(ctx context.Context, item History) error
	Read(ctx context.Context, ID UUID) (*History, error)
	Delete(ctx context.Context, ID UUID) error
	List(ctx context.Context, filter func(*History) bool) ([]*History, error)
}

type Storage interface {
	Groups() GroupsManager
	Videos() VideosManager
	Categories() CategoriesManager
	Templates() TemplatesManager
	History() HistoryManager
}
