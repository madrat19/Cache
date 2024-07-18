package main

import (
	"container/list"
	"sync"
	"time"
)

// Интерфейс из условия
type ICache interface {
	Cap() int
	Len() int
	Clear()
	Add(key, value interface{})
	AddWithTTL(key, value interface{}, ttl time.Duration)
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
}

// Основная структура
type Cache struct {
	cap   int
	data  map[interface{}]*list.Element
	queue *list.List
	mutex sync.RWMutex
}

// Элемент в кэше
type Item struct {
	Key        interface{}
	Value      interface{}
	Expiration *time.Time
}

// Конструктор
func NewCache(cap int) *Cache {
	return &Cache{
		cap:   cap,
		data:  make(map[interface{}]*list.Element),
		queue: list.New(),
	}
}

func (c *Cache) Cap() int {
	return c.cap
}

func (c *Cache) Len() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.data)
}

func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[interface{}]*list.Element)
	c.queue = list.New()
}

func (c *Cache) Add(key, value interface{}) {
	c.addInternal(key, value, nil)
}

func (c *Cache) AddWithTTL(key, value interface{}, ttl time.Duration) {
	expiration := time.Now().Add(ttl)
	c.addInternal(key, value, &expiration)
}

func (c *Cache) addInternal(key, value interface{}, expiration *time.Time) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Если ключ уже существует
	if element, exist := c.data[key]; exist {
		c.queue.MoveToFront(element)
		element.Value.(*Item).Value = value
		element.Value.(*Item).Expiration = expiration
		return
	}

	// Если вышли за cap
	if c.queue.Len() == c.cap {
		c.purge()
	}

	item := &Item{
		Key:        key,
		Value:      value,
		Expiration: expiration,
	}

	element := c.queue.PushFront(item)
	c.data[key] = element
}

func (c *Cache) purge() {
	// Если очередь не пуста
	if element := c.queue.Back(); element != nil {
		item := c.queue.Remove(element).(*Item)
		delete(c.data, item.Key)
	}
}

func (c *Cache) Get(key interface{}) (value interface{}, ok bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	element, exist := c.data[key]
	// Если ключ не существует
	if !exist {
		return nil, false
	}

	expiration := element.Value.(*Item).Expiration
	// Если срок жизни ключа истек
	if expiration != nil && expiration.Before(time.Now()) {
		c.queue.Remove(element)
		delete(c.data, key)
		return nil, false
	}

	c.queue.MoveToFront(element)
	return element.Value.(*Item).Value, true
}

func (c *Cache) Remove(key interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	element, exist := c.data[key]
	// Если ключ существует
	if exist {
		c.queue.Remove(element)
		delete(c.data, key)
	}
}
