package mt

import (
	"fmt"
	"testing"
)

func TestSFMT_604(t *testing.T) {
	sfmt := NewSFMT(params604)
	fmt.Printf("%s\n32 bit generated randoms\n", sfmt.params.IDSTR)
	sfmt.InitGenRand(1234)
	fmt.Println("init_gen_rand__________")
	for i := 0; i < 1000; i++ {
		fmt.Printf("%10d ", sfmt.GenRand32())
		if i%5 == 4 {
			fmt.Println()
		}
	}
	fmt.Println()

	sfmt.InitByArray([]uint32{0x1234, 0x5678, 0x9abc, 0xdef0})
	fmt.Println("init_by_array__________")
	for i := 0; i < 1000; i++ {
		fmt.Printf("%10d ", sfmt.GenRand32())
		if i%5 == 4 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func TestSFMT_1279(t *testing.T) {
	sfmt := NewSFMT(params1279)
	fmt.Printf("%s\n32 bit generated randoms\n", sfmt.params.IDSTR)
	sfmt.InitGenRand(1234)
	fmt.Println("init_gen_rand__________")
	for i := 0; i < 1000; i++ {
		fmt.Printf("%10d ", sfmt.GenRand32())
		if i%5 == 4 {
			fmt.Println()
		}
	}
	fmt.Println()

	sfmt.InitByArray([]uint32{0x1234, 0x5678, 0x9abc, 0xdef0})
	fmt.Println("init_by_array__________")
	for i := 0; i < 1000; i++ {
		fmt.Printf("%10d ", sfmt.GenRand32())
		if i%5 == 4 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func TestSFMT_19937(t *testing.T) {
	sfmt := NewSFMT(params19937)
	fmt.Printf("%s\n32 bit generated randoms\n", sfmt.params.IDSTR)
	sfmt.InitGenRand(1234)
	fmt.Println("init_gen_rand__________")
	for i := 0; i < 1000; i++ {
		fmt.Printf("%10d ", sfmt.GenRand32())
		if i%5 == 4 {
			fmt.Println()
		}
	}
	fmt.Println()

	sfmt.InitByArray([]uint32{0x1234, 0x5678, 0x9abc, 0xdef0})
	fmt.Println("init_by_array__________")
	for i := 0; i < 1000; i++ {
		fmt.Printf("%10d ", sfmt.GenRand32())
		if i%5 == 4 {
			fmt.Println()
		}
	}
	fmt.Println()
}
