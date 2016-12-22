package mut

import (
	"testing"
	"time"
)

func TestSyncFuture(t *testing.T) {
	f := NewSyncFuture()
	f.SetValue("1234")
	obj, err := f.WaitFor(time.Second)
	if err != nil {
		logger.Err(err, "err")
	} else {
		logger.Info("%v", obj)
	}

	///////////////////////////////////

	f = NewSyncFuture()
	go func() {
		time.Sleep(time.Second / 2)
		f.Cancel()
	}()
	obj, err = f.WaitFor(time.Second * 20)
	if err != nil {
		logger.Err(err, "err")
	} else {
		logger.Info("%v", obj)
	}

	///////////////////////////////////
	f = NewSyncFuture()
	obj, err = f.WaitFor(time.Second * 1)
	if err != nil {
		logger.Err(err, "err")
	} else {
		logger.Info("%v", obj)
	}
}
