package app

import (
	"math/rand"
	"sort"
	"time"

	domain "github.com/elak/golang_home_work/hw16_playlist_maker/internal/storage"
)

type PlayList struct {
	Content    []*domain.Video
	contentIdx map[domain.UUID]*domain.Video

	lastSeenByGroup      map[domain.UUID]time.Time
	totalDurationByGroup map[domain.UUID]int
	totalAmountByGroup   map[domain.UUID]int

	durationByGroup    map[domain.UUID]int
	durationByCategory map[domain.UUID]int
	amountByGroup      map[domain.UUID]int
	amountByCategory   map[domain.UUID]int

	template *domain.Template
	videos   []*domain.Video
	groups   map[domain.UUID]*domain.Group
	history  map[domain.UUID]time.Time
}

func (list *PlayList) add(video *domain.Video, catID domain.UUID) {
	list.Content = append(list.Content, video)
	list.contentIdx[video.ID] = video

	list.durationByCategory[catID] += video.Duration
	list.amountByCategory[catID]++

	list.durationByGroup[video.ParentID] += video.Duration
	list.amountByGroup[video.ParentID]++
}

func (list *PlayList) processHistory() {
	for _, video := range list.videos {
		seen, hasRec := list.history[video.ID]
		if !hasRec {
			continue
		}

		list.totalDurationByGroup[video.ParentID] += video.Duration
		list.totalAmountByGroup[video.ParentID]++

		if seen.After(list.lastSeenByGroup[video.ParentID]) {
			list.lastSeenByGroup[video.ParentID] = seen
		}
	}
}

func NewPlaylist(videos []*domain.Video, groups map[domain.UUID]*domain.Group, history map[domain.UUID]time.Time, tpl domain.Template) *PlayList {
	list := PlayList{}

	list.contentIdx = make(map[domain.UUID]*domain.Video)

	list.lastSeenByGroup = make(map[domain.UUID]time.Time)
	list.totalDurationByGroup = make(map[domain.UUID]int)
	list.totalAmountByGroup = make(map[domain.UUID]int)

	list.durationByGroup = make(map[domain.UUID]int)
	list.durationByCategory = make(map[domain.UUID]int)
	list.amountByGroup = make(map[domain.UUID]int)
	list.amountByCategory = make(map[domain.UUID]int)

	list.template = &tpl
	list.videos = videos
	list.groups = groups
	list.history = history

	list.processHistory()

	return &list
}

func MakePlaylist(duration int, videos []*domain.Video, groups map[domain.UUID]*domain.Group, history map[domain.UUID]time.Time, tpl domain.Template) *PlayList {
	list := NewPlaylist(videos, groups, history, tpl)

	rest := duration

	// обязательный начальный блок
	rest = list.processTemplateItemsBlock(rest, tpl.StartItems)

	reservedForLastChunk := 0

	// резервируем время под обязательный финальный блок.
	for _, item := range tpl.EndItems {
		if item.Duration > rest {
			reservedForLastChunk += rest
			rest = 0
			break
		} else {
			reservedForLastChunk += item.Duration
			rest -= item.Duration
		}
	}

	rest = list.processTemplateItemsBlock(rest, tpl.Items)

	list.processTemplateItemsBlock(rest+reservedForLastChunk, tpl.EndItems)

	return list
}

func (list *PlayList) processTemplateItemsBlock(duration int, items []domain.TemplateItem) int {
	rest := duration

	for rest > 60 {
		for i := range items {
			rest = list.processTemplateItem(rest, &items[i])

			if rest <= 60 {
				return rest
			}
		}
	}

	return rest
}

func (list *PlayList) processTemplateItem(duration int, item *domain.TemplateItem) int {
	itemDuration := item.Duration
	duration -= itemDuration

	if duration < 0 {
		itemDuration += duration
		duration = 0
	}

	chunkRest := list.fillPlaylistChunk(itemDuration, item)

	return duration + chunkRest
}

func (list *PlayList) fillPlaylistChunk(duration int, tpl *domain.TemplateItem) int {
	rest := duration
	// в этот список добавляем те же видео, что и в основной для ведения статистики и обсчёта ограничений текущего шага заполнения
	chunkList := PlayList{}

	for _, filler := range tpl.Fillers {
		// Отбираем видео из категории.
		categoryVideos := list.filterVideos(filler.CategoryID, rest, filler.AllowRepeat)

		// Сортируем в порядке приоритетов.
		list.presortVideos(categoryVideos, filler)

		for _, video := range categoryVideos {
			// в блоке осталось достаточно мета для видео?
			if video.Duration-rest > 60 {
				continue
			}

			// проверить ограничения шага
			if !chunkList.checkRestrictions(video, filler.CategoryID, tpl.Restrictions) {
				continue
			}

			// проверить ограничения всего шаблона
			if !list.checkRestrictions(video, filler.CategoryID, list.template.Restrictions) {
				continue
			}

			chunkList.add(video, filler.CategoryID)
			list.add(video, filler.CategoryID)

			rest -= video.Duration
			if rest <= 60 {
				return rest
			}
		}
	}

	return rest
}

func (list *PlayList) checkRestrictions(video *domain.Video, catID domain.UUID, restrictions []domain.TemplateRestriction) bool {
	for i := range restrictions {
		if !list.checkRestriction(video, catID, &restrictions[i]) {
			return false
		}
	}

	return true
}

