package raft

import (
	"encoding/json"
	"github.com/lycying/pitydb/tcp/mut/codec/typelen"
)

// AppendEntriesReq Invoked by leader to replicate log entries (§5.3); also used as heartbeat (§5.2).
type AppendEntriesReq struct {
	Term         uint64 //leader’s term
	LeaderId     string //so follower can redirect clients
	PrevLogIndex uint64 //index of log entry immediately preceding new ones
	PrevLogTerm  uint64 //term of prevLogIndex entry
	Entries      []*Log //log entries to store (empty for heartbeat; may send more than one for efficiency)
	LeaderCommit uint64 //leader’s commitIndex
}

type AppendEntriesResp struct {
	Term        uint64 //currentTerm, for leader to update itself
	Index       uint64
	CommitIndex uint64
	Success     bool //true if follower contained entry matching prevLogIndex and prevLogTerm
	FollowerID  string
}

// VoteReq Invoked by candidates to gather votes (§5.2).
type VoteReq struct {
	Term         uint64 //candidate’s term
	CandidateId  string //candidate requesting vote
	LastLogIndex uint64 //index of candidate’s last log entry (§5.4)
	LastLogTerm  uint64 //term of candidate’s last log entry (§5.4)
}

type VoteResp struct {
	Term        uint64 // currentTerm, for candidate to update itself
	VoteGranted bool   //true means candidate received vote
}

const (
	TypeAppendEntriesReq  uint32 = 111
	TypeAppendEntriesResp uint32 = 112
	TypeVoteReq           uint32 = 113
	TypeVoteResp          uint32 = 114
)

func marshal(obj interface{}) (*typelen.Packet, error) {
	packet := typelen.NewPacket()
	switch obj.(type) {
	case *AppendEntriesReq:
		packet.Type = TypeAppendEntriesReq
	case *AppendEntriesResp:
		packet.Type = TypeAppendEntriesResp
	case *VoteReq:
		packet.Type = TypeVoteReq
	case *VoteResp:
		packet.Type = TypeVoteResp
	}
	buffer, err := json.Marshal(obj)
	if err != nil {
		return packet, err
	}
	packet.Len = uint32(len(buffer))
	packet.Data = buffer
	return packet, nil
}
