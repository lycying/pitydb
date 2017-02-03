package btree

import (
	"encoding/binary"
	"math/rand"
	"os"
	"testing"
)

const (
	testFileName  = "/tmp/test.btree"
	benchFileName = "/tmp/bench.btree"
)

var (
	capacity uint
	count    int
	delta    uint
)

type key struct {
	K uint
	V uint
}

func (this *key) Compare(buf []byte) (int, error) {
	k := uint(binary.LittleEndian.Uint32(buf[:4]))
	r := this.K - k
	return int(r), nil
}

func (this key) Size() uint {
	return 8
}

func (this *key) Read(buf []byte) error {
	this.K = uint(binary.LittleEndian.Uint32(buf[:4]))
	this.V = uint(binary.LittleEndian.Uint32(buf[4:]))
	return nil
}

func (this *key) Write(buf []byte) error {
	binary.LittleEndian.PutUint32(buf[:4], uint32(this.K))
	binary.LittleEndian.PutUint32(buf[4:], uint32(this.V))
	return nil
}

var magic [16]byte = [16]byte{'T', 'e', 's', 't', 'B', 'T', 'r', 'e', 'e'}
var bt Tree
var testMap map[uint]*uint

func Test10(t *testing.T) {
	capacity = 10
	count = 1000
	delta = 1
	testMap = make(map[uint]*uint, count)
	testCreate(t)
	t.Log("create test done")
	testOpen(t)
	t.Log("open test done")
	testInsert(t)
	t.Log("insert test done")
	testFind(t)
	t.Log("find test done")
	testUpdate(t)
	t.Log("update test done")
	testFind(t)
	t.Log("find after update test done")
	testEnum(t)
	t.Log("enum test done")
	testReverseEnum(t)
	t.Log("reverse enum test done")
	testDelete(t)
	t.Log("delete test done")
	testEnum(t)
	t.Log("enum test done")
	testReverseEnum(t)
	t.Log("reverse enum test done")
	count *= 2
	testInsert(t)
	t.Log("second insert test done")
	testEnum(t)
	t.Log("second enum test done")
	testReverseEnum(t)
	t.Log("second reverse enum test done")
	t.Log("second enum test done")
	testDelete(t)
	t.Log("second delete test done")
	testEnum(t)
	t.Log("second enum test done")
	testReverseEnum(t)
	t.Log("second reverse enum test done")
}

func Test2(t *testing.T) {
	capacity = 2
	count = 1000
	delta = 10
	testMap = make(map[uint]*uint, count)
	testCreate(t)
	t.Log("create test done")
	testOpen(t)
	t.Log("open test done")
	testInsert(t)
	t.Log("insert test done")
	testFind(t)
	t.Log("find test done")
	testUpdate(t)
	t.Log("update test done")
	testFind(t)
	t.Log("find after update test done")
	testEnum(t)
	t.Log("enum test done")
	testReverseEnum(t)
	t.Log("reverse enum test done")
	testDelete(t)
	t.Log("delete test done")
}

func Test100(t *testing.T) {
	capacity = 1000
	count = 10000
	delta = 100
	testMap = make(map[uint]*uint, count)
	testCreate(t)
	t.Log("create test done")
	testOpen(t)
	t.Log("open test done")
	testInsert(t)
	t.Log("insert test done")
	testFind(t)
	t.Log("find test done")
	testUpdate(t)
	t.Log("update test done")
	testFind(t)
	t.Log("find after update test done")
	testEnum(t)
	t.Log("enum test done")
	testReverseEnum(t)
	t.Log("reverse enum test done")
	testDelete(t)
	t.Log("delete test done")
}

func testCreate(t *testing.T) {
	f, err := os.Create(testFileName)
	if err != nil {
		t.Fatal(err)
	}
	_, err = NewBTree(f, magic, &key{0, 0}, capacity)
	if err != nil {
		t.Fatal(err)
	}
}

func testOpen(t *testing.T) {
	f, err := os.OpenFile(testFileName, os.O_RDWR, 0666)
	if err != nil {
		t.Fatal(err)
	}
	bt, err = OpenBTree(f, magic, &key{0, 0})
	if err != nil {
		t.Fatal(err)
	}
}

func testInsert(t *testing.T) {
	for i := 0; i < count; {
		r := uint(rand.Int31())
		if _, found := testMap[r]; found {
			continue
		}

		if k, err := bt.Insert(&key{r, r}); err != nil {
			t.Fatal(err)
		} else if k != nil {
			t.Fatalf("%#v is already inserted", k)
		}
		testMap[r] = &r
		i++
	}
	for k, v := range testMap {
		if r, err := bt.Insert(&key{k, *v}); err != nil {
			t.Fatal(err)
		} else if r == nil {
			t.Fatalf("duplicate %#v has been inserted", k)
		}
	}
}

