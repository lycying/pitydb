package autoinc

import (
	"sync/atomic"
	"runtime"
)

type AutoIncrementId64 struct {
	ops uint64
}

func (this *AutoIncrementId64) GetNext() uint64 {
	atomic.AddUint64(&this.ops, 1)
	runtime.Gosched()
	return atomic.LoadUint64(&this.ops)
}
func (this *AutoIncrementId64) GetCurrent() uint64 {
	return atomic.LoadUint64(&this.ops)
}
func NewAutoIncrementId64(init uint64) *AutoIncrementId64 {
	return &AutoIncrementId64{
		ops:init,
	}
}