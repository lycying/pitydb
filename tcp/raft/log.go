package raft

import "sync"

const (
	committed   int32 = 0
	uncommitted int32 = 1
)

type LogEntry struct {
	ID    uint64
	Term  uint64
	Cmd   string
	State int32
}

type Log interface {
	Get(uint64) *LogEntry
	GetRange(uint64, uint64) []*LogEntry
	Insert(*LogEntry)
	Inserts([]*LogEntry)
	RemoveRange(uint64, uint64) int
	LastCommitted() *LogEntry
	LastApplied() *LogEntry
}
type RaftLog struct {
	logs        map[uint64]*LogEntry
	commitIndex uint64 //from 1
	lastApplied uint64

	lock *sync.RWMutex
}

func NewLogEntry() *LogEntry {
	return &LogEntry{
		ID:    0,
		Term:  0,
		Cmd:   "",
		State: committed,
	}
}

func NewRaftLog() *RaftLog {
	rl := &RaftLog{
		logs: make(map[uint64]*LogEntry),
		lock: new(sync.RWMutex),
	}
	//default init
	e := NewLogEntry()
	rl.insert(e)
	return rl
}

func (rl *RaftLog) containsEntry(index uint64, term uint64) bool {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	if value, ok := rl.logs[index]; ok {
		return value.Term == term
	}
	return false
}

func (rl *RaftLog) getEntry(index uint64) *LogEntry {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	if value, ok := rl.logs[index]; ok {
		return value
	}
	return nil
}

func (rl *RaftLog) getRange(begin uint64, end uint64) []*LogEntry {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	if nil != rl.logs {
		ret := make([]*LogEntry, 0)
		for _, item := range rl.logs {
			if item.ID > begin && item.ID <= end {
				ret = append(ret, item)
			}
		}
		return ret
	}
	return nil
}

func (rl *RaftLog) discard(id uint64) {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	for k, _ := range rl.logs {
		if id >= k {
			delete(rl.logs, k)
		}
	}
}
func (rl *RaftLog) insert(item *LogEntry) {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	rl.logs[item.ID] = item

	if item.ID > rl.commitIndex {
		rl.commitIndex = item.ID
	}
}

func (rl *RaftLog) inserts(items []*LogEntry) {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	if nil != items {
		for _, item := range items {
			rl.logs[item.ID] = item

			if item.ID > rl.commitIndex {
				rl.commitIndex = item.ID
			}
		}
	}
}

func (rl *RaftLog) lastTerm() uint64 {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	r := uint64(0)
	if nil != rl.logs {
		for _, v := range rl.logs {
			if v.Term > r {
				r = v.Term
			}
		}
	}
	return r
}
func (rl *RaftLog) lastID() uint64 {
	rl.lock.Lock()
	defer rl.lock.Unlock()

	r := uint64(0)
	if nil != rl.logs {
		for k, _ := range rl.logs {
			if k > r {
				r = k
			}
		}
	}
	return r
}
