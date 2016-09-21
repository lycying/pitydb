package membuf

import "os"

//在数据刷新到磁盘之前，需要先在内存中进行缓冲
//这部分数据可以是查询出来的缓冲数据，也可以是更新后修改的簇信息
//每个缓冲区可能对应多个磁盘文件，缓冲需要知道所对应的磁盘文件是否已经修改
//如果没有修改，则不需要更新磁盘文件，直接丢弃即可
//否则，就需要写入更新磁盘文件的某个区域或者索引信息
//根据策略配置，缓冲区可直接作为内存数据库，不必再刷磁盘
type Mem struct {
	link   os.File
	offset int
	len    int
}

type MemBuf interface {
	AppendBuf(buf []byte)
	RemoveBuf(buf []byte)
	Flush()
	MakeIndex()
	Search()
	Serialization() []byte
}

func BufferAlloc() {
}
