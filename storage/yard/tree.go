package yard

import (
	"github.com/lycying/pitydb/dt"
	"os"
)

type PageTree struct {
	root *Page
	meta *dt.RowMeta
	link *os.File
	mgr  *PageManagement
}

func NewPageTree(meta *dt.RowMeta, link *os.File) *PageTree {
	mgr := NewPageMgr()
	tree := &PageTree{
		meta: meta,
		link: link,
		mgr:  mgr,
	}
	root := tree.NewDataPage(0)

	tree.root = root

	return tree

}
func (tree *PageTree) NewIndexPage(level uint32) *Page {
	return tree.NewPage(level, indexPageType)
}
func (tree *PageTree) NewDataPage(level uint32) *Page {
	return tree.NewPage(level, dataPageType)
}
func (tree *PageTree) NewPage(lvl uint32, typ byte) *Page {
	pgID := dt.NewUInt32()
	pgID.SetValue(tree.mgr.NextPageID())
	pgType := dt.NewByte()
	pgType.SetValue(typ)
	level := dt.NewUInt32()
	level.SetValue(lvl)
	left := dt.NewUInt32()
	left.SetValue(uint32(0))
	right := dt.NewUInt32()
	right.SetValue(uint32(0))
	checksum := dt.NewUInt32()
	checksum.SetValue(uint32(0)) //TODO
	lastModify := dt.NewUInt64()
	lastModify.SetValue(uint64(0))
	size := dt.NewUInt32()
	size.SetValue(uint32(0))
	pg := &Page{
		pageHeader: pageHeader{
			pgID:       pgID,
			pgType:     pgType,
			level:      level,
			left:       left,
			right:      right,
			checksum:   checksum,
			lastModify: lastModify,
			size:       size,
		},
		_len:   0,
		tree:   tree,
		parent: nil,
		rows:   []*Row{},
	}
	tree.mgr.AddPage(pg)
	return pg
}

func (tree *PageTree) Insert(r *Row) {
	key := r.GetKey()

	node, idx, find := tree.root.findOne(key)

	//the row is so big that one default can not hold it
	if r.GetLen() > DefaultPageSize {
		//TODO big row storage
	}
	node.insert(r, idx, find)

}
func (tree *PageTree) Delete(key uint32) bool {
	node, idx, find := tree.root.findOne(key)
	if find {
		node.delete(key, idx)
		return true
	}
	return false
}

func (tree *PageTree) GetRoot() *Page {
	return tree.root
}
func (tree *PageTree) Dump() {
	println("BEGIN")
	root := tree.root
	_dump(root)
	println("END")
	println("")
}
func _getParentID(pg *Page) uint32 {
	if pg.parent == nil {
		return 0
	} else {
		return pg.parent.pgID.GetValue().(uint32)
	}
}
func _dump(pg *Page) {
	if nil == pg {
		return
	}
	if pg.isDataPage() {
		print(pg.level.GetValue().(uint32), "D`", pg.pgID.GetValue().(uint32), "@", _getParentID(pg), "`\t:(")
		for _, x := range pg.rows {
			print(x.GetKey(), ",")
		}
		print(")")
		println()
	} else {
		print(pg.level.GetValue().(uint32), "I`", pg.pgID.GetValue().(uint32), "@", _getParentID(pg), "`\t:[")
		for _, x := range pg.rows {
			print(x.GetKey(), ",")
		}
		print("]")
		println()
		for _, x := range pg.rows {
			px := pg.tree.mgr.GetPage(x.cells[0].GetValue().(uint32))
			_dump(px)
		}
	}
}
