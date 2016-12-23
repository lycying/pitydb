package raft

import (
	"testing"
	"time"
)

func TestRaft(t *testing.T) {
	go func() {
		raftCfg := &RaftConfig{}
		raftCfg.Srv = ":14001"
		raftCfg.Peers = []string{"localhost:14001", "localhost:14002", "localhost:14003"}
		raft := NewRaft(raftCfg)
		raft.Startup()
		println(raft)
	}()
	go func() {
		raftCfg := &RaftConfig{}
		raftCfg.Srv = ":14002"
		raftCfg.Peers = []string{"localhost:14001", "localhost:14002", "localhost:14003"}
		raft := NewRaft(raftCfg)
		raft.Startup()
		println(raft)
	}()
	go func() {
		raftCfg := &RaftConfig{}
		raftCfg.Srv = ":14003"
		raftCfg.Peers = []string{"localhost:14001", "localhost:14002", "localhost:14003"}
		raft := NewRaft(raftCfg)
		raft.Startup()
		println(raft)
	}()

	time.Sleep(time.Hour)
}
