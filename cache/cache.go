package cache 

import (
	"sync"
	"time"
)

type CacheItem struct {
	Value interface{}
	Expiration int64
	LastAccess int64
}

type InMemoryCache struct {
	sync.RWMutex
	items map[string]*CacheItem
	size int

}

func NewInMemoryCache(size int) *InMemoryCache {
	return &InMemoryCache{
		items: make(map[string]*CacheItem),
		size: size,
	}
}

func (cache *InMemoryCache) Set(key string, value interface{}, duration time.Duration){
	 cache.Lock()
	 defer cache.Unlock()

	 if len(cache.items) >= cache.size {
		cache.evictItems()
	 }

	 cache.items[key] = &CacheItem{
		Value: value, 
		Expiration: time.Now().Add(duration).UnixNano(),
		LastAccess: time.Now().UnixNano(),
	 }
}

func (cache *InMemoryCache) Get(key string) (interface{}, bool) {
	cache.Lock()
	defer cache.Unlock()

	item, found := cache.items[key]
	if !found {
		return nil, false
	}
	if time.Now().UnixNano() > item.Expiration {
		delete(cache.items, key)
		return nil, false
	}
	item.LastAccess = time.Now().UnixNano()
	return item.Value, true
}

func (cache *InMemoryCache) evictItems() {
	var oldestAccess int64
	var oldestKey string

	for k, v := range cache.items {
		if oldestAccess == 0 || v.LastAccess < oldestAccess {
			oldestAccess = v.LastAccess
			oldestKey = k
		}
	}
	delete(cache.items, oldestKey)
}

func (cache *InMemoryCache) StartEvictionTimer(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			cache.evictExpiredItems()
		}
	}()
}

func (cache *InMemoryCache) evictExpiredItems() {
	cache.Lock()
	defer cache.Unlock()

	now := time.Now().UnixNano()
	for k, v := range cache.items {
		if now > v.Expiration {
			delete(cache.items, k)
		}
	}
}

func (cache *InMemoryCache) Delete(key string) {
	cache.Lock()
	defer cache.Unlock()

	delete(cache.items, key)
}

func (cache *InMemoryCache) Flush() {
	cache.Lock()
	defer cache.Unlock()

	cache.items = make(map[string]*CacheItem)
}