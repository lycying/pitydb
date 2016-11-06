package utils

import "testing"

func TestSum32(t *testing.T) {
	b := make([]byte, 1024*16)

	if 2874462854 != Sum32(b) {
		t.Fatal("sum32 failed")
	}
}
