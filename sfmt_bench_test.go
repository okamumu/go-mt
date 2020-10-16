package gosfmt

import (
	"testing"
)

func Benchmark1(b *testing.B) {
	sfmt := NewSFMT32(params604)
	for i := 0; i < b.N; i++ {
		sfmt.initGenRand(1234)
		for i := 0; i < 1000; i++ {
			sfmt.genRand32()
		}
		sfmt.initByArray([]uint32{0x1234, 0x5678, 0x9abc, 0xdef0})
		for i := 0; i < 1000; i++ {
			sfmt.genRand32()
		}
	}
}

// func Benchmark2(b *testing.B) {
// 	sfmt := NewSFMT(params2_604)
// 	for i := 0; i < b.N; i++ {
// 		sfmt.InitGenRand(1234)
// 		for i := 0; i < 1000; i++ {
// 			sfmt.GenRand32()
// 		}
// 		sfmt.InitByArray([]uint32{0x1234, 0x5678, 0x9abc, 0xdef0})
// 		for i := 0; i < 1000; i++ {
// 			sfmt.GenRand32()
// 		}
// 	}
// }
