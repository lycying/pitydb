package yard

import (
	"bytes"
	"github.com/lycying/pitydb/dt"
)

type RowRefer interface {
	dt.Encoder
	dt.DeCoder

	GetLen() int
}
type Row struct {
	RowRefer

	meta  *dt.RowMeta  //meta data for loop data
	key   *dt.UInt32   //the key used for b+ tree
	cells []dt.DtRefer //the data part
}

func NewRow(meta *dt.RowMeta) *Row {
	return &Row{
		meta:  meta,
		key:   dt.NewUInt32(),
		cells: make([]dt.DtRefer, 0),
	}
}

func (r *Row) WithDefaultValues() {
	idx := int(0)
	data := make([]dt.DtRefer, r.meta.GetCellSize())
	for i, item := range r.meta.GetItems() {
		cell := item.NewCell()
		if nil != item.GetDefaultValue() {
			cell.SetValue(item.GetDefaultValue())
		}
		data[i] = cell
		idx += cell.GetLen()
	}
	r.cells = data
}

func (r *Row) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, item := range r.cells {
		b, _ := item.Encode()
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func (r *Row) Decode(buf []byte, offset int) (int, error) {
	idx := int(0)
	data := make([]dt.DtRefer, r.meta.GetCellSize())

	for i, item := range r.meta.GetItems() {
		cell := item.NewCell()
		cLen, _ := cell.Decode(buf, idx+offset)
		data[i] = cell
		idx += cLen
	}
	r.cells = data
	return idx, nil
}

func (r *Row) GetLen() int {
	cl := 0
	for _, item := range r.cells {
		cl += item.GetLen()
	}
	return cl
}

func (r *Row) GetCellAt(meta *dt.CellMeta) dt.DtRefer {
	return r.cells[meta.GetPos()]
}

func (r *Row) GetKey() uint32 {
	return r.key.GetValue().(uint32)
}

func (r *Row) SetCellValueForTest(meta *dt.CellMeta, value dt.ValueRefer) {
	if nil != value {
		r.GetCellAt(meta).SetValue(value)
	}
}

func (r *Row) SetKey(key uint32) {
	r.key.SetValue(key)
}
