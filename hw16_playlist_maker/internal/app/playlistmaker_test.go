package app

import (
	"testing"
	"time"

	domain "github.com/elak/golang_home_work/hw16_playlist_maker/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatchVideo(t *testing.T) {
	categories := make(map[string]*domain.Category)
	groups := make(map[domain.UUID]*domain.Group)
	videosByID := make(map[domain.UUID]*domain.Video)

	videos := dataSetUp(categories, groups, videosByID)

	history := make(map[domain.UUID]time.Time)
	seen := time.Now()

	history["GyV3rxnxIQY"] = seen
	history["387zMh6Zvig"] = seen.Add(time.Minute)

	list := NewPlaylist(videos, groups, history, domain.Template{})

	conditions := make([]domain.TemplateRestriction, 1)

	conditions[0] = domain.TemplateRestriction{
		ID:         "0",
		Title:      "Подходи по одному",
		Scope:      domain.ScopeGroup,
		CategoryID: "",
		GroupID:    "",
		Duration:   0,
		Amount:     1,
	}
	// первое видео из группы проходит проверку
	res := list.checkRestrictions(videos[1], categories["Познавательные"].ID, conditions)
	assert.True(t, res)

	// второе видео из группы Не проходит проверку
	list.add(videos[1], categories["Познавательные"].ID)

	res = list.checkRestrictions(videos[2], categories["Познавательные"].ID, conditions)
	assert.False(t, res)

	// второе видео из категории проходит проверку
	res = list.checkRestrictions(videosByID["XPwAl48IIH0"], categories["Познавательные"].ID, conditions)
	assert.True(t, res)

	// меняем область ограничений
	conditions[0].Scope = domain.ScopeCategory

	// второе видео из категории НЕ проходит проверку
	res = list.checkRestrictions(videosByID["XPwAl48IIH0"], categories["Познавательные"].ID, conditions)
	assert.False(t, res)

	// меняем ограничение на конкретную категорию
	conditions[0].CategoryID = categories["Развлекательные"].ID

	// второе видео из другой категории проходит проверку
	res = list.checkRestrictions(videosByID["XPwAl48IIH0"], categories["Познавательные"].ID, conditions)
	assert.True(t, res)

	list.add(videos[0], categories["Развлекательные"].ID)

	// меняем ограничение с количества на длительность
	conditions[0].Amount = 0
	conditions[0].Duration = 10000

	// видео НЕ превышает ограничений длительности
	res = list.checkRestrictions(videos[1], categories["Развлекательные"].ID, conditions)
	assert.True(t, res)

	conditions[0].Duration = 5000
	// видео превышает ограничения длительности
	res = list.checkRestrictions(videos[1], categories["Развлекательные"].ID, conditions)
	assert.False(t, res)
}

func TestFilterVideos(t *testing.T) {
	list := PlayList{}

	videosByID := make(map[domain.UUID]*domain.Video)
	categories := make(map[string]*domain.Category)
	groups := make(map[domain.UUID]*domain.Group)

	list.videos = dataSetUp(categories, groups, videosByID)

	list.history = make(map[domain.UUID]time.Time)
	seen := time.Now()
	// Домики - Золотые ворота
	list.history["GyV3rxnxIQY"] = seen
	// Домики - Храм Неба
	list.history["387zMh6Zvig"] = seen.Add(time.Minute)

	// фильтр на категорию
	filteredVideos := list.filterVideos(categories["Познавательные"].ID, 600, true)
	require.Equal(t, 8, len(filteredVideos))

	// фильтр на категорию и НЕ вхождение в историю просмотра
	filteredVideos = list.filterVideos(categories["Познавательные"].ID, 600, false)
	require.Equal(t, 6, len(filteredVideos))

	// фильтр на категорию и максимально возможную длительность
	filteredVideos = list.filterVideos(categories["Познавательные"].ID, 400, true)
	require.Equal(t, 7, len(filteredVideos))
}

func TestPresortVideos(t *testing.T) {
	videosByID := make(map[domain.UUID]*domain.Video)
	categories := make(map[string]*domain.Category)
	groups := make(map[domain.UUID]*domain.Group)

	videos := dataSetUp(categories, groups, videosByID)

	list := NewPlaylist(videos, groups, nil, domain.Template{})
	list.presortVideos(videos, domain.TemplateFiller{})

	require.Equal(t, "Малыш Коала (The Outback) Мультфильм HD", videos[0].Title)

	list.history = make(map[domain.UUID]time.Time)
	seen := time.Now()
	// Домики - Золотые ворота
	list.history["GyV3rxnxIQY"] = seen
	// Домики - Храм Неба
	list.history["387zMh6Zvig"] = seen.Add(time.Minute)

	videos2 := videos[1:8]

	list.presortVideos(videos2, domain.TemplateFiller{})

	// просмотренные должны уйти в конец списка
	// чем раньше был просмотрен ролик - тем выше он должен быть в части  списка просмотренных

	require.Equal(t, "Домики - Храм Неба – Обучающий мультфильм для детей - про Китай", videos2[6].Title)
	require.Equal(t, "Домики - Золотые ворота – Обучающий мультфильм для детей - про Владимир", videos2[5].Title)
}

