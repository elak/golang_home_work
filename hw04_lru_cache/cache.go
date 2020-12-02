package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	mutex    sync.Mutex
	capacity int               // - ёмкость (количество сохраняемых в кэше элементов)
	queue    List              // - очередь \[последних используемых элементов\] на основе двусвязного списка
	items    map[Key]*ListItem // - словарь, отображающий ключ (строка) на элемент очереди

}
type cacheItem struct {
	key   Key
	value interface{}
}

// Добавить значение в кэш по ключу.
func (cache *lruCache) Set(key Key, value interface{}) bool {
	var newItem cacheItem
	newItem.key = key
	newItem.value = value

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cachedValue, wasInCache := cache.items[key]

	if wasInCache {
		cache.queue.MoveToFront(cachedValue)
		cachedValue.Value = newItem
		return true
	}

	if cache.capacity == cache.queue.Len() {
		delete(cache.items, cache.queue.Back().Value.(cacheItem).key)
		cache.queue.Remove(cache.queue.Back())
	}

	cache.items[key] = cache.queue.PushFront(newItem)

	return false
}

// Получить значение из кэша по ключу.
func (cache *lruCache) Get(key Key) (interface{}, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cachedValue, wasInCache := cache.items[key]

	if !wasInCache {
		return nil, false
	}

	cache.queue.MoveToFront(cachedValue)

	return cachedValue.Value.(cacheItem).value, true
}

// Очистить кэш.
func (cache *lruCache) Clear() {
	cache.queue = NewList()
	cache.items = make(map[Key]*ListItem, cache.capacity)
}

func NewCache(capacity int) Cache {
	var cache lruCache

	cache.capacity = capacity
	cache.Clear()

	return &cache
}
