package utils

import "testing"

func TestBitmap(t *testing.T) {
	var flag uint64 = 11
	bp := NewBitmap(100)
	bp.SetBit(112, true)
	bp.SetBit(440, false)
	bp.SetBit(999, true)
	bp.SetBit(998, false)
	bp.SetBit(flag, false)

	if !bp.GetBit(999) || bp.GetBit(998) || !bp.GetBit(112) || bp.GetBit(440) || bp.GetBit(flag) {
		t.Fail()
	}
}
