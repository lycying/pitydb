package raft

import "github.com/lycying/pitydb/dt"

type Seq struct {
	seq *dt.UInt64
}

func NewSeq() *Seq {
	return &Seq{
		seq: dt.NewUInt64(),
	}
}

func (seq *Seq) NextID() uint64 {
	return seq.seq.IncrementAndGet()
}
