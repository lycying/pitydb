package dt

var _indexRowMeta *RowMeta = nil

func init() {
	slot0 := NewCellMetaRaw(0, UInt32Type, "pageID", "its the pageID to find", uint32(0))
	_indexRowMeta = NewRowMeta()
	_indexRowMeta.AddCellMeta(slot0)
}

// DefaultIndexRowMeta is the default RowMeta for IndexPageRow
// It has only one instance because it can be use more times
func DefaultIndexRowMeta() *RowMeta {
	return _indexRowMeta
}