func testFind(t *testing.T) {
	for k, v := range testMap {
		r, err := bt.Find(&key{k, 0})
		if err != nil {
			t.Fatal(err)
		}
		if r.(*key).V != *v {
			t.Fatalf("result of find %#v is mismatch: %#v, must be %#v\n", key{k, k}, r, *v)
		}
	}
	for i := 0; i < count; i++ {
		r := uint(rand.Int31())
		if _, found := testMap[r]; found {
			continue
		}
		if k, err := bt.Find(&key{r, 0}); err != nil {
			t.Fatal(err)
		} else if k != nil {
			t.Fatalf("result of find %#v is mismatch: %#v, must be nil\n", key{r, 0}, k)
		}
	}
}

func testUpdate(t *testing.T) {
	for k, v := range testMap {
		r, err := bt.Update(&key{k, (*v) + delta})
		if err != nil {
			t.Fatal(err)
		}
		if r == nil || r.(*key).V != *v {
			t.Fatalf("result of update is mismatch: %#v, must be %#v\n", r, *v)
		}
		(*v) = (*v) + delta
	}
}

func testDelete(t *testing.T) {
	count := 0
	for k, v := range testMap {
		r, err := bt.Delete(&key{k, 0})
		if err != nil {
			t.Fatal(err)
		}
		if r == nil || *v != r.(*key).V {
			t.Fatalf("result of delete is mismatch: %#v, must be %#v\n", r, *v)
		}
		count++
		delete(testMap, k)
	}
}

func testEnum(t *testing.T) {
	count := 0
	var last uint = 0
	f := bt.Enum(nil)
	var begin uint
	for k, e := f(); k != nil && e == nil; k, e = f() {
		if k.(*key).K <= last {
			t.Fatalf("wrong sequence of keys: current key: %#v, previous key: %#v\n", k, last)
		}
		last = k.(*key).K
		if count == len(testMap)/2 {
			begin = last
		}
		if v, found := testMap[last]; !found {
			t.Fatalf("key not found: %#v\n", last)
		} else if k.(*key).V != *v {
			t.Fatalf("value mismatch for key %#v, must be %#v\n", k, *v)
		}
		count++
	}
	if count != len(testMap) {
		t.Fatalf("count of values mismatch: %#v, must be %#v\n", count, len(testMap))
	}
	count = 0
	last = 0
	f = bt.Enum(nil)
	for k, e := f(); k != nil && e == nil && k.(*key).K < begin; k, e = f() {
		if k.(*key).K <= last {
			t.Fatalf("wrong sequence of keys: current key: %#v, previous key: %#v\n", k, last)
		}
		last = k.(*key).K
		if v, found := testMap[last]; !found {
			t.Fatalf("key not found: %#v\n", last)
		} else if k.(*key).V != *v {
			t.Fatalf("value mismatch for key %#v, must be %#v\n", k, *v)
		}
		count++
	}
	var end uint = begin + 1000000
	f = bt.Enum(&key{begin, 0})
	for k, e := f(); k != nil && e == nil && k.(*key).K < end; k, e = f() {
		if k.(*key).K <= last {
			t.Fatalf("wrong sequence of keys: current key: %#v, previous key: %#v\n", k, last)
		}
		last = k.(*key).K
		if v, found := testMap[last]; !found {
			t.Fatalf("key not found: %#v\n", last)
		} else if k.(*key).V != *v {
			t.Fatalf("value mismatch for key %#v, must be %#v\n", k, *v)
		}
		count++
	}
	f = bt.Enum(&key{end, 0})
	for k, e := f(); k != nil && e == nil; k, e = f() {
		if k.(*key).K <= last {
			t.Fatalf("wrong sequence of keys: current key: %#v, previous key: %#v\n", k, last)
		}
		last = k.(*key).K
		if v, found := testMap[last]; !found {
			t.Fatalf("key not found: %#v\n", last)
		} else if k.(*key).V != *v {
			t.Fatalf("value mismatch for key %#v, must be %#v\n", k, *v)
		}
		count++
	}
	if count != len(testMap) {
		t.Fatalf("count of values mismatch: %#v, must be %#v\n", count, len(testMap))
	}
}

