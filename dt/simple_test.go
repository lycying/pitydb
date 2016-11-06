package dt

import (
	"bytes"
	"testing"
)

func TestInt32(t *testing.T) {
	i1 := int32(-999)
	mI1 := ValidNewInt32(i1)
	b1, _ := mI1.Encode()

	mI2 := NewInt32()
	lenMI2, _ := mI2.Decode(b1, 0)

	if lenMI2 != 4 || lenMI2 != mI2.GetLen() {
		t.Fatal("lenMI2 should be 4,but ", lenMI2)
	}
	if mI2.GetValue().(int32) != -999 {
		t.Fatal("mI2 value should be -99,but ", mI2.GetValue().(int32))
	}
}

func TestUInt32(t *testing.T) {
	i1 := uint32(999)
	mI1 := NewUInt32()
	mI1.SetValue(i1)
	b1, _ := mI1.Encode()

	mI2 := NewUInt32()
	lenMI2, _ := mI2.Decode(b1, 0)

	if lenMI2 != 4 || lenMI2 != mI2.GetLen() {
		t.Fatal("lenMI2 should be 4,but ", lenMI2)
	}
	if mI2.GetValue().(uint32) != 999 {
		t.Fatal("mI2 value should be 999,but ", mI2.GetValue().(uint32))
	}
}
func TestInt64(t *testing.T) {
	i1 := int64(-999)
	mI1 := NewInt64()
	mI1.SetValue(i1)
	b1, _ := mI1.Encode()

	mI2 := NewInt64()
	lenMI2, _ := mI2.Decode(b1, 0)

	if lenMI2 != 8 || lenMI2 != mI2.GetLen() {
		t.Fatal("lenMI2 should be 8,but ", lenMI2)
	}
	if mI2.GetValue().(int64) != -999 {
		t.Fatal("mI2 value should be -99,but ", mI2.GetValue().(int64))
	}
}

func TestUInt64(t *testing.T) {
	i1 := uint64(999)
	mI1 := NewUInt64()
	mI1.SetValue(i1)
	b1, _ := mI1.Encode()

	mI2 := NewUInt64()
	lenMI2, _ := mI2.Decode(b1, 0)

	if lenMI2 != 8 || lenMI2 != mI2.GetLen() {
		t.Fatal("lenMI2 should be 8,but ", lenMI2)
	}
	if mI2.GetValue().(uint64) != 999 {
		t.Fatal("mI2 value should be 999,but ", mI2.GetValue().(uint64))
	}
}

func TestFloat64(t *testing.T) {
	i1 := float64(999.000000001)
	mI1 := NewFloat64()
	mI1.SetValue(i1)
	b1, _ := mI1.Encode()

	mI2 := NewFloat64()
	lenMI2, _ := mI2.Decode(b1, 0)

	if lenMI2 != 8 || lenMI2 != mI2.GetLen() {
		t.Fatal("lenMI2 should be 8,but ", lenMI2)
	}
	if mI2.GetValue().(float64) != i1 {
		t.Fatal("mI2 value should be 999.000000001 but ", mI2.GetValue().(float64))
	}
}
func TestFloat32(t *testing.T) {
	i1 := float32(-999.909)
	mI1 := NewFloat32()
	mI1.SetValue(i1)
	b1, _ := mI1.Encode()

	mI2 := NewFloat32()
	lenMI2, _ := mI2.Decode(b1, 0)

	if lenMI2 != 4 || lenMI2 != mI2.GetLen() {
		t.Fatal("lenMI2 should be 4,but ", lenMI2)
	}
	if mI2.GetValue().(float32) != float32(-999.909) {
		t.Fatal("mI2 value should be -99.909,but ", mI2.GetValue().(float32))
	}
}

func TestByte(t *testing.T) {
	i1 := byte(0xfe)
	mI1 := NewByte()
	mI1.SetValue(i1)
	b1, _ := mI1.Encode()

	mI2 := NewByte()
	lenMI2, _ := mI2.Decode(b1, 0)

	if lenMI2 != 1 || lenMI2 != mI2.GetLen() {
		t.Fatal("lenMI2 should be 1,but ", lenMI2)
	}
	if mI2.GetValue().(byte) != byte(0xfe) {
		t.Fatal("mI2 value should be 0xfe ,but ", mI2.GetValue().(byte))
	}
}
func TestBool(t *testing.T) {
	i1 := true
	mI1 := NewBool()
	mI1.SetValue(i1)
	b1, _ := mI1.Encode()

	mI2 := NewBool()
	lenMI2, _ := mI2.Decode(b1, 0)

	if lenMI2 != 1 || lenMI2 != mI2.GetLen() {
		t.Fatal("lenMI2 should be 1,but ", lenMI2)
	}
	if mI2.GetValue().(bool) != true {
		t.Fatal("mI2 value should be true ,but ", mI2.GetValue().(bool))
	}
}

