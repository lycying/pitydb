package autoinc

import (
	"testing"
	"time"
)

func tesGetValConcurrent(from int, p *AutoIncrementId32, t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Log(from, p.GetCurrent(), p.GetNext())
	}
}
func TestAutoIncrementId32_GetNextVal(t *testing.T) {
	x := &AutoIncrementId32{
		ops:99,
	}

	go tesGetValConcurrent(1, x, t)
	go tesGetValConcurrent(2, x, t)

	time.Sleep(time.Second)

	t.Log("Now the value shold be 299 :", x.GetCurrent() == 299)

}
