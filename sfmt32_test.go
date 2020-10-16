package gosfmt

import (
	"fmt"
	"testing"
)

func TestSFMT32_604(t *testing.T) {
	sfmt := NewSFMT32(params604)
    fmt.Printf("%s\n32 bit generated randoms\n", sfmt.params.idstr)
	sfmt.initGenRand(1234)
	fmt.Println("init_gen_rand__________")
	for i := 0; i < 1000; i++ {
		fmt.Printf("%10d ", sfmt.genRand32())
		if i%5 == 4 {
			fmt.Println()
		}
	}
	fmt.Println()

	sfmt.initByArray([]uint32{0x1234, 0x5678, 0x9abc, 0xdef0})
	fmt.Println("init_by_array__________")
	for i := 0; i < 1000; i++ {
		fmt.Printf("%10d ", sfmt.genRand32())
		if i%5 == 4 {
			fmt.Println()
		}
	}
	fmt.Println()
}
