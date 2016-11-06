package utils

type Bitmap struct {
	data []byte
	size uint64
}

func NewBitmap(size uint64) *Bitmap {
	return &Bitmap{
		data: make([]byte, size),
		size: size,
	}
}

func max(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func (bp *Bitmap) SetBit(offset uint64, value bool) {
	index, pos := offset/8, offset%8
	if bp.size < offset {
		size := max(bp.size<<1, offset)
		tmp := make([]byte, size)
		copy(tmp, bp.data)
		bp.data = tmp
		bp.size = size
	}

	if value {
		bp.data[index] |= 0x01 << pos
	} else {
		bp.data[index] &^= 0x01 << pos
	}
}

func (bp *Bitmap) GetBit(offset uint64) bool {
	index, pos := offset/8, offset%8

	if bp.size < offset {
		return false
	}

	return 0 != (bp.data[index]>>pos)&0x01
}

func (bp *Bitmap) Size() uint64 {
	return bp.size
}
