package raft

import (
	"errors"
	"sync"
)

// InmemStore implements the LogStore and StableStore interface.
// It should NOT EVER be used for production. It is used only for
// unit tests. Use the MDBStore implementation instead.
type InmemStore struct {
	lock      sync.RWMutex
	lowIndex  uint64
	highIndex uint64
	logs      map[uint64]*Log
	kv        map[string][]byte
	kvInt     map[string]uint64
}

// NewInmemStore returns a new in-memory backend. Do not ever
// use for production. Only for testing.
func NewInmemStore() *InmemStore {
	i := &InmemStore{
		logs:  make(map[uint64]*Log),
		kv:    make(map[string][]byte),
		kvInt: make(map[string]uint64),
	}
	return i
}

// FirstIndex implements the LogStore interface.
func (i *InmemStore) FirstIndex() (uint64, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.lowIndex, nil
}

// LastIndex implements the LogStore interface.
func (i *InmemStore) LastIndex() (uint64, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.highIndex, nil
}

// GetLog implements the LogStore interface.
func (i *InmemStore) GetLog(index uint64) (*Log, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	l, ok := i.logs[index]
	if !ok {
		return nil, errors.New("not found")
	}
	return l, nil
}

// StoreLog implements the LogStore interface.
func (i *InmemStore) StoreLog(log *Log) error {
	return i.StoreLogs([]*Log{log})
}

// StoreLogs implements the LogStore interface.
func (i *InmemStore) StoreLogs(logs []*Log) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	for _, l := range logs {
		i.logs[l.Index] = l
		if i.lowIndex == 0 {
			i.lowIndex = l.Index
		}
		if l.Index > i.highIndex {
			i.highIndex = l.Index
		}
	}
	return nil
}

// DeleteRange implements the LogStore interface.
func (i *InmemStore) DeleteRange(min, max uint64) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	for j := min; j <= max; j++ {
		delete(i.logs, j)
	}
	i.lowIndex = max + 1
	return nil
}

// Set implements the StableStore interface.
func (i *InmemStore) Set(key []byte, val []byte) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.kv[string(key)] = val
	return nil
}

// Get implements the StableStore interface.
func (i *InmemStore) Get(key []byte) ([]byte, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.kv[string(key)], nil
}

// SetUint64 implements the StableStore interface.
func (i *InmemStore) SetUint64(key []byte, val uint64) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	i.kvInt[string(key)] = val
	return nil
}

// GetUint64 implements the StableStore interface.
func (i *InmemStore) GetUint64(key []byte) (uint64, error) {
	i.lock.RLock()
	defer i.lock.RUnlock()
	return i.kvInt[string(key)], nil
}