func dataSetUp(categories map[string]*domain.Category, groupsByID map[domain.UUID]*domain.Group, videosByID map[domain.UUID]*domain.Video) []*domain.Video {
	categories["Развлекательные"] = &domain.Category{
		ID:    "1",
		Title: "Развлекательные",
	}

	categories["Познавательные"] = &domain.Category{
		ID:    "2",
		Title: "Познавательные",
	}

	categories["Музыкальные"] = &domain.Category{
		ID:    "3",
		Title: "Музыкальные",
	}

	groups := make(map[string]*domain.Group)

	groups["Полномеражки"] = &domain.Group{
		ID:         "1",
		Title:      "Полномеражки",
		ParentID:   "",
		CategoryID: categories["Развлекательные"].ID,
	}

	groups["Домики"] = &domain.Group{
		ID:         "2",
		Title:      "Домики",
		ParentID:   "",
		CategoryID: categories["Познавательные"].ID,
	}

	groups["Биографии"] = &domain.Group{
		ID:         "3",
		Title:      "Биографии",
		ParentID:   "",
		CategoryID: categories["Познавательные"].ID,
	}

	groups["Капитан Краб"] = &domain.Group{
		ID:         "4",
		Title:      "Капитан Краб",
		ParentID:   "",
		CategoryID: categories["Музыкальные"].ID,
	}

	videos := make([]*domain.Video, 0)

	videos = append(videos, &domain.Video{Order: 0, CategoryID: categories["Развлекательные"].ID, ParentID: groups["Полномеражки"].ID, ID: "6ZPOFjSBn-g", Title: "Малыш Коала (The Outback) Мультфильм HD", Duration: 4923})

	videos = append(videos, &domain.Video{Order: 6, CategoryID: categories["Познавательные"].ID, ParentID: groups["Домики"].ID, ID: "EZ9tKEs_Fcs", Title: "Домики - Кинкаку -Дзи – Обучающий мультфильм для детей - Япония", Duration: 336})
	videos = append(videos, &domain.Video{Order: 5, CategoryID: categories["Познавательные"].ID, ParentID: groups["Домики"].ID, ID: "hDB5S6hYrZ8", Title: "Домики - Дом музыки - Обучающий мультфильм для детей - Москва", Duration: 336})
	videos = append(videos, &domain.Video{Order: 4, CategoryID: categories["Познавательные"].ID, ParentID: groups["Домики"].ID, ID: "KUrXteVPWZI", Title: "Домики - Башня Сююмбике – Обучающий мультфильм для детей - про Казань", Duration: 344})
	videos = append(videos, &domain.Video{Order: 3, CategoryID: categories["Познавательные"].ID, ParentID: groups["Домики"].ID, ID: "387zMh6Zvig", Title: "Домики - Храм Неба – Обучающий мультфильм для детей - про Китай", Duration: 355})
	videos = append(videos, &domain.Video{Order: 2, CategoryID: categories["Познавательные"].ID, ParentID: groups["Домики"].ID, ID: "MC6fFNPc6rg", Title: "Домики - Кижи – Обучающий мультфильм для детей - про Карелию", Duration: 345})
	videos = append(videos, &domain.Video{Order: 1, CategoryID: categories["Познавательные"].ID, ParentID: groups["Домики"].ID, ID: "ZPEYwbkUnE4", Title: "Домики - Мост Понте-Веккьо – Обучающий мультфильм для детей - про Италию", Duration: 343})
	videos = append(videos, &domain.Video{Order: 0, CategoryID: categories["Познавательные"].ID, ParentID: groups["Домики"].ID, ID: "GyV3rxnxIQY", Title: "Домики - Золотые ворота – Обучающий мультфильм для детей - про Владимир", Duration: 347})
	videos = append(videos, &domain.Video{CategoryID: categories["Познавательные"].ID, ParentID: groups["Биографии"].ID, ID: "XPwAl48IIH0", Title: "Веселые биографии -  Архимед  – обучающий мультфильм для детей", Duration: 403})
	videos = append(videos, &domain.Video{CategoryID: categories["Музыкальные"].ID, ParentID: groups["Капитан Краб"].ID, ID: "IE5Goo2XCV0", Title: "Капитан Краб: Мечты\" (Колыбельная для детей)\"", Duration: 415})

	for _, video := range videos {
		videosByID[video.ID] = video
	}

	for _, group := range groups {
		groupsByID[group.ID] = group
	}

	return videos
}
