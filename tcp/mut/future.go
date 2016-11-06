package mut

import (
	"context"
	"time"
)

type funcResult struct {
	value interface{}
	err   error
}
type Future struct {
	result chan *funcResult
}

func NewFuture(f func() (interface{}, error)) *Future {
	future := &Future{
		result: make(chan *funcResult),
	}

	go func() {
		defer close(future.result)
		value, err := f()
		future.result <- &funcResult{value: value, err: err}

	}()
	return future
}

func (future *Future) GetWithTimeout(timeout time.Duration) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-future.result:
		return result.value, result.err

	}
}
