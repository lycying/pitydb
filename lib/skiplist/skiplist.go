package skiplist

import (
	"fmt"
	"math/rand"
	"sync"
)

const SkipListMaxLevel = 32
const SkipListBranch = 4

type iter struct {
	current *node
	list    *SkipList
	k, v    interface{}
}

type SkipList struct {
	update []*column
	rank   []int
	header *column
	tail   *column
	length int
	level  int

	lessThan func(l, r interface{}) bool
	sync.RWMutex
}

func New(lessThan func(l, r interface{}) bool) *SkipList {
	sk := new(SkipList)
	sk.length = 0
	sk.level = 1
	sk.header = newColumn(SkipListMaxLevel)
	sk.tail = newColumn(0)
	sk.rank = make([]int, SkipListMaxLevel)
	sk.update = make([]*column, SkipListMaxLevel)
	sk.lessThan = func(l, r interface{}) bool {
		if l == nil {
			return true
		}
		return lessThan(l, r)
	}

	for i := 0; i < SkipListMaxLevel; i++ {
		sk.header.forward[i].ptr = sk.tail
	}
	return sk
}

func randomLevel() int {
	level := 1
	for (rand.Int31()&0xFFFF)%SkipListBranch == 0 {
		level += 1
	}

	if level < SkipListMaxLevel {
		return level
	} else {
		return SkipListMaxLevel
	}
}

func (sk *SkipList) Find(k interface{}) interface{} {
	sk.RLock()
	defer sk.RUnlock()

	it := sk.find(k)
	if it.forward[0].ptr == sk.tail {
		return nil
	}
	if it.forward[0].ptr.k == k {
		return it.forward[0].ptr.v
	}
	return nil
}

func (sk *SkipList) find(k interface{}) *column {
	it := sk.header
	for i := sk.level - 1; i >= 0; i-- {
		for it.forward[i].ptr != sk.tail && sk.lessThan(it.forward[i].ptr.k, k) {
			it = it.forward[i].ptr
		}
	}
	return it
}

func (sk *SkipList) Add(k, v interface{}) {
	sk.Lock()
	defer sk.Unlock()

	it := sk.header
	for i := sk.level - 1; i >= 0; i-- {
		if i == sk.level-1 {
			sk.rank[i] = 0
		} else {
			sk.rank[i] = sk.rank[i+1]
		}
		for it.forward[i].ptr != sk.tail && sk.lessThan(it.forward[i].ptr.k, k) {
			sk.rank[i] = sk.rank[i] + it.forward[i].span
			it = it.forward[i].ptr
		}
		sk.update[i] = it
	}

	//get a random level
	level := randomLevel()

	if level > sk.level {
		for i := sk.level; i < level; i++ {
			sk.update[i] = sk.header
			sk.rank[i] = 0
			sk.update[i].forward[i].span = sk.length
		}
		sk.level = level
	}

	col := newColumn(level)
	col.k = k
	col.v = v

	for i := 0; i < level; i++ {
		colPre := sk.update[i]
		colNext := colPre.forward[i].ptr
		colPre.forward[i].ptr = col
		col.forward[i].ptr = colNext

	}

	sk.length++

}

func (sk *SkipList) Del(k interface{}) bool {
	sk.Lock()
	defer sk.Unlock()

	it := sk.header
	for i := sk.level - 1; i >= 0; i-- {
		if i == sk.level-1 {
			sk.rank[i] = 0
		} else {
			sk.rank[i] = sk.rank[i+1]
		}
		for it.forward[i].ptr != sk.tail && sk.lessThan(it.forward[i].ptr.k, k) {
			sk.rank[i] = sk.rank[i] + it.forward[i].span
			it = it.forward[i].ptr
		}
		sk.update[i] = it
	}

	if it.forward[0].ptr.k == k {
		for i := 0; i < sk.level; i++ {
			colPre := sk.update[i]
			if i >= len(it.forward) {
				break
			}
			nextCol := it.forward[i].ptr
			colPre.forward[i].ptr = nextCol.forward[i].ptr
		}
		sk.length--
		return true
	} else {
		return false
	}
}

func (sk *SkipList) Len() int {
	sk.RLock()
	defer sk.RUnlock()

	return sk.length
}

func (sk *SkipList) Level() int {
	sk.RLock()
	defer sk.RUnlock()

	return sk.level
}

func (sk *SkipList) Dump() {
	sk.Lock()
	defer sk.Unlock()

	for i := 0; i < SkipListMaxLevel; i++ {
		if sk.header.forward[i].ptr != sk.tail {
			fmt.Print("##Level", i+1, "##")
			it := sk.header.forward[i].ptr
			for it != sk.tail {
				if it.k != nil {
					fmt.Print(it.k.(int), ",")
				}
				it = it.forward[i].ptr
			}
			fmt.Println()
		}
	}
}
