package skiplist_test

import (
	"fmt"
	"github.com/lycying/pitydb/lib/skiplist"
	"testing"
	"time"
)

func TestSkipList_Find(t *testing.T) {
	sk := skiplist.New(func(l, r interface{}) bool {
		a := l.(int)
		b := r.(int)
		return a < b
	})

	start := time.Now()
	for i := 0; i < 100; i++ {
		sk.Add(i, fmt.Sprintf("pitydb-%d", i))
	}
	fmt.Println("add 0-100 cost:", time.Now().Sub(start).String())
	start = time.Now()
	for i := 2000000; i > 50; i-- {
		sk.Add(i, fmt.Sprintf("pitydb-%d", i))
	}
	fmt.Println("add 50-2000000 cost:", time.Now().Sub(start).String())
	start = time.Now()
	//sk.Dump()
	fmt.Println("length:", sk.Len())
	fmt.Println("level :", sk.Level())
	fmt.Printf("find %v:%v\n", 113, sk.Find(113))
	fmt.Println("find 113 cost:", time.Now().Sub(start).String())
	start = time.Now()
	fmt.Printf("find %v:%v\n", 99999, sk.Find(99999))
	fmt.Println("find 99999 cost:", time.Now().Sub(start).String())
	start = time.Now()
	fmt.Printf("find %v:%v\n", 1000090, sk.Find(1000090))
	fmt.Println("find 1000090 cost:", time.Now().Sub(start).String())
	start = time.Now()
	fmt.Printf("find %v:%v\n", -1, sk.Find(-1))
	sk.Del(20000000)
	sk.Del(2000000)
	sk.Del(0)
	sk.Del(-1)
	for i := 2000; i > 999; i-- {
		sk.Del(i)
	}
	fmt.Printf("find %v:%v\n", 0, sk.Find(0))
	fmt.Printf("find %v:%v\n", 2000000, sk.Find(2000000))
	fmt.Printf("find %v:%v\n", 1000, sk.Find(1000))
	fmt.Printf("find %v:%v\n", 999, sk.Find(999))
	//sk.Dump()
}
