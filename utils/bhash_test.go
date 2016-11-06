package utils

import (
	"fmt"
	"testing"
)

func TestHashKey(t *testing.T) {
	fmt.Printf("0x%X \n", HashString("unitneutralacritter.grp", 0))
	fmt.Printf("0x%X \n", HashString("unitneutralacritter.grp", 1))
	fmt.Printf("0x%X \n", HashString("unitneutralacritter.grp", 2))
	fmt.Printf("0x%X \n", HashString("unitneutralacritter.grp", 3))
	fmt.Printf("0x%X \n", HashString("unitneutralacritter.grp", 4))
}
