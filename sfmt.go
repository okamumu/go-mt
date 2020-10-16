package gosfmt

import "fmt"

type Params2 struct {
	MEXP   uint
	POS1   uint
	SL1    uint
	SL2    uint
	SR1    uint
	SR2    uint
	MSK1   uint32
	MSK2   uint32
	MSK3   uint32
	MSK4   uint32
	parity [4]uint32
	N      uint
	N32    uint
	IDSTR  string
}

func NewParams2(mexp uint, pos1 uint, sl1 uint, sl2 uint, sr1 uint, sr2 uint,
	msk1 uint32, msk2 uint32, msk3 uint32, msk4 uint32, parity [4]uint32) Params2 {
	params := Params2{
		MEXP:   mexp,
		POS1:   pos1,
		SL1:    sl1,
		SL2:    sl2,
		SR1:    sr1,
		SR2:    sr2,
		MSK1:   msk1,
		MSK2:   msk2,
		MSK3:   msk3,
		MSK4:   msk4,
		parity: parity,
	}
	params.N = params.MEXP/128 + 1
	params.N32 = params.N * 4
	params.IDSTR = fmt.Sprintf("SFMT-%d:%d-%d-%d-%d-%d:%08x-%08x-%08x-%08x",
		params.MEXP,
		params.POS1,
		params.SL1,
		params.SL2,
		params.SR1,
		params.SR2,
		params.MSK1,
		params.MSK2,
		params.MSK3,
		params.MSK4)
	return params
}

var params2_604 Params2

func init() {
	params2_604 = NewParams2(607, 2, 15, 3, 13, 3,
		0xfdff37ff, 0xef7f3f7d, 0xff777b7d, 0x7ff7fb2f,
		[4]uint32{0x00000001, 0x00000000, 0x00000000, 0x5986f054})
}

type w128 [4]uint32

type SFMT struct {
	params Params2
	idx    uint
	state  []w128
}

func NewSFMT(params Params2) *SFMT {
	sfmt := &SFMT{
		params: params,
	}
	sfmt.state = make([]w128, sfmt.params.N32, sfmt.params.N32)
	return sfmt
}

// func (sfmt *SFMT) GetState(i uint) uint32 {
// 	return sfmt.state[i>>2][i&3]
// }

// func (sfmt *SFMT) setState(i uint, x uint32) {
// 	sfmt.state[i>>2][i&3] = x
// }

func rshift128(out, in *w128, shift uint) {
	th := (uint64(in[3]) << 32) | (uint64(in[2]))
	tl := (uint64(in[1]) << 32) | (uint64(in[0]))

	oh := th >> (shift * 8)
	ol := tl >> (shift * 8)
	ol |= th << (64 - shift*8)
	out[1] = uint32(ol >> 32)
	out[0] = uint32(ol)
	out[3] = uint32(oh >> 32)
	out[2] = uint32(oh)
}

func lshift128(out, in *w128, shift uint) {
	th := (uint64(in[3]) << 32) | (uint64(in[2]))
	tl := (uint64(in[1]) << 32) | (uint64(in[0]))

	oh := th << (shift * 8)
	ol := tl << (shift * 8)
	oh |= tl >> (64 - shift*8)
	out[1] = uint32(ol >> 32)
	out[0] = uint32(ol)
	out[3] = uint32(oh >> 32)
	out[2] = uint32(oh)
}

func (sfmt *SFMT) do_recursion(r, a, b, c, d *w128) {
	params := sfmt.params
	var x, y w128
	lshift128(&x, a, params.SL2)
	rshift128(&y, c, params.SR2)
	r[0] = a[0] ^ x[0] ^ ((b[0] >> params.SR1) & params.MSK1) ^ y[0] ^ (d[0] << params.SL1)
	r[1] = a[1] ^ x[1] ^ ((b[1] >> params.SR1) & params.MSK2) ^ y[1] ^ (d[1] << params.SL1)
	r[2] = a[2] ^ x[2] ^ ((b[2] >> params.SR1) & params.MSK3) ^ y[2] ^ (d[2] << params.SL1)
	r[3] = a[3] ^ x[3] ^ ((b[3] >> params.SR1) & params.MSK4) ^ y[3] ^ (d[3] << params.SL1)
}

func (sfmt *SFMT) GenRand32() uint32 {
	if sfmt.idx >= sfmt.params.N32 {
		sfmt.GenRandAll()
		sfmt.idx = 0
	}
	r := sfmt.state[sfmt.idx>>2][sfmt.idx&3]
	sfmt.idx++
	return r
}

func (sfmt *SFMT) GenRand64() uint64 {
	if sfmt.idx >= sfmt.params.N32 {
		sfmt.GenRandAll()
		sfmt.idx = 0
	}
	r1 := sfmt.state[sfmt.idx>>2][sfmt.idx&3]
	sfmt.idx++
	r2 := sfmt.state[sfmt.idx>>2][sfmt.idx&3]
	sfmt.idx++
	return (uint64(r2) << 32) | uint64(r1)
}

