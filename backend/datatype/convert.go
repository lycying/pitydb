package datatype

import "encoding/binary"

type Serialization interface {
	Bytes(e interface{}) []byte
	Make(buf []byte) interface{}
}
//convert an int64 to []byte and no use to create a []byte outside
func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

//convert an int64 from []byte
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

//convert an int64 to []byte and no use to create a []byte outside
func Int32ToBytes(i int32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

//convert an int32 from []byte
func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

//convert a string to []byte
func StringToBytes(str string) []byte {
	return []byte(str)
}

//convert a []byte to string, utf8 only
func BytesToString(buf []byte) string {
	return string(buf)
}