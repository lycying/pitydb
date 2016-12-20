package raft

type Raft struct {
	cluster  *Cluster
	logStore LogStore
}

func NewRaft() *Raft {
	rf := &Raft{}
	return rf
}
