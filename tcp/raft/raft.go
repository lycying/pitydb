package raft

import (
	"encoding/json"
	"github.com/lycying/pitydb/dt"
	"github.com/lycying/pitydb/tcp/mut"
	"github.com/lycying/pitydb/tcp/mut/codec/typelen"
	"math"
	"math/rand"
	"time"
)

const (
	stateFollower  = 0
	stateCandidate = 1
	stateLeader    = 2
)

type followerInfo struct {
}

type candidateInfo struct {
	votes *dt.UInt32
}

type leaderInfo struct {
	nextIndex  map[string]uint64
	matchIndex map[string]uint64
}

type Node struct {
	state int

	currentTerm uint64
	votedFor    string
	raftLog     *RaftLog

	leaderInfo
	candidateInfo
	followerInfo
}

type Raft struct {
	clientsStr []string
	serverStr  string
	clients    map[string]*mut.Client
	server     *mut.Server

	electionTimer  *time.Timer
	heartbeatTimer *time.Timer

	Node
}

func NewRaft(str string) *Raft {
	r := &Raft{
		serverStr: str,
		clients:   make(map[string]*mut.Client),
		server:    nil,
		Node: Node{
			state:       stateFollower,
			currentTerm: 0,
			votedFor:    "",
			raftLog:     NewRaftLog(),

			candidateInfo: candidateInfo{
				votes: dt.NewUInt32(),
			},

			leaderInfo: leaderInfo{
				nextIndex:  make(map[string]uint64),
				matchIndex: make(map[string]uint64),
			},
		},
	}
	r.readConfig()
	return r
}

func (rf *Raft) readConfig() {
	rf.clientsStr = []string{
		"0.0.0.0:14001",
		"0.0.0.0:14002",
		"0.0.0.0:14003",
	}
}

func (rf *Raft) Start() {
	cfg := mut.DefaultConfig()
	cfg.SetCodec(typelen.NewCodec())
	cfg.SetCallback(rf)

	rf.startServer(cfg)
	rf.startClients(cfg)
	rf.startTimer()

}

func (rf *Raft) AddCmd(cmd string) {
	logger.Debug("raft# cmd recv : %v", cmd)
	e := NewLogEntry()
	e.Cmd = cmd
	e.Term = rf.currentTerm
	e.ID = rf.raftLog.lastID() + 1
	e.State = uncommitted

	rf.raftLog.insert(e)
}

func (rf *Raft) IsLeader() bool {
	return rf.state == stateLeader
}

func (rf *Raft) startServer(cfg *mut.Config) {
	//the server
	server := mut.NewServer(rf.serverStr, cfg)
	err := server.Servo()
	if err != nil {
		logger.Err(err, "raft# server init error")
		return
	}

	rf.server = server

}
func (rf *Raft) startClients(cfg *mut.Config) {
	for _, item := range rf.clientsStr {
		if item != rf.serverStr {
			client := mut.NewClient(item, cfg)
			go func() {
				err := client.DialAsync()
				if err != nil {
					logger.Err(err, "error while connect")
				}
			}()
			rf.clients[item] = client
		}

	}
}

func (rf *Raft) startTimer() {
	rf.electionTimer = time.NewTimer(0)
	rf.electionTimer.Reset(time.Millisecond * time.Duration(rand.Int63n(200)+100))

	rf.heartbeatTimer = time.NewTimer(0)
	rf.heartbeatTimer.Reset(math.MaxInt64)

	//todo: stop along the raft
	go func() {
		for {
			select {
			case <-rf.electionTimer.C:
				rf.doVoteReq()
			case <-rf.heartbeatTimer.C:
				rf.doAppendEntriesReq()
			}
		}
	}()
}
func (rf *Raft) onlineNumber() int {
	ret := 1
	for _, v := range rf.clients {
		if v.IsConnected() {
			ret++
		}
	}
	return ret
}

func (rf *Raft) toStateLeader() {

}

func (rf *Raft) OnMessage(c *mut.Conn, p mut.Packet) {
	packet := p.(*typelen.Packet)
	switch packet.Type {
	case TypeAppendEntriesReq:
		obj := &AppendEntriesReq{}
		json.Unmarshal(packet.Data, obj)
		rf.onAppendEntriesReq(c, obj)
	case TypeAppendEntriesResp:
		obj := &AppendEntriesResp{}
		json.Unmarshal(packet.Data, obj)
		rf.onAppendEntriesResp(c, obj)
	case TypeVoteReq:
		obj := &VoteReq{}
		json.Unmarshal(packet.Data, obj)
		rf.onVoteReq(c, obj)
	case TypeVoteResp:
		obj := &VoteResp{}
		json.Unmarshal(packet.Data, obj)
		rf.onVoteResp(c, obj)
	}
}
func (rf *Raft) OnConnect(c *mut.Conn) {
	logger.Debug("raft# on connect")
}
func (rf *Raft) OnClose(c *mut.Conn) {
	logger.Debug("raft# on close")
}
func (rf *Raft) OnError(c *mut.Conn, err error) {
	logger.Err(err, "raft# on error %+v,%+v", c.Socket().RemoteAddr(), c.Socket().LocalAddr())
	c.Close()
}
