package dt

import "sync/atomic"

func (p *AtomicBool) GetAndSet(v bool) bool {
	for {
		next := int32(0)
		if v {
			next = int32(1)
		}
		current := p.value
		if atomic.CompareAndSwapInt32(&p.value, current, next) {
			if current == 0x0 {
				return false
			} else {
				return true
			}
		}
	}
}

func (p *Int32) GetAndSet(v int32) int32 {
	for {
		current := p.value
		if atomic.CompareAndSwapInt32(&p.value, current, v) {
			return current
		}
	}
}

func (p *Int32) GetAndIncrement() int32 {
	for {
		current := p.value
		next := current + 1
		if atomic.CompareAndSwapInt32(&p.value, current, next) {
			return current
		}
	}
}

func (p *Int32) GetAndDecrement() int32 {
	for {
		current := p.value
		next := current - 1
		if atomic.CompareAndSwapInt32(&p.value, current, next) {
			return current
		}
	}
}
func (p *Int32) GetAndAdd(delta int32) int32 {
	for {
		current := p.value
		next := current + delta
		if atomic.CompareAndSwapInt32(&p.value, current, next) {
			return current
		}
	}
}

func (p *Int32) IncrementAndGet() int32 {
	for {
		current := p.value
		next := current + 1
		if atomic.CompareAndSwapInt32(&p.value, current, next) {
			return next
		}
	}
}
func (p *Int32) DecrementAndGet() int32 {
	for {
		current := p.value
		next := current - 1
		if atomic.CompareAndSwapInt32(&p.value, current, next) {
			return next
		}
	}
}
func (p *Int32) AddAndGet(delta int32) int32 {
	for {
		current := p.value
		next := current + delta
		if atomic.CompareAndSwapInt32(&p.value, current, next) {
			return next
		}
	}
}
func (p *UInt32) GetAndSet(v uint32) uint32 {
	for {
		current := p.value
		if atomic.CompareAndSwapUint32(&p.value, current, v) {
			return current
		}
	}
}

func (p *UInt32) GetAndIncrement() uint32 {
	for {
		current := p.value
		next := current + 1
		if atomic.CompareAndSwapUint32(&p.value, current, next) {
			return current
		}
	}
}

func (p *UInt32) GetAndDecrement() uint32 {
	for {
		current := p.value
		next := current - 1
		if atomic.CompareAndSwapUint32(&p.value, current, next) {
			return current
		}
	}
}
func (p *UInt32) GetAndAdd(delta uint32) uint32 {
	for {
		current := p.value
		next := current + delta
		if atomic.CompareAndSwapUint32(&p.value, current, next) {
			return current
		}
	}
}

func (p *UInt32) IncrementAndGet() uint32 {
	for {
		current := p.value
		next := current + 1
		if atomic.CompareAndSwapUint32(&p.value, current, next) {
			return next
		}
	}
}
func (p *UInt32) DecrementAndGet() uint32 {
	for {
		current := p.value
		next := current - 1
		if atomic.CompareAndSwapUint32(&p.value, current, next) {
			return next
		}
	}
}
func (p *UInt32) AddAndGet(delta uint32) uint32 {
	for {
		current := p.value
		next := current + delta
		if atomic.CompareAndSwapUint32(&p.value, current, next) {
			return next
		}
	}
}
func (p *UInt32) SubAndGet(delta uint32) uint32 {
	for {
		current := p.value
		next := current - delta
		if atomic.CompareAndSwapUint32(&p.value, current, next) {
			return next
		}
	}
}

func (p *Int64) GetAndSet(v int64) int64 {
	for {
		current := p.value
		if atomic.CompareAndSwapInt64(&p.value, current, v) {
			return current
		}
	}
}

func (p *Int64) GetAndIncrement() int64 {
	for {
		current := p.value
		next := current + 1
		if atomic.CompareAndSwapInt64(&p.value, current, next) {
			return current
		}
	}
}

func (p *Int64) GetAndDecrement() int64 {
	for {
		current := p.value
		next := current - 1
		if atomic.CompareAndSwapInt64(&p.value, current, next) {
			return current
		}
	}
}
func (p *Int64) GetAndAdd(delta int64) int64 {
	for {
		current := p.value
		next := current + delta
		if atomic.CompareAndSwapInt64(&p.value, current, next) {
			return current
		}
	}
}

func (p *Int64) IncrementAndGet() int64 {
	for {
		current := p.value
		next := current + 1
		if atomic.CompareAndSwapInt64(&p.value, current, next) {
			return next
		}
	}
}
func (p *Int64) DecrementAndGet() int64 {
	for {
		current := p.value
		next := current - 1
		if atomic.CompareAndSwapInt64(&p.value, current, next) {
			return next
		}
	}
}
func (p *Int64) AddAndGet(delta int64) int64 {
	for {
		current := p.value
		next := current + delta
		if atomic.CompareAndSwapInt64(&p.value, current, next) {
			return next
		}
	}
}

func (p *UInt64) GetAndSet(v uint64) uint64 {
	for {
		current := p.value
		if atomic.CompareAndSwapUint64(&p.value, current, v) {
			return current
		}
	}
}

func (p *UInt64) GetAndIncrement() uint64 {
	for {
		current := p.value
		next := current + 1
		if atomic.CompareAndSwapUint64(&p.value, current, next) {
			return current
		}
	}
}

func (p *UInt64) GetAndDecrement() uint64 {
	for {
		current := p.value
		next := current - 1
		if atomic.CompareAndSwapUint64(&p.value, current, next) {
			return current
		}
	}
}
func (p *UInt64) GetAndAdd(delta uint64) uint64 {
	for {
		current := p.value
		next := current + delta
		if atomic.CompareAndSwapUint64(&p.value, current, next) {
			return current
		}
	}
}

func (p *UInt64) IncrementAndGet() uint64 {
	for {
		current := p.value
		next := current + 1
		if atomic.CompareAndSwapUint64(&p.value, current, next) {
			return next
		}
	}
}
func (p *UInt64) DecrementAndGet() uint64 {
	for {
		current := p.value
		next := current - 1
		if atomic.CompareAndSwapUint64(&p.value, current, next) {
			return next
		}
	}
}
func (p *UInt64) AddAndGet(delta uint64) uint64 {
	for {
		current := p.value
		next := current + delta
		if atomic.CompareAndSwapUint64(&p.value, current, next) {
			return next
		}
	}
}
func (p *UInt64) SubAndGet(delta uint64) uint64 {
	for {
		current := p.value
		next := current - delta
		if atomic.CompareAndSwapUint64(&p.value, current, next) {
			return next
		}
	}
}
