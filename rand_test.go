package mt

import (
	"fmt"
	"testing"
)

func TestRandIF1(t *testing.T) {
	mt := NewMT64()
	mt.InitByArray([]uint64{0x12345, 0x23456, 0x34567, 0x45678})
	rng := MTRand{mt: mt}
	for i := 0; i < 1000; i++ {
		fmt.Println(rng.Float64())
	}
}
