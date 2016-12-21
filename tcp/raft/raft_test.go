package raft

import (
	"testing"
	"time"
)

func TestRaft(t *testing.T) {
	raftCfg1 := &RaftConfig{}
	raftCfg1.Srv = ":14001"
	raftCfg1.Peers = []string{"localhost:14001", "localhost:14002"}
	raft1 := NewRaft(raftCfg1)
	println(raft1)

	time.Sleep(time.Hour)
}
