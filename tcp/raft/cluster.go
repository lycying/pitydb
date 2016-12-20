package raft

import (
	"github.com/lycying/pitydb/tcp/mut"
	"strings"
	"sync"
)

type Cluster struct {
	collectionLock *sync.RWMutex
	cfg            *mut.Config
	endPoints      map[string]*mut.Client
	srv            *mut.Server
}

func NewCluster(cfg *mut.Config) *Cluster {
	cluster := &Cluster{}
	cluster.cfg = cfg
	cluster.collectionLock = new(sync.RWMutex)
	cluster.endPoints = make(map[string]*mut.Client)
	return cluster
}

func (cluster *Cluster) InitCluster(address string, peers []string) error {
	cluster.srv = mut.NewServer(address, cluster.cfg)
	err := cluster.srv.Servo()
	if err != nil {
		return err
	}
	for _, peer := range peers {
		cluster.JoinCluster(peer)
	}
	return nil
}

func (cluster *Cluster) JoinCluster(peer string) {
	cluster.collectionLock.Lock()
	defer cluster.collectionLock.Unlock()

	peer = strings.TrimSpace(peer)
	if _, ok := cluster.endPoints[peer]; ok {
		logger.Warn("raft# cluster repeat join %v , ignore", peer)
		return
	}
	client := mut.NewClient(peer, cluster.cfg)
	cluster.endPoints[peer] = client
	go client.DialAsync()
	logger.Info("raft# cluster join %v , welcome", peer)
}

func (cluster *Cluster) LeaveCluster(peer string) {
	cluster.collectionLock.Lock()
	defer cluster.collectionLock.Unlock()

	peer = strings.TrimSpace(peer)
	if val, ok := cluster.endPoints[peer]; ok {
		logger.Warn("raft# cluster %v dispose", val)
		val.Close()
		delete(cluster.endPoints, peer)
		return
	}
}

func (cluster *Cluster) Broadcast(packet *mut.Packet) {
	for _, value := range cluster.endPoints {
		if value.IsConnected() {
			value.Conn().WriteAsync(packet)
		}
	}
}
