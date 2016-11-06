package dt

type DType int32

const (
	ByteType DType = iota
	Int32Type
	UInt32Type
	Int64Type
	UInt64Type
	Float32Type
	Float64Type
	BoolType //1 is true , 0 is false
	StringType
	DecimalType
	TimeType
	ArrayType
	JsonType
)

// Encoder is the interface representing objects that can encode themselves.
type Encoder interface {
	Encode() ([]byte, error)
}

// Decoder is the interface representing objects that can decode themselves
type DeCoder interface {
	Decode(buf []byte, offset int) (int, error)
}

// ValueRefer is the interface that refer to the real value
type ValueRefer interface {
}

type Comparator interface {
	Compare(Comparator) int
}

type DtRefer interface {
	Encoder
	DeCoder
	Comparator

	SetValue(ValueRefer)
	GetValue() ValueRefer
	GetLen() int
	Copy() DtRefer
}

func NewDtRefer(dType DType) DtRefer {
	return NewDtReferWithSub(dType, -1)
}
func NewDtReferWithSub(dType DType, subType DType) DtRefer {
	var r DtRefer = nil
	switch dType {
	case Int32Type:
		r = NewInt32()
	case UInt32Type:
		r = NewUInt32()
	case ByteType:
		r = NewByte()
	case BoolType:
		r = NewBool()
	case Int64Type:
		r = NewInt64()
	case UInt64Type:
		r = NewUInt64()
	case Float32Type:
		r = NewFloat32()
	case Float64Type:
		r = NewFloat64()
	case StringType:
		r = NewString()
	case ArrayType:
		r = NewArray(subType)
	case DecimalType:
	//TODO
	case TimeType:
	//TODO
	case JsonType:
		//TODO
	}
	return r
}