func (sfmt *SFMT) GenRandAll() {
	params := sfmt.params
	state := sfmt.state
	r1 := &state[params.N-2]
	r2 := &state[params.N-1]
	var i uint
	for i = 0; i < params.N-params.POS1; i++ {
		sfmt.do_recursion(&state[i], &state[i], &state[i+params.POS1], r1, r2)
		r1, r2 = r2, &state[i]
	}
	for ; i < params.N; i++ {
		sfmt.do_recursion(&state[i], &state[i], &state[i+params.POS1-params.N], r1, r2)
		r1, r2 = r2, &state[i]
	}
}

func (sfmt *SFMT) InitGenRand(seed uint32) {
	sfmt.state[0][0] = seed
	for i := uint(1); i < sfmt.params.N32; i++ {
		sfmt.state[i>>2][i&3] = 1812433253*(sfmt.state[(i-1)>>2][(i-1)&3]^(sfmt.state[(i-1)>>2][(i-1)&3]>>30)) + uint32(i)
	}
	sfmt.idx = sfmt.params.N32
	sfmt.periodCertification()
}

// func func1(x uint32) uint32 {
// 	return (x ^ (x >> 27)) * 1664525
// }

// func func2(x uint32) uint32 {
// 	return (x ^ (x >> 27)) * 1566083941
// }

func (sfmt *SFMT) InitByArray(initKey []uint32) {
	params := sfmt.params
	keyLength := uint(len(initKey))
	size := params.N * 4

	var lag, count, i, j uint
	var r uint32

	if size >= 623 {
		lag = 11
	} else if size >= 68 {
		lag = 7
	} else if size >= 39 {
		lag = 5
	} else {
		lag = 3
	}
	mid := (size - lag) / 2

	for i := 0; i < len(sfmt.state); i++ {
		sfmt.state[i][0] = 0x8b8b8b8b
		sfmt.state[i][1] = 0x8b8b8b8b
		sfmt.state[i][2] = 0x8b8b8b8b
		sfmt.state[i][3] = 0x8b8b8b8b
	}
	if keyLength+1 > params.N32 {
		count = keyLength + 1
	} else {
		count = params.N32
	}

	r = func1(sfmt.state[0][0] ^ sfmt.state[mid>>2][mid&3] ^ sfmt.state[(params.N32-1)>>2][(params.N32-1)&3])
	sfmt.state[mid>>2][mid&3] += r
	r += uint32(keyLength)
	sfmt.state[(mid+lag)>>2][(mid+lag)&3] += r
	sfmt.state[0][0] = r

	count--
	for i, j = 1, 0; j < count && j < keyLength; j++ {
		r = func1(sfmt.state[i>>2][i&3] ^ sfmt.state[(i+mid)>>2][(i+mid)&3] ^
			sfmt.state[((i+params.N32-1)%params.N32)>>2][((i+params.N32-1)%params.N32)&3])
		sfmt.state[(i+mid)%params.N32>>2][(i+mid)%params.N32%3] += r
		r += initKey[j] + uint32(i)
		sfmt.state[(i+mid+lag)%params.N32>>2][(i+mid+lag)%params.N32&3] += r
		sfmt.state[i>>2][i&3] = r
		i = (i + 1) % params.N32
	}
	for ; j < count; j++ {
		r = func1(sfmt.state[i>>2][i&3] ^ sfmt.state[(i+mid)%params.N32>>2][(i+mid)%params.N32&3] ^
			sfmt.state[((i+params.N32-1)%params.N32)>>2][((i+params.N32-1)%params.N32)&3])
		sfmt.state[(i+mid)%params.N32>>2][(i+mid)%params.N32%3] += r
		r += uint32(i)
		sfmt.state[(i+mid+lag)%params.N32>>2][(i+mid+lag)%params.N32&3] += r
		sfmt.state[i>>2][i&3] = r
		i = (i + 1) % params.N32
	}
	for j = 0; j < params.N32; j++ {
		r = func1(sfmt.state[i>>2][i&3] + sfmt.state[(i+mid)%params.N32>>2][(i+mid)%params.N32&3] +
			sfmt.state[((i+params.N32-1)%params.N32)>>2][((i+params.N32-1)%params.N32)&3])
		sfmt.state[(i+mid)%params.N32>>2][(i+mid)%params.N32%3] ^= r
		r -= uint32(i)
		sfmt.state[(i+mid+lag)%params.N32>>2][(i+mid+lag)%params.N32&3] ^= r
		sfmt.state[i>>2][i&3] = r
		i = (i + 1) % params.N32
	}
	sfmt.idx = params.N32
	sfmt.periodCertification()
}

func (sfmt *SFMT) periodCertification() {
	parity := sfmt.params.parity
	inner := uint32(0)
	for i := uint(0); i < 4; i++ {
		inner ^= sfmt.state[i>>2][i&3] & parity[i]
	}
	for i := uint(16); i > 0; i >>= 1 {
		inner ^= inner >> i
	}
	inner &= 1
	if inner == 1 {
		return
	}
	for i := uint(0); i < 4; i++ {
		work := uint32(1)
		for j := 0; j < 32; j++ {
			if (work & parity[i]) != 0 {
				sfmt.state[i>>2][i&3] ^= work
				return
			}
			work = work << 1
		}
	}
}
