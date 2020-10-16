package gosfmt

import (
	"fmt"
	"testing"
)

func TestSFMT64_604(t *testing.T) {
	sfmt32 := NewSFMT32(params604)
	fmt.Printf("%s\n64 bit generated randoms with SFMT32\n", sfmt32.params.idstr)
	sfmt32.initGenRand(4321)
	fmt.Println("init_gen_rand__________")
	for i := 0; i < 10; i++ {
		fmt.Printf("%20d ", sfmt32.genRand64())
		if i%5 == 4 {
			fmt.Println()
		}
	}
	fmt.Println()

	sfmt := NewSFMT64(params604)
	fmt.Printf("%s\n64 bit generated randoms\n", sfmt.params.idstr)
	sfmt.initGenRand(4321)
	fmt.Println("init_gen_rand__________")
	for i := 0; i < 10; i++ {
		fmt.Printf("%20d ", sfmt.genRand64())
		if i%5 == 4 {
			fmt.Println()
		}
	}
	fmt.Println()

	// sfmt.initByArray([]uint32{0x1234, 0x5678, 0x9abc, 0xdef0})
	// fmt.Println("init_by_array__________")
	// for i := 0; i < 1000; i++ {
	// 	fmt.Printf("%10d ", sfmt.genRand32())
	// 	if i%5 == 4 {
	// 		fmt.Println()
	// 	}
	// }
	// fmt.Println()
}
