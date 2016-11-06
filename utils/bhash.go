package utils

//https://github.com/julycoding/The-Art-Of-Programming-By-July/blob/master/ebook/zh/%E5%80%92%E6%8E%92%E7%B4%A2%E5%BC%95%E5%85%B3%E9%94%AE%E8%AF%8D%E4%B8%8D%E9%87%8D%E5%A4%8DHash%E7%BC%96%E7%A0%81.md
func init() {
	cryptTable = make([]uint64, 1280)
	prepareCryptTable()
}

var cryptTable []uint64

func prepareCryptTable() {
	var seed, index1, index2 uint64 = 0x00100001, 0, 0
	i := 0
	for index1 = 0; index1 < 0x100; index1 += 1 {
		for index2, i = index1, 0; i < 5; index2 += 0x100 {
			seed = (seed*125 + 3) % 0x2aaaab
			temp1 := (seed & 0xffff) << 0x10
			seed = (seed*125 + 3) % 0x2aaaab
			temp2 := seed & 0xffff
			cryptTable[index2] = temp1 | temp2
			i += 1
		}
	}
}

func HashString(str string, dwHashType int) uint64 {
	i, ch := 0, 0
	var seed1, seed2 uint64 = 0x7FED7FED, 0xEEEEEEEE
	var key uint8
	strLen := len(str)
	for i < strLen {
		key = str[i]
		ch = int(key)
		i += 1
		seed1 = cryptTable[(dwHashType<<8)+ch] ^ (seed1 + seed2)
		seed2 = uint64(ch) + seed1 + seed2 + (seed2 << 5) + 3
	}
	return uint64(seed1)
}
