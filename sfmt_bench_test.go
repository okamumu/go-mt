package mt

import (
	"testing"
)

func Benchmark1(b *testing.B) {
	sfmt := NewSFMT(params604)
	for i := 0; i < b.N; i++ {
		sfmt.InitGenRand(1234)
		for i := 0; i < 1000; i++ {
			sfmt.GenRand32()
		}
		sfmt.InitByArray([]uint32{0x1234, 0x5678, 0x9abc, 0xdef0})
		for i := 0; i < 1000; i++ {
			sfmt.GenRand32()
		}
	}
}

func Benchmark2(b *testing.B) {
	sfmt := NewSFMT(params604)
	for i := 0; i < b.N; i++ {
		sfmt.InitGenRand(1234)
		for i := 0; i < 1000; i++ {
			sfmt.GenRand64()
		}
		sfmt.InitByArray([]uint32{0x1234, 0x5678, 0x9abc, 0xdef0})
		for i := 0; i < 1000; i++ {
			sfmt.GenRand64()
		}
	}
}

func Benchmark3(b *testing.B) {
	sfmt := NewSFMT(params19937)
	for i := 0; i < b.N; i++ {
		sfmt.InitGenRand(1234)
		for i := 0; i < 1000; i++ {
			sfmt.GenRand64()
		}
		sfmt.InitByArray([]uint32{0x1234, 0x5678, 0x9abc, 0xdef0})
		for i := 0; i < 1000; i++ {
			sfmt.GenRand64()
		}
	}
}

func Benchmark4(b *testing.B) {
	sfmt := NewMT64()
	for i := 0; i < b.N; i++ {
		sfmt.InitGenRand(1234)
		for i := 0; i < 1000; i++ {
			sfmt.GenRand64()
		}
		sfmt.InitByArray([]uint64{0x1234, 0x5678, 0x9abc, 0xdef0})
		for i := 0; i < 1000; i++ {
			sfmt.GenRand64()
		}
	}
}
