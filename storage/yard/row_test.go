package yard

import (
	"github.com/lycying/pitydb/dt"
	"testing"
)

func TestNewRow(t *testing.T) {
	rowMeta := dt.NewRowMeta()
	slot0 := dt.NewCellMetaRaw(0, dt.UInt32Type, "id", "the auto incrementID", nil)
	slot1 := dt.NewCellMetaRaw(1, dt.Float64Type, "col2", "a flot64...", float64(0.9999999))
	slot2 := dt.NewCellMetaRaw(2, dt.StringType, "col3", "a string with default value pitydb...", "pitydb")

	rowMeta.AddCellMeta(slot0)
	rowMeta.AddCellMeta(slot1)
	rowMeta.AddCellMeta(slot2)

	row := NewRow(rowMeta)
	row.WithDefaultValues()
	//int32 + float64 + string def + byte (len) + string
	if row.GetLen() != (4 + 8 + 1 + 1 + len("pitydb")) {
		t.Fatal("the default value is not set because the length is not excepted!")
	}

	row.GetCellAt(slot1).SetValue(float64(-99999999.1234567890))
	row.GetCellAt(slot2).SetValue("cross the GreatWall we can reach every corner in the world!")

	buf, _ := row.Encode()

	row2 := NewRow(rowMeta)
	row2Len, _ := row2.Decode(buf, 0)

	if row2Len != len(buf) || len(buf) != row2.GetLen() {
		t.Fatal("row1 and row2 should have the same length")
	}

}
