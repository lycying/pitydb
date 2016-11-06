package skiplist

type node struct {
	ptr  *column
	span int
}

func newNode() *node {
	return new(node)
}

type column struct {
	forward  []*node
	backward *column
	k, v     interface{}
}

func newColumn(level int) *column {
	ele := new(column)
	ele.backward = nil
	ele.forward = make([]*node, level)
	for i := 0; i < level; i++ {
		ele.forward[i] = newNode()
	}
	return ele
}
