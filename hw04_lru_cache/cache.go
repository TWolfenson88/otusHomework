package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc lruCache) Set(key Key, value interface{}) bool {
	if lc.capacity == 0 {
		return false
	}

	if _, ok := lc.items[key]; ok {
		lc.items[key].Value.(*cacheItem).value = value
		lc.queue.MoveToFront(lc.items[key])
		lc.items[key] = lc.queue.Front()
		return ok
	}

	lc.queue.PushFront(&cacheItem{
		key:   key,
		value: value,
	})
	lc.items[key] = lc.queue.Front()

	if lc.queue.Len() > lc.capacity {
		lc.Clear()
	}
	return false
}

func (lc lruCache) Get(key Key) (interface{}, bool) {
	if _, ok := lc.items[key]; ok {
		lc.queue.MoveToFront(lc.items[key])
		lc.items[key] = lc.queue.Front()
		return lc.items[key].Value.(*cacheItem).value, ok
	}
	return nil, false
}

func (lc lruCache) Clear() {
	// lastItem := lc.queue.Back()

	if item, ok := lc.queue.Back().Value.(*cacheItem); ok {
		delete(lc.items, item.key)
		lc.queue.Remove(lc.queue.Back())
	}
}
