package raft

import (
	"github.com/lycying/pitydb/tcp/mut"
	"math"
	"time"
)

func (rf *Raft) doAppendEntriesReq() {
	logger.Debug("raft# %v broadcast append entries ", rf.serverStr)

	for k, v := range rf.clients {
		if k != rf.serverStr && v.IsConnected() {
			req := &AppendEntriesReq{}
			req.Term = rf.currentTerm
			req.LeaderCommit = rf.raftLog.lastID()
			req.LeaderId = rf.serverStr
			req.PrevLogIndex = rf.nextIndex[k]

			last := rf.raftLog.getEntry(req.PrevLogIndex)

			req.PrevLogTerm = last.Term
			req.Entries = rf.raftLog.getRange(req.PrevLogIndex, rf.raftLog.commitIndex)

			logger.Debug("b from %v to %v %+v", rf.nextIndex[k], rf.raftLog.commitIndex, req.Entries)
			p, _ := marshal(req)
			v.Conn().WriteAsync(p)
		}
	}

	rf.heartbeatTimer.Reset(time.Second * 3)
}

func (rf *Raft) onAppendEntriesReq(c *mut.Conn, req *AppendEntriesReq) {
	logger.Debug("raft# %v rev append req %+v", rf.serverStr, req)
	resp := &AppendEntriesResp{}
	resp.Term = rf.currentTerm
	resp.FollowerID = rf.serverStr
	resp.Success = true

	if req.Term < rf.currentTerm {
		logger.Warn("raft# %v append req error req.Term=%+v,currentTerm=%+v", rf.serverStr, req.Term, rf.currentTerm)
		resp.Success = false
		p, _ := marshal(resp)
		c.WriteAsync(p)
		return
	}

	if req.Term == rf.currentTerm {
		if rf.state == stateCandidate {
			rf.state = stateFollower
		}
	} else { //req.Term > rf.currentTerm
		rf.currentTerm = req.Term
	}

	item := rf.raftLog.getEntry(req.PrevLogIndex)
	if item.Term != req.PrevLogTerm {
		rf.raftLog.discard(req.PrevLogIndex)
		resp.Success = false

		p, _ := marshal(resp)
		c.WriteAsync(p)
		return
	}

	if len(req.Entries) > 0 {
		rf.raftLog.inserts(req.Entries)

		if req.LeaderCommit > rf.raftLog.commitIndex {
			rf.raftLog.commitIndex = uint64(math.Max(float64(req.LeaderCommit), float64(req.Entries[0].ID)))
		}
	}

	resp.CommitIndex = rf.raftLog.commitIndex

	p, _ := marshal(resp)
	c.WriteAsync(p)
}
func (rf *Raft) onAppendEntriesResp(c *mut.Conn, resp *AppendEntriesResp) {
	logger.Debug("raft# %v rev append resp %+v", rf.serverStr, resp)
	rf.nextIndex[resp.FollowerID] = resp.CommitIndex
}
