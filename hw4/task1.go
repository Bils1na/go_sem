package main

import (
	"fmt"
 	"time"
)

type Cache interface {
	Get(k string) (string, bool)
	Set(k, v string)
}

var _ Cache = (*cacheImpl)(nil)

// Доработаем конструктор и методы кеша, чтобы они соответствовали интерфейсу Cache
func newCacheImpl() *cacheImpl {
 	return &cacheImpl{data: make(map[string]cacheItem)}
}

type cacheItem struct {
	value     string
	timestamp time.Time
}

type cacheImpl struct {
 	data map[string]cacheItem
}

func (c *cacheImpl) Get(k string) (string, bool) {
	item, found := c.data[k]
	if !found {
		return "", false
	}
	// Удаляем старые ключи
	if time.Since(item.timestamp) > time.Minute {
		delete(c.data, k)
		return "", false
	}
	return item.value, true
}

func (c *cacheImpl) Set(k, v string) {
 	c.data[k] = cacheItem{value: v, timestamp: time.Now()}
}

func newDbImpl(cache Cache) *dbImpl {
 	return &dbImpl{cache: cache, dbs: map[string]string{"hello": "world", "test": "test"}}
}

type dbImpl struct {
	cache Cache
	dbs   map[string]string
}

func (d *dbImpl) Get(k string) (string, bool) {
	v, ok := d.cache.Get(k)
	if ok {
		return fmt.Sprintf("answer from cache: key: %s, val: %s", k, v), ok
	}

	v, ok = d.dbs[k]
	if ok {
		d.cache.Set(k, v)
		return fmt.Sprintf("answer from dbs: key: %s, val: %s", k, v), ok
	}

	return "", false
}

func main() {
	c := newCacheImpl()
	db := newDbImpl(c)
	fmt.Println(db.Get("test"))
	fmt.Println(db.Get("hello"))
	fmt.Println(db.Get("nonexistent"))
}