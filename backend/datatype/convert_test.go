package datatype

import (
	"testing"
)

func TestDataType_Convert(t *testing.T) {
	//this should be like long in java
	{
		var val_int int64 = 10

		val_bytes := Int64ToBytes(val_int)
		t.Log("len should be 8 =", len(val_bytes))

		val_int_2 := BytesToInt64(val_bytes)
		t.Log(val_int_2 == val_int)
	}
	//this should be like int in java
	{
		var val_int int32 = 10

		val_bytes := Int32ToBytes(val_int)
		t.Log("len should be 4 =", len(val_bytes))

		val_int_2 := BytesToInt32(val_bytes)
		t.Log(val_int_2 == val_int)
	}
	{
		str := "One2龙" //3+1+3 = 7
		val_bytes := StringToBytes(str)
		t.Log("len should be 7 =", len(val_bytes))
		t.Log(BytesToString(val_bytes), str == BytesToString(val_bytes))

	}


}
