package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T){
	cache := NewInMemoryCache(5)
	key := "keyXYZ"
	value := "valueXYZ"

	cache.Set(key, value, 10*time.Second)

	if gotValue, found := cache.Get(key); !found || gotValue != value {
		t.Errorf("Get() = %v, %v; want %v, %v", gotValue, found, value, true)
	}
}

func TestExpiration(t *testing.T){
	cache := NewInMemoryCache(8)
	key := "keyXYZ"
	value := "valueXYZ"
	cache.Set(key, value, 1*time.Millisecond)

	time.Sleep(2 * time.Millisecond)

	if _, found := cache.Get(key); found{
		t.Errorf("Expected items expire")
	}
}

func TestEviction(t *testing.T){
	cache := NewInMemoryCache(1)
	key1 := "keyXYZ"
	value1 := "valueXYZ"
	key2 := "keyZYX"
	value2 := "valueZYX"
	cache.Set(key1, value1, 10*time.Second)
	cache.Set(key2, value2, 9*time.Second)
	if _, found := cache.Get(key1); found {
		t.Errorf("Expected key1 to be removed")
	}	
}

func TestDelete (t *testing.T) {
	cache := NewInMemoryCache(1)
	key := "keyXYZ"
	value := "valueXYZ"
	cache.Set(key, value, 10*time.Second)

	cache.Delete(key)
	if _, found := cache.Get(key); found {
		t.Errorf("Expected key to be deleted")
	}
}

func TestFlush(t *testing.T) {
	cache := NewInMemoryCache(4)
	cache.Set("keyXYZ", "valueXYZ", 10*time.Second)
	cache.Set("keyZYX", "valueZYX", 10*time.Second)
	cache.Set("keyXYZ", "valueXYZ", 10*time.Second)
	cache.Set("keyZYX", "valueZYX", 10*time.Second)
	cache.Flush()

	if len(cache.items) != 0 {
		t.Errorf("Expected cache to be empty, still have %d items", len(cache.items))
	}
}

func TestAutomaticEvictionFunction (t *testing.T){
	cache := NewInMemoryCache(5)
	cache.Set("keyXYZ", "valueXYZ", 1*time.Millisecond)
	time.Sleep(10 * time.Millisecond)

	cache.Lock()
	if len(cache.items) != 0 {
		cache.Unlock()
		t.Errorf("Expected expired items to be removed from cache")
	}
}

func TestCacheSize(t *testing.T){
	cacheSize := 2
	cache := NewInMemoryCache(cacheSize)

	for i := 0; i < cacheSize + 1; i++ {
		key:= fmt.Sprintf("key%d", i)
		cache.Set(key, fmt.Sprintf("value%d", i) , 1*time.Minute)
	}
	if len(cache.items) > cacheSize {
		t.Errorf("Cache size limit exceeded")
	}
}