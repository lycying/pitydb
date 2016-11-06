package dt

import "bytes"

// Array is array that has the same data type
type Array struct {
	DtRefer
	value []DtRefer
	size  *UInt32

	//it's a template flag that never be saved
	mType DType
}

func NewArray(aType DType) *Array {
	return &Array{
		mType: aType,
		size:  ValidNewUInt32(0),
		value: make([]DtRefer, 0),
	}
}

func (p *Array) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	sizeArr, _ := p.size.Encode()
	buf.Write(sizeArr)
	if p.size.GetValue().(uint32) > 0 {
		for _, item := range p.value {
			b, _ := item.Encode()
			buf.Write(b)
		}
	}
	return buf.Bytes(), nil
}
func (p *Array) Decode(buf []byte, offset int) (int, error) {
	index := 0
	sizeArr, _ := p.size.Decode(buf, offset)
	index = index + sizeArr

	length := int(p.size.GetValue().(uint32))
	if length > 0 {
		for i := 0; i < length; i++ {
			dt := NewDtRefer(p.mType)
			dtLen, _ := dt.Decode(buf, offset+index)
			index = index + dtLen
		}
	}
	return index, nil
}

func (p *Array) SetValue(v ValueRefer) {
	p.value = v.([]DtRefer)
	p.size.SetValue(uint32(len(p.value)))
}

func (p *Array) GetValue() ValueRefer {
	return p.value
}
func (p *Array) GetLen() int {
	i := 0
	i += p.size.GetLen()
	for _, item := range p.value {
		i += item.GetLen()
	}
	return i
}
func (p *Array) Copy() DtRefer {
	ret := NewArray(p.mType)
	ret.size.SetValue(p.size.GetValue())
	length := int(p.size.GetValue().(uint32))
	if length > 0 {
		ret.value = make([]DtRefer, length)
		for i := 0; i < length; i++ {
			ret.value[i] = p.value[i]
		}
	}
	return ret
}
func (p *Array) DeepCopy() DtRefer {
	ret := NewArray(p.mType)
	ret.size.SetValue(p.size.GetValue())
	length := int(p.size.GetValue().(uint32))
	if length > 0 {
		ret.value = make([]DtRefer, length)
		for i := 0; i < length; i++ {
			ret.value[i] = p.value[i].Copy()
		}
	}
	return ret
}

func (p *Array) GetSize() int {
	return int(p.size.GetValue().(uint32))
}

func (p *Array) Add(v DtRefer, canEquals bool) bool {
	if canEquals {
		p.value = append(p.value, v)
		return true
	} else {
		insert := true
		for _, item := range p.value {
			if item.Compare(v) == 0 {
				insert = false
				break
			}
		}
		if insert {
			p.value = append(p.value, v)
			return true
		}
		return false
	}
}
