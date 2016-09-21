package tuple

import (
	"../slot"
)

const S_INT byte = 0x01
const S_LONG byte = 0x02
const S_FLOAT byte = 0x03
const S_DOUBLE byte = 0x04
const S_BOOL byte = 0x05
//提供记录元组的结构描述
type TupleDesc struct {
	items []string
}

//最小的元组存储单元，元组不能跨数据快存储
type Tuple struct {
	pre   *Tuple
	next  *Tuple
	slots []*slot.Slot
}

type PersistTuple struct {
	key  int
	pre  int
	next int
}

func Test() {}