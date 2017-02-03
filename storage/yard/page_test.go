package yard

import (
	"github.com/lycying/pitydb/dt"
	"os"
	"testing"
)

func TestNewPage(t *testing.T) {
	rowMeta := dt.NewRowMeta()
	slot0 := dt.NewCellMetaRaw(0, dt.UInt32Type, "id", "the auto incrementID", nil)
	slot1 := dt.NewCellMetaRaw(1, dt.Float64Type, "col2", "a flot64...", float64(0.9999999))
	slot2 := dt.NewCellMetaRaw(2, dt.StringType, "col3", "a string with default value pitydb...", "pitydb")

	rowMeta.AddCellMeta(slot0)
	rowMeta.AddCellMeta(slot1)
	rowMeta.AddCellMeta(slot2)

	link, _ := os.OpenFile("/tmp/b", os.O_RDWR, 0666)
	tree := NewPageTree(rowMeta, link)

	for i := 1; i <= 1; i++ {
		r := NewRow(rowMeta)
		r.WithDefaultValues()
		r.SetKey(uint32(i))
		r.SetCellValueForTest(rowMeta.GetItems()[0], uint32(i))
		r.SetCellValueForTest(rowMeta.GetItems()[1], -99999999.99999999)
		r.SetCellValueForTest(rowMeta.GetItems()[2], "Hard work make bird stupid!")

		tree.Insert(r)
	}

	for i := 100; i >= 1; i-- {
		r := NewRow(rowMeta)
		r.WithDefaultValues()
		r.SetKey(uint32(i))
		r.SetCellValueForTest(rowMeta.GetItems()[0], uint32(i))
		r.SetCellValueForTest(rowMeta.GetItems()[1], -99999999.99999999)
		r.SetCellValueForTest(rowMeta.GetItems()[2], "Hard work make bird stupid!")

		tree.Insert(r)
	}
	tree.Dump()
	//for i := 1; i < 20000; i++ {
	//	_, _, found := tree.root.findOne(uint32(i))
	//	assert.Equal(t, found, true)
	//}

}
