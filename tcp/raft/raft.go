package raft

import (
	"github.com/lycying/pitydb/tcp/mut"
	"github.com/lycying/pitydb/tcp/mut/codec/typelen"
)

type RaftState uint32

const (
	Follower RaftState = iota
	Candidate
	Leader
	Shutdown
)

type Raft struct {
	cluster  *Cluster
	logStore LogStore
}

type RaftConfig struct {
	Srv   string
	Peers []string
}

func NewRaft(raftCfg *RaftConfig) *Raft {
	rf := &Raft{}
	rf.logStore = NewInmemStore()
	cfg := mut.DefaultConfig()
	cfg.SetCodec(typelen.NewCodec())
	cfg.SetCallback(rf)
	rf.cluster = NewCluster(cfg)
	rf.cluster.InitCluster(raftCfg.Srv, raftCfg.Peers)
	return rf
}

func (r *Raft) OnConnect(c *mut.Conn) {
}
func (r *Raft) OnClose(c *mut.Conn) {
}
func (r *Raft) OnError(c *mut.Conn, err error) {
}
func (r *Raft) OnMessage(c *mut.Conn, p mut.Packet) {
}
