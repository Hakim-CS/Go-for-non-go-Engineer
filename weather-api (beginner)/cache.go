package main

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data   WheatherRespond
	Expiry time.Time
}

type Cache struct {
	items map[string]CacheItem
	mu    sync.Mutex
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}
func (c *Cache) Get(key string) (WheatherRespond, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// check whether if key(city) exist or not
	item, exist := c.items[key]
	if !exist {
		return WheatherRespond{}, false
	}
	// what it return : c.Items[key] , true/false
	// if exist is true then item will have the value of c.Items[key]

	// check whether if the item is expired or not
	if time.Now().After(item.Expiry) { // explain this line please
		// if current time is after the expiry time of the item
		// then the item is expired
		delete(c.items, key) // delete the expired item from the cache
		return WheatherRespond{}, false
	}

	return item.Data, true

}

func (c *Cache) Set(key string, data WheatherRespond, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheItem{
		Data:   data,
		Expiry: time.Now().Add(duration),
	}
}
