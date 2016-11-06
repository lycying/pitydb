package mut

import (
	"context"
	"testing"
	"time"
)

func TestFuture_GetWithTimeout(t *testing.T) {
	var a string = "test"
	future := NewFuture(func() (interface{}, error) {
		time.Sleep(time.Second * 2)
		a = "fucking going "

		return nil, nil
	})

	_, err := future.GetWithTimeout(time.Second / 3)

	if err != context.DeadlineExceeded {
		t.Fatal()
	}

	time.Sleep(time.Second * 3)
	t.Log(a)

}
