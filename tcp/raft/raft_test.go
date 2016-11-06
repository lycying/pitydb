package raft

import (
	"fmt"
	"testing"
	"time"
)

func TestRaft_Start(t *testing.T) {
	raft1 := NewRaft("0.0.0.0:14001")
	go raft1.Start()

	raft2 := NewRaft("0.0.0.0:14002")
	go raft2.Start()

	raft3 := NewRaft("0.0.0.0:14003")
	go raft3.Start()

	//send append log
	go func() {
		shotme_count := 1
		for {
			time.Sleep(time.Second * 2)
			var rf *Raft
			if raft1.IsLeader() {
				rf = raft1
			}

			if raft2.IsLeader() {
				rf = raft2
			}

			if raft3.IsLeader() {
				rf = raft3
			}

			if nil != rf {
				shotme_count++
				rf.AddCmd(fmt.Sprintf("shot me %d", shotme_count))
			}
		}
	}()

	time.Sleep(time.Second * 60)
}
