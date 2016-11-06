package dt

type CellMeta struct {
	pos          *Int32
	name         *String
	comment      *String
	mType        *Int32
	defaultValue interface{}
}

func NewCellMetaRaw(pos int, typ DType, name string, comment string, defaultValue interface{}) *CellMeta {
	tPos := NewInt32()
	tPos.SetValue(int32(pos))
	tName := NewString()
	tName.SetValue(name)
	tComment := NewString()
	tComment.SetValue(comment)
	tTyp := NewInt32()
	tTyp.SetValue(int32(typ))

	return NewCellMeta(tPos, tTyp, tName, tComment, defaultValue)
}

func NewCellMeta(pos *Int32, mType *Int32, name *String, comment *String, defaultValue interface{}) *CellMeta {
	return &CellMeta{
		pos:          pos,
		name:         name,
		comment:      comment,
		mType:        mType,
		defaultValue: defaultValue,
	}
}

func (s *CellMeta) GetPos() int {
	return int(s.pos.value)
}

func (s *CellMeta) GetName() string {
	return s.name.value
}
func (s *CellMeta) GetComment() string {
	return s.comment.value
}
func (s *CellMeta) GetMType() DType {
	return DType(s.mType.value)
}
func (s *CellMeta) GetDefaultValue() interface{} {
	return s.defaultValue
}

func (s *CellMeta) NewCell() DtRefer {
	r := NewDtRefer(DType(s.mType.value))
	if s.defaultValue != nil {
		r.SetValue(s.defaultValue)
	}
	return r
}

type RowMeta struct {
	items   []*CellMeta
	comment *String
}

func NewRowMeta() *RowMeta {
	return &RowMeta{
		items:   make([]*CellMeta, 0),
		comment: NewString(),
	}
}

func (meta *RowMeta) SetComment(str string) {
	meta.comment.SetValue(str)
}

func (meta *RowMeta) AddCellMeta(item *CellMeta) {
	pos := int(item.pos.value)
	oldLen := len(meta.items)
	if pos >= oldLen {
		newItems := make([]*CellMeta, pos+1)
		if pos > 0 && oldLen > 0 {
			copy(newItems[:oldLen], meta.items)
		}
		meta.items = newItems
	}
	meta.items[pos] = item
}

func (meta *RowMeta) GetItems() []*CellMeta {
	return meta.items
}

func (meta *RowMeta) GetCellSize() int {
	return len(meta.items)
}
