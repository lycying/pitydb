package mut

import (
	"errors"
	"time"
)

type SyncFuture struct {
	ch chan interface{}
}

func NewSyncFuture() *SyncFuture {
	return &SyncFuture{
		ch: make(chan interface{}, 1),
	}
}

func (f *SyncFuture) WaitFor(timeout time.Duration) (interface{}, error) {
	timeoutCh := time.After(timeout)
	select {
	case obj := <-f.ch:
		return obj, nil
	case <-timeoutCh:
		return nil, errors.New("timeout")
	}
}

func (f *SyncFuture) SetValue(obj interface{}) {
	f.ch <- obj
}
