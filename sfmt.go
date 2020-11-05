package mt

import "fmt"

type Params struct {
	MEXP   uint32
	POS1   uint32
	SL1    uint32
	SL2    uint32
	SR1    uint32
	SR2    uint32
	MSK1   uint32
	MSK2   uint32
	MSK3   uint32
	MSK4   uint32
	parity [4]uint32
	N      uint32
	N32    uint32
	N64    uint32
	IDSTR  string
	fSR2   uint32
	fSL2   uint32
	rSR2   uint32
	rSL2   uint32
	POS    uint32
}

func NewParams(mexp uint32, pos1 uint32, sl1 uint32, sl2 uint32, sr1 uint32, sr2 uint32,
	msk1 uint32, msk2 uint32, msk3 uint32, msk4 uint32, parity [4]uint32) Params {
	params := Params{
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
	params.N64 = params.N * 2
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

var params604, params1279, params2281, params4253, params11213, params19937, params44497, params86243, params132049, params216091 Params

func init() {
	params604 = NewParams(607,
		2, 15, 3, 13, 3,
		0xfdff37ff, 0xef7f3f7d, 0xff777b7d, 0x7ff7fb2f,
		[4]uint32{0x00000001, 0x00000000, 0x00000000, 0x5986f054})
	params1279 = NewParams(1279,
		7, 14, 3, 5, 1,
		0xf7fefffd, 0x7fefcfff, 0xaff3ef3f, 0xb5ffff7f,
		[4]uint32{0x00000001, 0x00000000, 0x00000000, 0x20000000})
	params2281 = NewParams(2281,
		12, 19, 1, 5, 1,
		0xbff7ffbf, 0xfdfffffe, 0xf7ffef7f, 0xf2f7cbbf,
		[4]uint32{0x00000001, 0x00000000, 0x00000000, 0x41dfa600})
	params4253 = NewParams(4253,
		17, 20, 1, 7, 1,
		0x9f7bffff, 0x9fffff5f, 0x3efffffb, 0xfffff7bb,
		[4]uint32{0xa8000001, 0xaf5390a3, 0xb740b3f8, 0x6c11486d})
	params11213 = NewParams(11213,
		68, 14, 3, 7, 3,
		0xeffff7fb, 0xffffffef, 0xdfdfbfff, 0x7fffdbfd,
		[4]uint32{0x00000001, 0x00000000, 0xe8148000, 0xd0c7afa3})
	params19937 = NewParams(19937,
		122, 18, 1, 11, 1,
		0xdfffffef, 0xddfecb7f, 0xbffaffff, 0xbffffff6,
		[4]uint32{0x00000001, 0x00000000, 0x00000000, 0x13c9e684})
	params44497 = NewParams(44497,
		330, 5, 3, 9, 3,
		0xeffffffb, 0xdfbebfff, 0xbfbf7bef, 0x9ffd7bff,
		[4]uint32{0x00000001, 0x00000000, 0xa3ac4000, 0xecc1327a})
	params86243 = NewParams(86243,
		366, 6, 7, 19, 1,
		0xfdbffbff, 0xbff7ff3f, 0xfd77efff, 0xbf9ff3ff,
		[4]uint32{0x00000001, 0x00000000, 0x00000000, 0xe9528d85})
	params132049 = NewParams(132049,
		110, 19, 1, 21, 1,
		0xffffbb5f, 0xfb6ebf95, 0xfffefffa, 0xcff77fff,
		[4]uint32{0x00000001, 0x00000000, 0xcb520000, 0xc7e91c7d})
	params216091 = NewParams(216091,
		627, 11, 3, 10, 1,
		0xbff7bff7, 0xbfffffff, 0xbffffa7f, 0xffddfbfb,
		[4]uint32{0xf8000001, 0x89e80709, 0x3bd2b64b, 0x0c64b1e4})
}

type SFMT struct {
	params Params
	idx    uint32
	state  []uint32
}

func int2long(u, l uint32) uint64 {
	return uint64(u)<<32 | uint64(l)
}

func NewSFMT(params Params) *SFMT {
	sfmt := &SFMT{
		params: params,
	}
	sfmt.params.fSR2 = 8 * sfmt.params.SR2
	sfmt.params.fSL2 = 8 * sfmt.params.SL2
	sfmt.params.rSR2 = 32 - sfmt.params.fSR2
	sfmt.params.rSL2 = 32 - sfmt.params.fSL2
	sfmt.params.POS = 4 * sfmt.params.POS1
	sfmt.state = make([]uint32, sfmt.params.N32, sfmt.params.N32)
	return sfmt
}

func (sfmt *SFMT) InitGenRand(seed uint32) {
	sfmt.state[0] = seed
	for i := uint32(1); i < sfmt.params.N32; i++ {
		sfmt.state[i] = 1812433253*(sfmt.state[i-1]^(sfmt.state[i-1]>>30)) + i
	}
	sfmt.idx = sfmt.params.N32
	sfmt.periodCertification()
}

func func1(x uint32) uint32 {
	return (x ^ (x >> 27)) * 1664525
}

func func2(x uint32) uint32 {
	return (x ^ (x >> 27)) * 1566083941
}

func (sfmt *SFMT) InitByArray(initKey []uint32) {
	params := sfmt.params
	state := sfmt.state
	keyLength := uint32(len(initKey))
	size := params.N * 4

	var lag uint32
	var count uint32
	var i, j, r uint32

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

	for i = 0; i < params.N32; i++ {
		state[i] = 0x8b8b8b8b
	}
	if keyLength+1 > params.N32 {
		count = keyLength + 1
	} else {
		count = params.N32
	}

	r = func1(state[0] ^ state[mid] ^ state[params.N32-1])
	state[mid] += r
	r += keyLength
	state[mid+lag] += r
	state[0] = r

	count--
	for i, j = 1, 0; j < count && j < keyLength; j++ {
		r = func1(state[i] ^ state[(i+mid)%params.N32] ^ state[(i+params.N32-1)%params.N32])
		state[(i+mid)%params.N32] += r
		r += initKey[j] + i
		state[(i+mid+lag)%params.N32] += r
		state[i] = r
		i = (i + 1) % params.N32
	}
	for ; j < count; j++ {
		r = func1(state[i] ^ state[(i+mid)%params.N32] ^ state[(i+params.N32-1)%params.N32])
		state[(i+mid)%params.N32] += r
		r += i
		state[(i+mid+lag)%params.N32] += r
		state[i] = r
		i = (i + 1) % params.N32
	}
	for j = 0; j < params.N32; j++ {
		r = func2(state[i] + state[(i+mid)%params.N32] + state[(i+params.N32-1)%params.N32])
		state[(i+mid)%params.N32] ^= r
		r -= i
		state[(i+mid+lag)%params.N32] ^= r
		state[i] = r
		i = (i + 1) % params.N32
	}
	sfmt.idx = params.N32
	sfmt.periodCertification()
}

func (sfmt *SFMT) periodCertification() {
	parity := sfmt.params.parity
	inner := uint32(0)
	for i := 0; i < 4; i++ {
		inner ^= sfmt.state[i] & parity[i]
	}
	for i := uint(16); i > 0; i >>= 1 {
		inner ^= inner >> i
	}
	inner &= 1
	if inner == 1 {
		return
	}
	for i := 0; i < 4; i++ {
		work := uint32(1)
		for j := 0; j < 32; j++ {
			if (work & parity[i]) != 0 {
				sfmt.state[i] ^= work
				return
			}
			work = work << 1
		}
	}
}

func (sfmt *SFMT) GenRandAll() {
	params := sfmt.params
	state := sfmt.state
	r1 := params.N32 - 8
	r2 := params.N32 - 4
	for i, pos := uint32(0), params.POS; i < params.N32; i, pos = i+4, pos+4 {
		if pos >= params.N32 {
			pos = 0
		}
		state[i+3] ^= ((state[i+3] << params.fSL2) | (state[i+2] >> params.rSL2)) ^
			((state[pos+3] >> params.SR1) & params.MSK4) ^
			(state[r1+3] >> params.fSR2) ^
			(state[r2+3] << params.SL1)
		state[i+2] ^= ((state[i+2] << params.fSL2) | (state[i+1] >> params.rSL2)) ^
			((state[pos+2] >> params.SR1) & params.MSK3) ^
			((state[r1+2] >> params.fSR2) | (state[r1+3] << params.rSR2)) ^
			(state[r2+2] << params.SL1)
		state[i+1] ^= ((state[i+1] << params.fSL2) | (state[i] >> params.rSL2)) ^
			((state[pos+1] >> params.SR1) & params.MSK2) ^
			((state[r1+1] >> params.fSR2) | (state[r1+2] << params.rSR2)) ^
			(state[r2+1] << params.SL1)
		state[i] ^= (state[i] << params.fSL2) ^
			((state[pos] >> params.SR1) & params.MSK1) ^
			((state[r1] >> params.fSR2) | (state[r1+1] << params.rSR2)) ^
			(state[r2] << params.SL1)
		r1, r2 = r2, i
	}
}

func (sfmt *SFMT) GenRand32() uint32 {
	if sfmt.idx >= sfmt.params.N32 {
		sfmt.GenRandAll()
		sfmt.idx = 0
	}
	result := sfmt.state[sfmt.idx]
	sfmt.idx++
	return result
}

func (sfmt *SFMT) GenRand64() uint64 {
	if sfmt.idx >= sfmt.params.N32 {
		sfmt.GenRandAll()
		sfmt.idx = 0
	}
	r1 := sfmt.state[sfmt.idx]
	r2 := sfmt.state[sfmt.idx+1]
	sfmt.idx += 2
	return int2long(r2, r1)
}