func testReverseEnum(t *testing.T) {
	count := 0
	var last uint = 0xFFFFFFFF
	var begin uint
	f := bt.ReverseEnum(nil)
	for k, e := f(); k != nil && e == nil; k, e = f() {
		if k.(*key).K > last {
			t.Fatalf("wrong sequence of keys: current key: %#v, previous key: %#v\n", k, last)
		}
		last = k.(*key).K
		if count == len(testMap)/2 {
			begin = last
		}
		if v, found := testMap[last]; !found {
			t.Fatalf("key not found: %#v\n", last)
		} else if k.(*key).V != *v {
			t.Fatalf("value mismatch for key %#v, must be %#v\n", k, *v)
		}
		count++
	}
	if count != len(testMap) {
		t.Fatalf("count of values mismatch: %#v, must be %#v\n", count, len(testMap))
	}
	count = 0
	last = 0xFFFFFFFF
	f = bt.ReverseEnum(nil)
	for k, e := f(); k != nil && e == nil && k.(*key).K > begin; k, e = f() {
		if k.(*key).K > last {
			t.Fatalf("wrong sequence of keys: current key: %#v, previous key: %#v\n", k, last)
		}
		last = k.(*key).K
		if v, found := testMap[last]; !found {
			t.Fatalf("key not found: %#v\n", last)
		} else if k.(*key).V != *v {
			t.Fatalf("value mismatch for key %#v, must be %#v\n", k, *v)
		}
		count++
	}
	var end uint = begin - 100000000
	f = bt.ReverseEnum(&key{begin, 0})
	for k, e := f(); k != nil && e == nil && k.(*key).K > end; k, e = f() {
		if k.(*key).K > last {
			t.Fatalf("wrong sequence of keys: current key: %#v, previous key: %#v\n", k, last)
		}
		last = k.(*key).K
		if v, found := testMap[last]; !found {
			t.Fatalf("key not found: %#v\n", last)
		} else if k.(*key).V != *v {
			t.Fatalf("value mismatch for key %#v, must be %#v\n", k, *v)
		}
		count++
	}
	f = bt.ReverseEnum(&key{end, 0})
	for k, e := f(); k != nil && e == nil; k, e = f() {
		if k.(*key).K > last {
			t.Fatalf("wrong sequence of keys: current key: %#v, previous key: %#v\n", k, last)
		}
		last = k.(*key).K
		if v, found := testMap[last]; !found {
			t.Fatalf("key not found: %#v\n", last)
		} else if k.(*key).V != *v {
			t.Fatalf("value mismatch for key %#v, must be %#v\n", k, *v)
		}
		count++
	}
	if count != len(testMap) {
		t.Fatalf("count of values mismatch: %#v, must be %#v\n", count, len(testMap))
	}
}

var benchList []uint

func BenchmarkInsert(b *testing.B) {
	b.StopTimer()
	capacity = 100
	count = 100000
	b.N = count
	delta = 1
	testMap := make(map[uint]uint, count)
	benchList = make([]uint, 0, count)
	for i := 0; i < count; {
		r := uint(rand.Int31())
		if _, found := testMap[r]; found {
			continue
		}
		i++
		testMap[r] = r
		benchList = append(benchList, r)
	}
	testMap = nil
	f, err := os.Create(benchFileName)
	if err != nil {
		panic(err)
	}
	bt, err = NewBTree(f, magic, &key{0, 0}, capacity)
	if err != nil {
		panic(err)
	}
	for i := 0; i < count; i++ {
		b.StartTimer()
		r, err := bt.Insert(&key{benchList[i], benchList[i]})
		b.StopTimer()
		if err != nil {
			panic(err)
		} else if r != nil {
			panic(r)
		}
	}
}

func BenchmarkFind(b *testing.B) {
	b.StopTimer()
	b.N = count
	for i := 0; i < count; i++ {
		b.StartTimer()
		k, err := bt.Find(&key{benchList[i], 0})
		b.StopTimer()
		if err != nil {
			panic(err)
		} else if k == nil {
			panic(k)
		}
	}
}

func BenchmarkFailedFind(b *testing.B) {
	b.StopTimer()
	b.N = count
	for i := 0; i < count; i++ {
		b.StartTimer()
		_, err := bt.Find(&key{benchList[i] + 1, 0})
		b.StopTimer()
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkUpdate(b *testing.B) {
	b.StopTimer()
	b.N = count
	for i := 0; i < count; i++ {
		b.StartTimer()
		k, err := bt.Update(&key{benchList[i], benchList[i] + delta})
		b.StopTimer()
		if err != nil {
			panic(err)
		} else if k == nil {
			panic(k)
		}
	}
}

func BenchmarkDelete(b *testing.B) {
	b.StopTimer()
	b.N = count
	for i := 0; i < count; i++ {
		b.StartTimer()
		k, err := bt.Delete(&key{benchList[i], 0})
		b.StopTimer()
		if err != nil {
			panic(err)
		} else if k == nil {
			panic(k)
		}
	}
}
