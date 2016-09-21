package bitmap

import "testing"

func TestBitmap_SetGet(t *testing.T) {
	var flag uint64 = 11
	bp := NewBitmap(100)
	bp.SetBit(112, true)
	t.Log(bp.Size(), bp.Size() == 200)
	bp.SetBit(440, false)
	t.Log(bp.Size(), bp.Size() == 440)
	bp.SetBit(999, true)
	t.Log(bp.Size())
	bp.SetBit(998, false)
	t.Log(bp.Size())
	bp.SetBit(flag, false)

	t.Log(bp.GetBit(999) == true)
	t.Log(bp.GetBit(998) == false)
	t.Log(bp.GetBit(112) == true)
	t.Log(bp.GetBit(440) == false)
	t.Log(bp.GetBit(flag) == false)
}
