package mt

import (
	"fmt"
	"testing"
)

func TestMT64(t *testing.T) {
	mt := NewMT64()
	mt.InitByArray([]uint64{0x12345, 0x23456, 0x34567, 0x45678})
	fmt.Println("init_gen_rand__________")
	for i := 0; i < 1000; i++ {
		fmt.Printf("%20d ", mt.GenRand64())
		if i%5 == 4 {
			fmt.Println()
		}
	}
	fmt.Println()

	// 	sfmt.InitByArray([]uint32{0x1234, 0x5678, 0x9abc, 0xdef0})
	// 	fmt.Println("init_by_array__________")
	// 	for i := 0; i < 1000; i++ {
	// 		fmt.Printf("%10d ", sfmt.GenRand32())
	// 		if i%5 == 4 {
	// 			fmt.Println()
	// 		}
	// 	}
	// 	fmt.Println()
}
