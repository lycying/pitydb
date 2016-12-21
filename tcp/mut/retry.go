package mut

import (
	"time"
)

const defaultMaxRetryTime = time.Second * 60

type retry struct {
	maxRetryTime time.Duration
	tempDelay    time.Duration

	retryTime int
}

func newValidRetry(maxRetryTime time.Duration) *retry {
	return &retry{
		maxRetryTime: maxRetryTime,
		tempDelay:    0,
		retryTime:    0,
	}
}
func newRetry() *retry {
	return newValidRetry(defaultMaxRetryTime)
}

func (r *retry) retryAfter(info interface{}) {
	if r.tempDelay == 0 {
		r.tempDelay = 60 * time.Millisecond
	} else {
		r.tempDelay *= 2
	}
	if max := r.maxRetryTime; r.tempDelay > max {
		r.tempDelay = max
	}

	logger.Warn("mut# %v retry after %+v ,has retried %d times", info, r.tempDelay, r.retryTime)

	time.Sleep(r.tempDelay)
	r.retryTime++
}

func (r *retry) reset() {
	r.tempDelay = 0
}
