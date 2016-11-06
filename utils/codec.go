package utils

import (
	"hash/crc32"
)

func Sum32(bs []byte) uint32 {
	ieee := crc32.NewIEEE()
	ieee.Write(bs)
	return ieee.Sum32()
}
