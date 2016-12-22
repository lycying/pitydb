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
}
