package raft

import (
	"encoding/json"
	"github.com/lycying/mut"
	"github.com/lycying/mut/codec/typelen"
	"github.com/lycying/pitydb/dt"
	"math/rand"
	"time"
)

type RaftState uint32

const (
	Follower RaftState = iota
	Candidate
	Leader
	Shutdown
)

type Raft struct {
	cluster        *Cluster
	logStore       LogStore
	state          RaftState
	cfg            *RaftConfig
	heartBeatTimer *time.Timer

	// The current term, cache of StableStore
	currentTerm uint64

	// Highest committed log entry
	commitIndex uint64

	// Last applied log to the FSM
	lastApplied uint64

	myLeader    string
	voteCounter *dt.UInt32
}

type RaftConfig struct {
	Srv   string
	Peers []string
}

func NewRaft(raftCfg *RaftConfig) *Raft {
	rf := &Raft{}

	cfg := mut.DefaultConfig()
	cfg.SetCodec(typelen.NewCodec())
	cfg.SetCallback(rf)
	rf.cfg = raftCfg

	rf.logStore = NewInmemStore()
	log0 := &Log{
		Index: 1,
		Term:  1,
		Data:  "",
	}
	rf.logStore.StoreLog(log0)

	rf.lastApplied = 1
	rf.commitIndex = 1
	rf.currentTerm = 1

	rf.myLeader = ""
	rf.voteCounter = dt.NewUInt32()

	rf.cluster = NewCluster(cfg)
	rf.cluster.InitCluster(raftCfg.Srv, raftCfg.Peers)
	return rf
}
func (r *Raft) Startup() {
	r.SetState(Follower)
	r.WaitForHeartBeat()
}

func (r *Raft) SetState(state RaftState) {
	r.state = state
}

//heartbeat timeout -> elect
//half than elect -> leader
//leader -> send heartbeat

func heartBeatSpec() time.Duration {
	return time.Millisecond*100 + time.Millisecond*time.Duration(rand.Int63n(100))
}
func (r *Raft) WaitForHeartBeat() {
	if nil == r.heartBeatTimer {
		r.heartBeatTimer = time.AfterFunc(heartBeatSpec(), func() {
			idx, _ := r.logStore.LastIndex()
			lastLog, _ := r.logStore.GetLog(idx)
			req := &VoteReq{
				Candidate:    r.cfg.Srv,
				Term:         r.currentTerm,
				LastLogIndex: lastLog.Index,
				LastLogTerm:  lastLog.Term,
			}
			p, _ := Marshal(req)

			r.cluster.Broadcast(p)

			r.myLeader = r.cfg.Srv
			r.voteCounter.SetValue(uint32(1))
		})
	} else {
		r.heartBeatTimer.Reset(heartBeatSpec())
	}
}

func (r *Raft) processVoteReq(c *mut.Conn, req *VoteReq) {
	logger.Debug("%+v", req)
	resp := &VoteResp{}
	p, _ := Marshal(resp)
	c.WriteAsync(p)
}
func (r *Raft) processVoteResp(c *mut.Conn, resp *VoteResp) {
	logger.Debug("%+v", resp)
}

func (r *Raft) OnConnect(c *mut.Conn) {
}
func (r *Raft) OnClose(c *mut.Conn) {
}
func (r *Raft) OnError(c *mut.Conn, err error) {
}
func (r *Raft) OnMessage(c *mut.Conn, p mut.Packet) {
	packet := p.(*typelen.Packet)
	switch packet.Type {
	case TypeVoteReq:
		req := &VoteReq{}
		json.Unmarshal(packet.Data, req)
		r.processVoteReq(c, req)
	case TypeVoteResp:
		resp := &VoteResp{}
		json.Unmarshal(packet.Data, resp)
		r.processVoteResp(c, resp)
	}

}
