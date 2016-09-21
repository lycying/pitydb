package autoinc

import (
	"sync/atomic"
	"runtime"
)

type AutoIncrementId32 struct {
	ops uint32
}

func (this *AutoIncrementId32) GetNext() uint32 {
	atomic.AddUint32(&this.ops, 1)
	runtime.Gosched()
	return atomic.LoadUint32(&this.ops)
}
func (this *AutoIncrementId32) GetCurrent() uint32 {
	return atomic.LoadUint32(&this.ops)
}
func NewAutoIncrementId32(init uint32) *AutoIncrementId32 {
	return &AutoIncrementId32{
		ops:init,
	}
}