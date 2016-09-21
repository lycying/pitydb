package bitmap

type Bitmap struct {
	data []byte
	size uint64
}

func max(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}else {
		return b
	}
}
func (this *Bitmap) SetBit(offset uint64, value bool) {
	index, pos := offset / 8, offset % 8
	if this.size < offset {
		size := max(this.size << 1, offset)
		tmp := make([]byte, size)
		copy(tmp, this.data)
		this.data = tmp
		this.size = size
	}

	if value {
		this.data[index] |= 0x01 << pos
	}else {
		this.data[index] &^= 0x01 << pos
	}
}

func (this *Bitmap) GetBit(offset uint64) bool {
	index, pos := offset / 8, offset % 8

	if this.size < offset {
		return false
	}

	return 0 != (this.data[index] >> pos) & 0x01
}

func NewBitmap(size uint64) *Bitmap {
	return &Bitmap{
		data:make([]byte, size),
		size:size,
	}
}

func (this *Bitmap) Size() uint64 {
	return this.size
}