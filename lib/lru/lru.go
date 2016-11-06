package lru

import (
	"container/list"
	"sync"
)

// Lru is a simple lru and is threadsafe
type Lru struct {
	cap int
	l   *list.List
	m   map[interface{}]*list.Element
	sync.RWMutex
}

type lruEntry struct {
	k interface{}
	v interface{}
}

func newLruEntry(k, v interface{}) *lruEntry {
	return &lruEntry{
		k: k,
		v: v,
	}

}

func New(cap int) *Lru {
	lru := new(Lru)
	lru.cap = cap
	lru.l = list.New()
	lru.m = make(map[interface{}]*list.Element)

	return lru
}

func (lru *Lru) Set(k, v interface{}) bool {
	lru.Lock()
	defer lru.Unlock()

	if ele, ok := lru.m[k]; ok {
		ele.Value.(*lruEntry).v = v
		lru.l.MoveToFront(ele)
	} else {
		entry := newLruEntry(k, v)
		ele = lru.l.PushFront(entry)
		lru.m[k] = ele
	}

	if lru.l.Len() > lru.cap {
		ele := lru.l.Back()
		lru.l.Remove(ele)
		delete(lru.m, ele.Value.(*lruEntry).k)
		return false
	}
	return true

}

func (lru *Lru) Get(k interface{}) interface{} {
	lru.RLock()
	defer lru.RUnlock()

	if ele, ok := lru.m[k]; ok {
		return ele.Value.(*lruEntry).v
	}

	return nil

}

func (lru *Lru) Remove(k interface{}) bool {
	lru.Lock()
	defer lru.Unlock()

	if ele, ok := lru.m[k]; ok {
		lru.l.Remove(ele)
		delete(lru.m, k)
		return true
	}
	return false
}

func (lru *Lru) Size() int {
	return lru.l.Len()
}
