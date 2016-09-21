package slot

type StringTuple struct {
	val string
}

func (this *StringTuple) ToBytes() []byte {
	return []byte(this.val)
}