package yard

import (
	"bytes"
	"github.com/lycying/pitydb/dt"
	"sort"
)

const DefaultPageSize = 1024 * 2

const (
	indexPageType byte = iota
	dataPageType
)

type pageHeader struct {
	dt.Encoder
	dt.DeCoder

	pgID       *dt.UInt32 //4294967295.0*16/1024/1024/1024 ~= 63.99999998509884 TiB
	pgType     *dt.Byte
	level      *dt.UInt32
	left       *dt.UInt32
	right      *dt.UInt32
	checksum   *dt.UInt32
	lastModify *dt.UInt64 //time.Now().UnixNano()
	size       *dt.UInt32 //this counter is used to read data from disk
}

func (header *pageHeader) Encode() ([]byte, error) {

	bPgID, _ := header.pgID.Encode()
	bPgType, _ := header.pgType.Encode()
	bLevel, _ := header.level.Encode()
	bLeft, _ := header.left.Encode()
	bRight, _ := header.right.Encode()
	bChecksum, _ := header.checksum.Encode()
	bLastModify, _ := header.lastModify.Encode()
	bSize, _ := header.size.Encode()

	buf := new(bytes.Buffer)
	buf.Write(bPgID)
	buf.Write(bPgType)
	buf.Write(bLevel)
	buf.Write(bLeft)
	buf.Write(bRight)
	buf.Write(bChecksum)
	buf.Write(bLastModify)
	buf.Write(bSize)

	return buf.Bytes(), nil
}

func (header *pageHeader) Decode(buf []byte, offset int) (int, error) {

	idx := 0
	lenPgID, _ := header.pgID.Decode(buf, idx+offset)
	idx += lenPgID
	lenPgType, _ := header.pgType.Decode(buf, idx+offset)
	idx += lenPgType
	lenLevel, _ := header.level.Decode(buf, idx+offset)
	idx += lenLevel
	lenLeft, _ := header.left.Decode(buf, idx+offset)
	idx += lenLeft
	lenRight, _ := header.right.Decode(buf, idx+offset)
	idx += lenRight
	lenChecksum, _ := header.checksum.Decode(buf, idx+offset)
	idx += lenChecksum
	lenLastModify, _ := header.lastModify.Decode(buf, idx+offset)
	idx += lenLastModify
	lenSize, _ := header.size.Decode(buf, idx+offset)
	idx += lenSize

	return idx, nil
}

type Page struct {
	pageHeader

	left   *Page
	right  *Page
	parent *Page

	tree *PageTree

	rows []*Row //the tuple data
	_len int    //finger if the size is larger than 16kb
}

func (p *Page) Decode(buf []byte, offset int) (int, error) {
	idx, _ := p.pageHeader.Decode(buf, offset)
	for _, row := range p.rows {
		rLen, _ := row.Decode(buf, idx+offset)
		idx += rLen
	}
	return idx, nil
}

func (p *Page) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	bHeader, _ := p.pageHeader.Encode()
	buf.Write(bHeader)
	for _, row := range p.rows {
		bRow, _ := row.Encode()
		buf.Write(bRow)
	}
	return buf.Bytes(), nil
}

func (p *Page) findIndexRow(key uint32) (*Page, int, bool) {
	pSize := int(p.size.GetValue().(uint32))

	i := sort.Search(pSize, func(i int) bool {
		return key <= p.rows[i].GetKey()
	})
	//the rows is empty
	if i == 0 && pSize == 0 {
		return p, 0, false
	}
	if i >= pSize {
		i = i - 1
	}
	if i > 0 {
		x0 := p.rows[i].GetKey()
		x1 := p.rows[i-1].GetKey()

		if key < x0 && key > x1 {
			i = i - 1
		}
	}
	return p, i + 1, true

}

