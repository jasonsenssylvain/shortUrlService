package shorturllib

import (
	"container/list"
	"errors"
	"fmt"
)

//EvictCallback
type EvictCallback func(key interface{}, value interface{})

//GetCallback
type GetCallback func(key interface{}) (interface{}, error)

type LRU struct {
	size    int
	list    *list.List
	items   map[interface{}]*list.Element
	onEvict EvictCallback
	onGet   GetCallback
}

type Entry struct {
	key   interface{}
	value interface{}
}

func NewLRU(itemSize int, evictCallback EvictCallback, getCallback GetCallback) (*LRU, error) {
	if itemSize <= 0 {
		return nil, errors.New("size should be > 0")
	}
	lru := &LRU{
		itemSize,
		list.New(),
		make(map[interface{}]*list.Element),
		evictCallback,
		getCallback,
	}
	return lru, nil
}

func (t *LRU) Clear() {
	for key, value := range t.items {
		if t.onEvict != nil {
			t.onEvict(key, value)
		}
		delete(t.items, key)
	}
	t.list.Init()
}

func (t *LRU) Add(key interface{}, value interface{}) bool {
	item, ok := t.items[key]
	if ok {
		t.list.MoveToFront(item)
		item.Value.(*Entry).value = value
		return false
	}

	entry := &Entry{key, value}
	element := t.list.PushFront(entry)
	t.items[key] = element

	if t.list.Len() > t.size {
		t.removeOldest()
	}

	fmt.Println("add item, key is " + key.(string) + ", value is " + value.(string))
	for k, v := range t.items {
		fmt.Println("key is " + k.(string) + ", value is " + v.Value.(*Entry).value.(string))
	}
	return true
}

func (t *LRU) Get(key interface{}) (interface{}, bool) {
	fmt.Println("get item, key is " + key.(string))
	ele, ok := t.items[key]
	// if ele.Value != nil {
	// 	fmt.Println("value is " + ele.Value.(*Entry).value.(string))
	// }
	if ok {
		t.list.MoveToFront(ele)
		return ele.Value.(*Entry).value, true
	} else if t.onGet != nil {
		fmt.Println("cannot find item, try to get ")
		value, err := t.onGet(key)
		if err != nil {
			return nil, false
		}

		ok = t.Add(key, value)
		if !ok {
			return nil, false
		}

	}
	return nil, false
}

func (t *LRU) GetOldest() (interface{}, interface{}, bool) {
	ele := t.list.Back()
	if ele != nil {
		entry := ele.Value.(*Entry)
		return entry.key, entry.value, true
	}
	return nil, nil, false
}

func (t *LRU) Keys() ([]interface{}, error) {
	keys := make([]interface{}, 0)
	for ele := t.list.Back(); ele != nil; ele = ele.Prev() {
		key := ele.Value.(*Entry).key
		keys = append(keys, key)
	}
	return keys, nil
}

func (t *LRU) Len() int {
	return t.list.Len()
}

func (t *LRU) Contains(key interface{}) bool {
	_, ok := t.items[key]
	return ok
}

func (t *LRU) RemoveElement(ele *list.Element) bool {
	t.list.Remove(ele)
	e := ele.Value.(*Entry)
	delete(t.items, e.key)
	if t.onEvict != nil {
		t.onEvict(e.key, e.value)
	}
	return true
}

func (t *LRU) RemoveKey(key interface{}) bool {
	ele, ok := t.items[key]
	if ok {
		t.list.Remove(ele)
		delete(t.items, key)
	}
	return true
}

func (t *LRU) removeOldest() bool {
	item := t.list.Back()
	if item != nil {
		t.list.Remove(item)
		return true
	}
	return false
}
