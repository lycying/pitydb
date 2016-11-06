package dt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"strings"
)

type Byte struct {
	DtRefer
	value byte
}
type Bool struct {
	DtRefer
	value bool
}
type AtomicBool struct {
	DtRefer
	value int32
}
type Int32 struct {
	DtRefer
	value int32
}

type Int64 struct {
	DtRefer
	value int64
}

type Float32 struct {
	DtRefer
	value float32
}

type Float64 struct {
	DtRefer
	value float64
}

type UInt32 struct {
	DtRefer
	value uint32
}

type UInt64 struct {
	DtRefer
	value uint64
}

type String struct {
	DtRefer
	value string
}

func NewInt32() *Int32 {
	return ValidNewInt32(0)
}
func NewUInt32() *UInt32 {
	return ValidNewUInt32(0)
}
func NewInt64() *Int64 {
	return ValidInt64(0)
}
func NewUInt64() *UInt64 {
	return ValidNewUInt64(0)
}
func NewFloat32() *Float32 {
	return ValidNewFloat32(0)
}
func NewFloat64() *Float64 {
	return ValidNewFloat64(0)
}
func NewBool() *Bool {
	return ValidNewBool(false)
}
func NewAtomicBool() *AtomicBool {
	return ValidNewAtomicBool(false)
}
func NewByte() *Byte {
	return ValidNewByte(0x0)
}
func NewString() *String {
	return ValidNewString("")
}

func ValidNewInt32(v int32) *Int32 {
	return &Int32{value: v}
}
func ValidNewUInt32(v uint32) *UInt32 {
	return &UInt32{value: v}
}
func ValidInt64(v int64) *Int64 {
	return &Int64{value: v}
}
func ValidNewUInt64(v uint64) *UInt64 {
	return &UInt64{value: v}
}
func ValidNewFloat32(v float32) *Float32 {
	return &Float32{value: v}
}
func ValidNewFloat64(v float64) *Float64 {
	return &Float64{value: v}
}
func ValidNewBool(v bool) *Bool {
	return &Bool{value: v}
}
func ValidNewAtomicBool(v bool) *AtomicBool {
	r := &AtomicBool{}
	r.SetValue(v)
	return r
}
func ValidNewByte(v byte) *Byte {
	return &Byte{value: v}
}
func ValidNewString(v string) *String {
	return &String{value: v}
}

func (p *Bool) Encode() ([]byte, error) {
	var b byte = 0x0
	if p.value {
		b = 0x1
	}
	return []byte{b}, nil
}

func (p *Bool) Decode(buf []byte, offset int) (int, error) {
	b := buf[offset]
	p.value = (b == 0x1)
	return 1, nil
}

func (p *Bool) SetValue(v ValueRefer) {
	p.value = v.(bool)
}

func (p *Bool) GetValue() ValueRefer {
	return p.value
}

func (p *Bool) GetLen() int {
	return 1
}
func (p *Bool) Copy() DtRefer {
	ret := NewBool()
	ret.value = p.value
	return ret
}
func (p *Bool) Compare(v Comparator) int {
	if p.value == v.(*Bool).value {
		return 0
	}
	return 1
}

func (p *Byte) Encode() ([]byte, error) {
	return []byte{p.value}, nil
}

func (p *Byte) Decode(buf []byte, offset int) (int, error) {
	b := buf[offset]
	p.value = b
	return 1, nil
}

func (p *Byte) SetValue(v ValueRefer) {
	p.value = v.(byte)
}

func (p *Byte) GetValue() ValueRefer {
	return p.value
}

func (p *Byte) GetLen() int {
	return 1
}

func (p *Byte) Copy() DtRefer {
	ret := NewByte()
	ret.value = p.value
	return ret
}

func (p *Byte) Compare(v Comparator) int {
	ov := v.(*Byte).value
	switch {
	case p.value == ov:
		return 0
	case p.value < ov:
		return -1
	default:
		return 1

	}
}
func (p *AtomicBool) Encode() ([]byte, error) {
	var b byte = 0x0
	if p.value == 1 {
		b = 0x1
	}
	return []byte{b}, nil
}

func (p *AtomicBool) Decode(buf []byte, offset int) (int, error) {
	b := buf[offset]
	if b == 0x1 {
		p.value = 1
	} else {
		p.value = 0
	}
	return 1, nil
}

