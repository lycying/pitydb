package mut

import (
	"sync"
	"sync/atomic"
)

type ConnMgr struct {
	connections map[uint64]*Conn
	lock        *sync.RWMutex
	nextID      uint64
}

func newConnectionMgr() *ConnMgr {
	mgr := new(ConnMgr)
	mgr.connections = make(map[uint64]*Conn)
	mgr.lock = new(sync.RWMutex)
	mgr.nextID = 0
	return mgr
}

func (mgr *ConnMgr) getNextId() uint64 {
	for {
		current := mgr.nextID
		next := current + 1
		if atomic.CompareAndSwapUint64(&mgr.nextID, current, next) {
			return next
		}
	}
}

func (mgr *ConnMgr) add(conn *Conn) {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	id := mgr.getNextId()
	logger.Debug("mut# mgr: id(%d) %v=>%v", id, conn.socket.RemoteAddr(), conn.socket.LocalAddr())
	mgr.connections[id] = conn
	conn.connID = id
}

func (mgr *ConnMgr) remove(conn *Conn) bool {
	mgr.lock.Lock()
	defer mgr.lock.Unlock()

	logger.Debug("mut# mgr: id(%d) %v=>%v", conn.connID, conn.socket.RemoteAddr(), conn.socket.LocalAddr())
	if _, ok := mgr.connections[conn.connID]; ok {
		delete(mgr.connections, conn.connID)
		return true
	}
	return false
}

func (mgr *ConnMgr) Iterate(callback func(uint64, *Conn) bool) {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()

	for key, conn := range mgr.connections {
		if !callback(key, conn) {
			break
		}
	}
}

func (mgr *ConnMgr) Count() int {
	mgr.lock.RLock()
	defer mgr.lock.RUnlock()

	return len(mgr.connections)
}
