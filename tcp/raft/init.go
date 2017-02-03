package raft

import "github.com/lycying/log"

var logger *log.Logger
var SEQ *Seq

func init() {
	logger, _ = log.New(log.DEBUG, "")

	SEQ = NewSeq()
}