func (p *AtomicBool) SetValue(v ValueRefer) {
	b := v.(bool)
	if b {
		p.value = 1
	} else {
		p.value = 0
	}
}

func (p *AtomicBool) GetValue() ValueRefer {
	if p.value == 0x1 {
		return true
	} else {
		return false
	}
}

func (p *AtomicBool) GetLen() int {
	return 1
}
func (p *AtomicBool) Copy() DtRefer {
	ret := &AtomicBool{}
	ret.value = p.value
	return ret
}
func (p *AtomicBool) Compare(v Comparator) int {
	ov := v.(*AtomicBool).value

	if ov == p.value {
		return 0
	}
	return 1
}

func (p *Int32) Encode() ([]byte, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(p.value))
	return buf, nil
}

func (p *Int32) Decode(buf []byte, offset int) (int, error) {
	p.value = int32(binary.BigEndian.Uint32(buf[offset : offset+4]))
	return 4, nil
}

func (p *Int32) SetValue(v ValueRefer) {
	p.value = v.(int32)
}

func (p *Int32) GetValue() ValueRefer {
	return p.value
}

func (p *Int32) GetLen() int {
	return 4
}
func (p *Int32) Copy() DtRefer {
	ret := NewInt32()
	ret.value = p.value
	return ret
}
func (p *Int32) Compare(v Comparator) int {
	ov := v.(*Int32).value
	switch {
	case p.value == ov:
		return 0
	case p.value < ov:
		return -1
	default:
		return 1

	}
}

func (p *UInt32) Encode() ([]byte, error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, p.value)
	return buf, nil
}

func (p *UInt32) Decode(buf []byte, offset int) (int, error) {
	p.value = binary.BigEndian.Uint32(buf[offset : offset+4])
	return 4, nil
}

func (p *UInt32) SetValue(v ValueRefer) {
	p.value = v.(uint32)
}

func (p *UInt32) GetValue() ValueRefer {
	return p.value
}

func (p *UInt32) GetLen() int {
	return 4
}

func (p *UInt32) Copy() DtRefer {
	ret := NewUInt32()
	ret.value = p.value
	return ret
}
func (p *UInt32) Compare(v Comparator) int {
	ov := v.(*UInt32).value
	switch {
	case p.value == ov:
		return 0
	case p.value < ov:
		return -1
	default:
		return 1
	}
}

func (p *Int64) Encode() ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(p.value))
	return buf, nil
}

func (p *Int64) Decode(buf []byte, offset int) (int, error) {
	p.value = int64(binary.BigEndian.Uint64(buf[offset : offset+8]))
	return 8, nil
}

func (p *Int64) SetValue(v ValueRefer) {
	p.value = v.(int64)
}

func (p *Int64) GetValue() ValueRefer {
	return p.value
}

func (p *Int64) GetLen() int {
	return 8
}
func (p *Int64) Copy() DtRefer {
	ret := NewInt64()
	ret.value = p.value
	return ret
}

func (p *Int64) Compare(v Comparator) int {
	ov := v.(*Int64).value
	switch {
	case p.value == ov:
		return 0
	case p.value < ov:
		return -1
	default:
		return 1
	}
}

func (p *UInt64) Encode() ([]byte, error) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, p.value)
	return buf, nil
}

func (p *UInt64) Decode(buf []byte, offset int) (int, error) {
	p.value = binary.BigEndian.Uint64(buf[offset : offset+8])
	return 8, nil
}

func (p *UInt64) SetValue(v ValueRefer) {
	p.value = v.(uint64)
}

func (p *UInt64) GetValue() ValueRefer {
	return p.value
}

func (p *UInt64) GetLen() int {
	return 8
}
func (p *UInt64) Copy() DtRefer {
	ret := NewUInt64()
	ret.value = p.value
	return ret
}

func (p *UInt64) Compare(v Comparator) int {
	ov := v.(*UInt64).value
	switch {
	case p.value == ov:
		return 0
	case p.value < ov:
		return -1
	default:
		return 1
	}
}

func (p *Float32) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p.value)

	return buf.Bytes(), nil
}

func (p *Float32) Decode(buf []byte, offset int) (int, error) {
	byteBuf := bytes.NewReader(buf[offset : offset+4])
	binary.Read(byteBuf, binary.BigEndian, &p.value)
	return 4, nil
}