func TestString(t *testing.T) {
	str1 := string("太阳落山了")
	mStr1 := ValidNewString(str1)
	b1, _ := mStr1.Encode()

	toBe := []byte{0x01, 0x0f, 0xe5, 0xa4, 0xaa, 0xe9, 0x98, 0xb3, 0xe8, 0x90, 0xbd, 0xe5, 0xb1, 0xb1, 0xe4, 0xba, 0x86}
	if 0 != bytes.Compare(toBe, b1) {
		t.Fatalf("should be %x ,but % x ! \n", toBe, b1)
	}

	str := `duang
	犀 隘 媒 媚 婿 缅 缆 缔 缕 骚 搀 搓 壹 搔 葫 募 蒋 蒂 韩 棱 椰 焚 椎 棺 榔 椭 粟 棘 酣 酥 硝 硫
	تشينغ  مسؤولية  اختيار  دان  تان  و  سحب  سحب  سحب  الفيلم الذي  اتجه إلى  أعلى  ضد  اعتقال  وتفكيك  عقد  عقد  مع  سحب  اعتراض  المحتملة  . 
	봉 놀다 링 무 푸른 책임을 현 시계 게이지 바르다 골라 뽑아 모으다 뽑다 지고 평탄한 걸어 꺾어 끌고 찍은 자 위 뜯어 밀려 저당 체포하여 퍼텐셜 안고 오물 끌고 막아 비비다
	다행히 모집 경사진 걸쳐 지출하다 고르다 들고 그 고통을 받다, 만약 무성하다 사과 모 英范 줄곧 줄기 가지 논하다 林枝 잔 카운터 분석 보드 소나무 총 구상 은결이 베고 진술하다.
	奉遊び環武青責任表現を拾ったり規則を拭いタン押吸って曲がって頂をはずして引っ張って落札者に拘勢を抱いてごみを遮る和え
	坂を上げた幸を選択を取るならその苦い茂規程（平成11年12月苗英范直と茄子の莖茅林枝杯櫃析板松銃構傑のように枕
	喪は棗になって棗を刺して棗を売ることができます
	duangduangduangduang`
	mStr1.SetValue(str)
	b2, _ := mStr1.Encode()

	mStr2 := NewString()
	mStr2.Decode(b2, 0)

	if str != mStr2.GetValue().(string) {
		t.Fatalf("mStr2 should be %s , but %s \n", str, mStr2.GetValue().(string))
	}
}

func TestTableMeta(t *testing.T) {
	tableMeta := NewRowMeta()
	slot0 := NewCellMetaRaw(0, UInt32Type, "id", "the auto incrementID", nil)
	slot1 := NewCellMetaRaw(1, Int64Type, "col1", "col1...", int64(-1000))
	slot2 := NewCellMetaRaw(2, BoolType, "col2", "col2...", false)
	slot3 := NewCellMetaRaw(3, ByteType, "col3", "col3...", byte(0xff))
	slot4 := NewCellMetaRaw(4, Float64Type, "col4", "col4...", float64(0.9999999))
	slot5 := NewCellMetaRaw(5, StringType, "col5", "col5...", "very very pity!")

	tableMeta.AddCellMeta(slot2)

	tableMeta.AddCellMeta(slot0)
	tableMeta.AddCellMeta(slot1)
	tableMeta.AddCellMeta(slot2)
	tableMeta.AddCellMeta(slot3)
	tableMeta.AddCellMeta(slot5)
	tableMeta.AddCellMeta(slot4)

	if 6 != tableMeta.GetCellSize() {
		t.Fatal("table size should be 6, but:", tableMeta.GetCellSize())
	}

}

func TestDefaultIndexRowMeta(t *testing.T) {
	rMeta := DefaultIndexRowMeta()
	if UInt32Type != rMeta.GetItems()[0].GetMType() {
		t.Fatal("meta type should be PTypeUInt32")
	}
}
