package btree

const K_SIZE = 3 //度

type Item interface {
	Less(than Item) bool
	Key() int
}
//对于文件系统来说，载入的node即为页，所以int类型存放的是页在其中的偏移量
type internal_node struct {
	prev     int
	next     int
	parent   int
	num      int
	children []int
	is_leaf  bool
	key      int
}

type BPlusTree struct {
	root   *internal_node
	degree int
}

func New() *BPlusTree {
	root := &internal_node{
		prev:0,
		next:0,
		parent:0,
		num:0,
		children:nil,
		is_leaf:false,
	}
	tree := &BPlusTree{
		root:root,
		degree:3,
	}
	return tree
}

func (this *BPlusTree) Insert(item Item) *internal_node {
	if this.root.num == 0 {
		this.root.num++
		this.root.is_leaf = true
		this.root.key = item.Key()
		return this.root
	}
	return nil
}

func (this *BPlusTree) Search(key int) {
}
