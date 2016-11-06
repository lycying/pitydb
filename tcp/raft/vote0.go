package raft

import (
	"github.com/lycying/pitydb/tcp/mut"
	"math"
	"time"
)

func (rf *Raft) doVoteReq() {
	logger.Debug("raft# %v request votes, and vote itself", rf.serverStr)
	rf.state = stateCandidate

	//vote itself
	rf.votedFor = rf.serverStr
	rf.votes.SetValue(uint32(0))
	rf.votes.GetAndIncrement()

	//loop send to other servers
	req := &VoteReq{}
	req.Term = rf.currentTerm
	req.LastLogIndex = rf.raftLog.lastID()
	req.LastLogTerm = rf.currentTerm
	req.CandidateId = rf.serverStr
	p, _ := marshal(req)

	for k, v := range rf.clients {
		if k != rf.serverStr && v.IsConnected() {
			v.Conn().WriteAsync(p)
		}
	}
}
func (rf *Raft) onVoteReq(c *mut.Conn, req *VoteReq) {
	resp := &VoteResp{}
	resp.Term = rf.currentTerm
	resp.VoteGranted = false

	//compare the lastlog index and the last logterm
	if (rf.currentTerm == req.LastLogTerm && rf.raftLog.lastID() > req.LastLogIndex) ||
		rf.currentTerm > req.LastLogTerm {
	} else {
		if rf.state == stateFollower && rf.votedFor == "" {
			rf.votedFor = req.CandidateId
			resp.VoteGranted = true
			logger.Info("raft# %v vote for %+v", rf.serverStr, req.CandidateId)
		}
	}

	p, _ := marshal(resp)
	c.WriteAsync(p)
}

func (rf *Raft) onVoteResp(c *mut.Conn, resp *VoteResp) {
	if resp.VoteGranted {
		vt := rf.votes.IncrementAndGet()
		total := rf.onlineNumber()
		if int(vt)*2 >= total {
			if rf.state == stateLeader {
				logger.Info("raft# %+v alreay a leader | total:%v,got:%v", rf.serverStr, total, vt)
			} else {
				rf.state = stateLeader
				rf.heartbeatTimer.Reset(time.Millisecond * 20)
				rf.electionTimer.Reset(math.MaxInt64)

				for k, _ := range rf.clients {
					rf.nextIndex[k] = rf.raftLog.lastID()
					rf.matchIndex[k] = rf.raftLog.lastID()
				}
				logger.Info("raft# %+v become a leader | total:%v,got:%v", rf.serverStr, total, vt)
			}
		}
	}
}
