package lru_test

import (
	"github.com/lycying/pitydb/lib/lru"
	"testing"
)

func TestLru(t *testing.T) {
	l := lru.New(5)
	l.Set("10000", "abc0")
	l.Set("10001", "abc1")
	l.Set("10002", "abc2")
	l.Set("10003", "abc3")
	l.Set("10004", "abc4")
	l.Set("10005", "abc5")

	if l.Get("10000") != nil {
		t.Fail()
	}
	if l.Get("10001").(string) != "abc1" {
		t.Fail()
	}
}
