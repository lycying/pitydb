package btree

import (
	"testing"
)

type Int int

// Less returns true if int(a) < int(b).
func (a Int) Less(b Item) bool {
	return a < b.(Int)
}
func (a Int) Key() int{
	return int(a)
}

func TestBTreeCreate(t *testing.T) {
	tree := New()
	tree.Insert(Int(1000))
}


