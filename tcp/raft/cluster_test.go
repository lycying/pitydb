package raft

import (
	"github.com/lycying/pitydb/tcp/mut"
	"github.com/lycying/pitydb/tcp/mut/codec/typelen"
	"testing"
	"time"
)

func TestCluster_InitCluster(t *testing.T) {
	cfg := mut.DefaultConfig()
	cfg.SetCallback(mut.NewHandlerSkeleton())
	cfg.SetCodec(typelen.NewCodec())
	cl := NewCluster(cfg)
	err := cl.InitCluster(":3122", []string{"127.0.0.1:3121", "127.0.0.1:3122", "127.0.0.1:3123"})
	if err != nil {
		logger.Err(err, "")
	}
	cl.JoinCluster("127.0.0.1:3121 ")
	time.Sleep(time.Second * 3)
	cl.LeaveCluster(" 127.0.0.1:3121")
	cl.LeaveCluster(" 127.0.0.1:3123")

	time.Sleep(time.Hour)
}
