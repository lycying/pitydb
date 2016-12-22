package mut

import (
	"errors"
	"time"
)

type SyncFuture struct {
	ch        chan interface{}
	timerChan <-chan time.Time
}

func NewSyncFuture() *SyncFuture {
	return &SyncFuture{
		ch: make(chan interface{}, 1),
	}
}

func (f *SyncFuture) WaitFor(timeout time.Duration) (interface{}, error) {
	f.timerChan = time.After(timeout)
	select {
	case obj := <-f.ch:
		return obj, nil
	case b := <-f.timerChan:
		logger.Debug("%+v", b)
		return nil, errors.New("timeout")
	}
}

func (f *SyncFuture) SetValue(obj interface{}) {
	f.ch <- obj
}

func (f *SyncFuture) Cancel() {
	f.ch <- nil
}
