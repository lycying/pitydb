package block

import (
	"testing"
	"os"
)

func TestBlock(t *testing.T) {
	f := "/tmp/b"
	file, _ := os.OpenFile(f, os.O_RDWR, 0666)
	b := NewBlock(0, file)
	b.Write(make([]byte, EightKB))
	for i := 0; i < 10000; i++ {
		b = NewBlock(int64(i), file)
		b.Write(make([]byte, EightKB))
	}
	file.Close()
}