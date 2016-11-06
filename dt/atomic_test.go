package dt

import (
	"sync"
	"testing"
)

func TestUInt64_GetAndAdd(t *testing.T) {
	var wg sync.WaitGroup
	i := NewUInt64()
	i.SetValue(uint64(0))

	loop := 100000

	for j := 0; j < loop; j++ {
		wg.Add(1)
		go func() {
			i.GetAndIncrement()
			wg.Add(-1)
		}()
	}

	wg.Wait()

	if int(i.GetValue().(uint64)) != loop {
		t.Fail()
	}

	for j := 0; j < loop; j++ {
		wg.Add(1)
		go func() {
			i.AddAndGet(10)
			wg.Add(-1)
		}()
	}
	wg.Wait()

	if int(i.GetValue().(uint64)) != 10*loop+loop {
		t.Fail()
	}
	for j := 0; j < loop; j++ {
		wg.Add(1)
		go func() {
			i.SubAndGet(10)
			wg.Add(-1)
		}()
	}
	wg.Wait()

	if int(i.GetValue().(uint64)) != loop {
		t.Fail()
	}
}