func (p *Float32) SetValue(v ValueRefer) {
	p.value = v.(float32)
}

func (p *Float32) GetValue() ValueRefer {
	return p.value
}

func (p *Float32) GetLen() int {
	return 4
}

func (p *Float32) Copy() DtRefer {
	ret := NewFloat32()
	ret.value = p.value
	return ret
}

func (p *Float32) Compare(v Comparator) int {
	ov := v.(*Float32).value
	switch {
	case p.value == ov:
		return 0
	case p.value < ov:
		return -1
	default:
		return 1
	}
}

func (p *Float64) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p.value)

	return buf.Bytes(), nil
}

func (p *Float64) Decode(buf []byte, offset int) (int, error) {
	byteBuf := bytes.NewReader(buf[offset : offset+8])
	binary.Read(byteBuf, binary.BigEndian, &p.value)
	return 8, nil
}

func (p *Float64) SetValue(v ValueRefer) {
	p.value = v.(float64)
}

func (p *Float64) GetValue() ValueRefer {
	return p.value
}

func (p *Float64) GetLen() int {
	return 8
}

func (p *Float64) Copy() DtRefer {
	ret := NewFloat64()
	ret.value = p.value
	return ret
}

func (p *Float64) Compare(v Comparator) int {
	ov := v.(*Float64).value
	switch {
	case p.value == ov:
		return 0
	case p.value < ov:
		return -1
	default:
		return 1
	}
}

func (p *String) Encode() ([]byte, error) {
	var retArr []byte
	var lenFlag byte
	var offset int

	strArr := []byte(p.value)
	strLen := len(strArr)

	switch {
	case strLen > math.MaxUint32:
		return nil, errors.New("string too long")
	case strLen > math.MaxUint16:
		retArr = make([]byte, 1+4+strLen)
		binary.BigEndian.PutUint32(retArr[1:5], uint32(strLen))
		lenFlag = 0x3
		offset = 5
	case strLen > math.MaxUint8:
		retArr = make([]byte, 1+2+strLen)
		binary.BigEndian.PutUint16(retArr[1:3], uint16(strLen))
		lenFlag = 0x2
		offset = 3
	case strLen > 0:
		retArr = make([]byte, 1+1+strLen)
		retArr[1] = byte(strLen)
		lenFlag = 0x1
		offset = 2
	case strLen == 0:
		retArr = make([]byte, 1)
		lenFlag = 0x0
		offset = 1
	}
	retArr[0] = lenFlag
	if strLen > 0 {
		copy(retArr[offset:], strArr)
	}
	return retArr, nil
}

func (p *String) Decode(buf []byte, offset int) (int, error) {
	lenFlag := buf[offset]
	retLen := 1

	offset = offset + 1 // the retLen has 1 byte
	switch lenFlag {
	case 0x0:
		p.value = ""
	case 0x1:
		size := int(buf[offset])
		p.value = string(buf[offset+1 : offset+1+size])
		retLen = retLen + 1 + size
	case 0x2:
		size := int(binary.BigEndian.Uint16(buf[offset : offset+2]))
		p.value = string(buf[offset+2 : offset+2+size])
		retLen = retLen + 2 + size
	case 0x3:
		size := int(binary.BigEndian.Uint32(buf[offset : offset+4]))
		p.value = string(buf[offset+4 : offset+4+size])
		retLen = retLen + 4 + size
	}

	return retLen, nil
}

func (p *String) GetLen() int {
	strLen := len(p.value)
	switch {
	case strLen > math.MaxUint16:
		return 1 + 4 + strLen
	case strLen > math.MaxUint8:
		return 1 + 2 + strLen
	case strLen > 0:
		return 1 + 1 + strLen
	case strLen == 0:
		return 1
	}
	return 0 // never reached
}
func (p *String) SetValue(v ValueRefer) {
	p.value = v.(string)
}

func (p *String) GetValue() ValueRefer {
	return p.value
}
func (p *String) Copy() DtRefer {
	ret := NewString()
	ret.value = p.value
	return ret
}

func (p *String) Compare(v Comparator) int {
	if strings.EqualFold(p.value, v.(*String).value) {
		return 0
	}
	return 1
}
