package raft

import "github.com/lycying/pitydb/lib/btree"

type BtreeMemLog struct {
	tree          *btree.BTree
	lastCommitted uint64
	lastApplied   uint64
}

func (a *LogEntry) Less(b btree.Item) bool {
	return a.ID < b.(*LogEntry).ID
}

func (l *BtreeMemLog) Get(id uint64) *LogEntry {
	return l.tree.Get(&LogEntry{ID: id}).(*LogEntry)
}
func (l *BtreeMemLog) GetRange(a, b uint64) []*LogEntry {
	got := make([]*LogEntry, 0)
	l.tree.AscendRange(&LogEntry{ID: a}, &LogEntry{ID: b}, func(a btree.Item) bool {
		got = append(got, a.(*LogEntry))
		return true
	})
	return got
}
func (l *BtreeMemLog) Insert(a *LogEntry) {
	l.tree.ReplaceOrInsert(a)
}
func (l *BtreeMemLog) Inserts(ls []*LogEntry) {
	for _, item := range ls {
		l.tree.ReplaceOrInsert(item)
	}
}
func (l *BtreeMemLog) RemoveRange(a, b uint64) int {
	items := l.GetRange(a, b)
	for _, item := range items {
		l.tree.Delete(item)
	}
	return len(items)
}
func (l *BtreeMemLog) LastCommitted() *LogEntry {
	return l.Get(l.lastCommitted)
}
func (l *BtreeMemLog) LastApplied() *LogEntry {
	return l.Get(l.lastApplied)
}