func (list *PlayList) checkRestriction(video *domain.Video, catID domain.UUID, restriction *domain.TemplateRestriction) bool {
	switch restriction.Scope {
	case domain.ScopeUndefined:
		return true
	case domain.ScopeCategory:
		return list.checkCategoryRestriction(video, catID, restriction)
	case domain.ScopeGroup:
		return list.checkGroupRestriction(video, restriction)
	}

	return true
}

func (list *PlayList) checkGroupRestriction(video *domain.Video, restriction *domain.TemplateRestriction) bool {
	if video.ParentID == "" {
		return true
	}

	if restriction.GroupID != "" && restriction.GroupID != video.ParentID {
		return true
	}

	if restriction.Duration != 0 && list.durationByGroup[video.ParentID]+video.Duration > restriction.Duration {
		return false
	}
	if restriction.Amount != 0 && list.amountByGroup[video.ParentID] == restriction.Amount {
		return false
	}
	return true
}

func (list *PlayList) checkCategoryRestriction(video *domain.Video, catID domain.UUID, restriction *domain.TemplateRestriction) bool {
	if restriction.CategoryID != "" && restriction.CategoryID != catID {
		return true
	}
	if restriction.Duration != 0 && list.durationByCategory[catID]+video.Duration > restriction.Duration {
		return false
	}
	if restriction.Amount != 0 && list.amountByCategory[catID] == restriction.Amount {
		return false
	}
	return true
}

func (list *PlayList) compareVideosGroups(iParentID, jParentID domain.UUID) bool {
	if iParentID == "" {
		return true
	}

	if jParentID == "" {
		return false
	}
	// сравниваем разные группы
	// суммарное время просмотра по группе

	iDuration := list.totalDurationByGroup[iParentID]
	jDuration := list.totalDurationByGroup[jParentID]
	if iDuration != jDuration {
		return iDuration > jDuration
	}

	// суммарное количество просмотров по группе

	iAmount := list.totalAmountByGroup[iParentID]
	jAmount := list.totalAmountByGroup[jParentID]
	if iAmount != jAmount {
		return iAmount > jAmount
	}

	// последний просмотр группы

	iSeen := list.lastSeenByGroup[iParentID]
	jSeen := list.lastSeenByGroup[jParentID]

	if iSeen != jSeen {
		return jSeen.Before(iSeen)
	}

	iGroup := list.groups[iParentID]
	jGroup := list.groups[jParentID]
	// для подгрупп одной группы - порядок в группе
	if iGroup.ParentID == jGroup.ParentID {
		if iGroup.Order != jGroup.Order {
			return iGroup.Order < jGroup.Order
		}
	}

	// и если совсем ничего не осталось порядок идентификаторов
	return iParentID < jParentID
}

func (list *PlayList) compareVideos(iVid, jVid *domain.Video, extPriority map[domain.UUID]int) bool {
	iSeen, iInHistory := list.history[iVid.ID]
	jSeen, jInHistory := list.history[jVid.ID]

	// видео не входят в историю просмотров или входят в один и тот же момент времени
	if iSeen == jSeen {
		// задано внешнее упорядочивание
		if extPriority != nil {
			return extPriority[iVid.ID] < extPriority[jVid.ID]
		}

		// сравниваем порядок
		return iVid.Order < jVid.Order
	}

	// ещё не просмотренные видео всегда выше в списке
	if !iInHistory {
		return true
	}
	if !jInHistory {
		return false
	}

	// сравниваем положение в истории просмотров
	return iSeen.Before(jSeen)
}

func (list *PlayList) presortVideos(videos []*domain.Video, options domain.TemplateFiller) {
	var videosPriority map[domain.UUID]int

	if options.VideosPriority == domain.PriorityRandom {
		videosPriority = make(map[domain.UUID]int)

		rand.Seed(time.Now().UnixNano())

		for _, video := range list.videos {
			videosPriority[video.ID] = rand.Int()
		}
	}

	sort.SliceStable(videos, func(i, j int) bool {
		iVid := videos[i]
		jVid := videos[j]

		// сравниваем элементы одной группы
		if iVid.ParentID == jVid.ParentID {
			return list.compareVideos(iVid, jVid, videosPriority)
		}

		return list.compareVideosGroups(iVid.ParentID, jVid.ParentID)
	})
}

func (list *PlayList) filterVideos(catID domain.UUID, maxDuration int, allowAlreadySeen bool) (res []*domain.Video) {
	res = make([]*domain.Video, 0)

	for _, video := range list.videos {
		if video.Duration <= 0 {
			continue
		}

		if video.Duration > maxDuration {
			continue
		}

		if !allowAlreadySeen {
			if _, seen := list.history[video.ID]; seen {
				continue
			}
		}

		if list.contentIdx[video.ID] != nil {
			continue
		}

		videoCategoryID := video.CategoryID
		if videoCategoryID == "" && list.groups[video.ParentID] != nil {
			videoCategoryID = list.groups[video.ParentID].CategoryID
		}

		if videoCategoryID == "" {
			continue
		}

		if catID == videoCategoryID {
			res = append(res, video)
		}
	}

	return
}