func (p *Page) findOne(key uint32) (*Page, int, bool) {
	if p.isIndexPage() {
		_, count, _ := p.findIndexRow(key)

		count = count - 1
		pgID := p.rows[count].cells[0].GetValue().(uint32)
		next := p.tree.mgr.GetPage(pgID)

		return next.findOne(pgID)
	}

	pSize := int(p.size.GetValue().(uint32))

	i := sort.Search(pSize, func(i int) bool {
		return key <= p.rows[i].GetKey()
	})
	//the rows is empty
	if i == 0 && pSize == 0 {
		return p, 0, false
	}

	//should put at the tail of the row array
	if i >= pSize {
		return p, pSize, false
	}

	if p.rows[i].GetKey() == key {
		return p, i, true
	}

	return p, i, false
}

func (p *Page) insert(row *Row, index int, find bool) int {
	bs := p._len + row.GetLen()
	if find {
		bs = bs - p.rows[index].GetLen()
		p.rows[index] = row
	} else {
		p.rows = append(p.rows[:index], append([]*Row{row}, p.rows[index:]...)...)
		pSize := p.size.GetValue().(uint32)
		p.size.SetValue(pSize + 1)
	}
	p._len = bs

	if p.shouldSplit() {
		//should split here
		i := 0
		counter := 0

		pSize := p.size.GetValue().(uint32)
		for ; i < int(pSize); i++ {
			nextCounter := counter + p.rows[i].GetLen()
			if nextCounter > DefaultPageSize {
				bs = counter
				break
			}
			counter = nextCounter
		}

		newPage := p.tree.NewPage(p.level.GetValue().(uint32), p.pgType.GetValue().(byte))
		//copy [:i-1] to newNode
		newPage.copyRightPart(p, i-1)
		//only left [i-1:] part
		p.deleteRightPart(i - 1)

		if p.hasParent() {
			indexRow := newPage.NewIndexRow()
			_, toIndex, _ := p.parent.findIndexRow(indexRow.GetKey())
			p.parent.insert(indexRow, toIndex, false)
			newPage.parent = p.parent

		} else {
			newRoot := p.tree.NewIndexPage(p.level.GetValue().(uint32) + 1)

			indexRow0 := p.NewIndexRow()
			newRoot.insert(indexRow0, 0, false)
			p.parent = newRoot

			indexRow1 := newPage.NewIndexRow()
			newRoot.insert(indexRow1, 1, false)
			newPage.parent = newRoot

			p.tree.root = newRoot
		}

	}
	return bs
}

func (p *Page) delete(key uint32, index int) {
	p.rows = append(p.rows[:index], p.rows[index+1:]...)
	pSize := p.size.GetValue().(uint32)
	p.size.SetValue(pSize - 1)
	p._len = p.GetLen()
}

func (p *Page) GetLen() int {
	ret := 0
	for _, row := range p.rows {
		ret = ret + row.GetLen()
	}
	return ret
}

func (p *Page) copyRightPart(from *Page, index int) {
	p.rows = append(p.rows, from.rows[index:]...)
	p.size.SetValue(uint32(len(p.rows)))
	p._len = p.GetLen()
}

func (p *Page) deleteRightPart(index int) {
	p.rows = p.rows[:index]
	p.size.SetValue(uint32(len(p.rows)))
	p._len = p.GetLen()
}

func (p *Page) shouldSplit() bool {
	return p._len > DefaultPageSize
}

func (p *Page) isIndexPage() bool {
	return p.pgType.GetValue().(byte) == indexPageType
}
func (p *Page) isDataPage() bool {
	return p.pgType.GetValue().(byte) == dataPageType
}

func (p *Page) hasParent() bool {
	return p.parent != nil
}

func (p *Page) NewIndexRow() *Row {
	row := NewRow(dt.DefaultIndexRowMeta())
	row.WithDefaultValues()
	row.cells[0].SetValue(p.pgID.GetValue())
	row.SetKey(p.getMinKey())
	return row
}
func (p *Page) getMinKey() uint32 {
	return p.rows[0].GetKey()
}
