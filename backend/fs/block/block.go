package block

import "os"

const (
	EightKB int64 = 1024 * 8
)

type Block interface {
	Read() (int, error)
	Write(buf []byte) (int, error)
}
type blockhead struct {
	//为块所分配的唯一bid,可以根据其seek到具体的磁盘块,bid有0开始
	bid  int64
	//物理存储文件
	link *os.File
}

type blockentry struct {
	blockhead

	//是否已经删除，删除一个block并不降其真正删除，遍历每个block去真正删除这些区间无疑是效率缓慢的，需要一种机制来解决
	del bool
	//是否有脏数据，根据此状态会判断是否进行一些整理操作
	dirty  bool
	buf    []byte
}

func NewBlock(bid int64, link *os.File) Block {
	return &blockentry{
		blockhead:blockhead{
			bid:bid,
			link:link,
		},
		buf:make([]byte, EightKB),
		dirty:false,
		del:false,
	}
}
func (this *blockentry) Read() (int, error) {
	bid := this.blockhead.bid
	return this.blockhead.link.ReadAt(this.buf, bid * EightKB)
}
func (this *blockentry) Write(buf []byte) (int, error) {
	bid := this.blockhead.bid
	return this.blockhead.link.WriteAt(buf, bid * EightKB)
}