package gosfmt

type SFMT32 struct {
	params Params
	idx    uint32
	x      []uint32
}

func int2long(u, l uint32) uint64 {
	return uint64(u)<<32 | uint64(l)
}

func NewSFMT32(params Params) *SFMT32 {
	sfmt := &SFMT32{
		params: params,
	}
	sfmt.params.fSR2 = 8 * sfmt.params.sr2
	sfmt.params.fSL2 = 8 * sfmt.params.sl2
	sfmt.params.rSR2 = 32 - sfmt.params.fSR2
	sfmt.params.rSL2 = 32 - sfmt.params.fSL2
	sfmt.params.pos = 4 * sfmt.params.pos1
	sfmt.x = make([]uint32, sfmt.params.n32, sfmt.params.n32)
	return sfmt
}

func (sfmt *SFMT32) initGenRand(seed uint32) {
	sfmt.x[0] = seed
	for i := uint32(1); i < sfmt.params.n32; i++ {
		sfmt.x[i] = 1812433253*(sfmt.x[i-1]^(sfmt.x[i-1]>>30)) + i
	}
	sfmt.idx = sfmt.params.n32
	sfmt.periodCertification()
}

func func1(x uint32) uint32 {
	return (x ^ (x >> 27)) * 1664525
}

func func2(x uint32) uint32 {
	return (x ^ (x >> 27)) * 1566083941
}

func (sfmt *SFMT32) initByArray(initKey []uint32) {
	params := sfmt.params
	x := sfmt.x
	keyLength := uint32(len(initKey))
	size := params.n * 4

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

	for i = 0; i < params.n32; i++ {
		x[i] = 0x8b8b8b8b
	}
	if keyLength+1 > params.n32 {
		count = keyLength + 1
	} else {
		count = params.n32
	}

	r = func1(x[0] ^ x[mid] ^ x[params.n32-1])
	x[mid] += r
	r += keyLength
	x[mid+lag] += r
	x[0] = r

	count--
	for i, j = 1, 0; j < count && j < keyLength; j++ {
		r = func1(x[i] ^ x[(i+mid)%params.n32] ^ x[(i+params.n32-1)%params.n32])
		x[(i+mid)%params.n32] += r
		r += initKey[j] + i
		x[(i+mid+lag)%params.n32] += r
		x[i] = r
		i = (i + 1) % params.n32
	}
	for ; j < count; j++ {
		r = func1(x[i] ^ x[(i+mid)%params.n32] ^ x[(i+params.n32-1)%params.n32])
		x[(i+mid)%params.n32] += r
		r += i
		x[(i+mid+lag)%params.n32] += r
		x[i] = r
		i = (i + 1) % params.n32
	}
	for j = 0; j < params.n32; j++ {
		r = func2(x[i] + x[(i+mid)%params.n32] + x[(i+params.n32-1)%params.n32])
		x[(i+mid)%params.n32] ^= r
		r -= i
		x[(i+mid+lag)%params.n32] ^= r
		x[i] = r
		i = (i + 1) % params.n32
	}
	sfmt.idx = params.n32
	sfmt.periodCertification()
}

func (sfmt *SFMT32) periodCertification() {
	parity := sfmt.params.parity
	inner := uint32(0)
	for i := 0; i < 4; i++ {
		inner ^= sfmt.x[i] & parity[i]
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
				sfmt.x[i] ^= work
				return
			}
			work = work << 1
		}
	}
}

func (sfmt *SFMT32) genRandAll() {
	params := sfmt.params
	x := sfmt.x
	r1 := params.n32 - 8
	r2 := params.n32 - 4
	for i, pos := uint32(0), params.pos; i < params.n32; i, pos = i+4, pos+4 {
		if pos >= params.n32 {
			pos = 0
		}
		x[i+3] ^= ((x[i+3] << params.fSL2) | (x[i+2] >> params.rSL2)) ^
			((x[pos+3] >> params.sr1) & params.msk4) ^
			(x[r1+3] >> params.fSR2) ^
			(x[r2+3] << params.sl1)
		x[i+2] ^= ((x[i+2] << params.fSL2) | (x[i+1] >> params.rSL2)) ^
			((x[pos+2] >> params.sr1) & params.msk3) ^
			((x[r1+2] >> params.fSR2) | (x[r1+3] << params.rSR2)) ^
			(x[r2+2] << params.sl1)
		x[i+1] ^= ((x[i+1] << params.fSL2) | (x[i] >> params.rSL2)) ^
			((x[pos+1] >> params.sr1) & params.msk2) ^
			((x[r1+1] >> params.fSR2) | (x[r1+2] << params.rSR2)) ^
			(x[r2+1] << params.sl1)
		x[i] ^= (x[i] << params.fSL2) ^
			((x[pos] >> params.sr1) & params.msk1) ^
			((x[r1] >> params.fSR2) | (x[r1+1] << params.rSR2)) ^
			(x[r2] << params.sl1)
		r1, r2 = r2, i
	}
}

func (sfmt *SFMT32) genRand32() uint32 {
	if sfmt.idx >= sfmt.params.n32 {
		sfmt.genRandAll()
		sfmt.idx = 0
	}
	result := sfmt.x[sfmt.idx]
	sfmt.idx++
	return result
}

func (sfmt *SFMT32) genRand64() uint64 {
	if sfmt.idx >= sfmt.params.n32 {
		sfmt.genRandAll()
		sfmt.idx = 0
	}
	r1 := sfmt.x[sfmt.idx]
	r2 := sfmt.x[sfmt.idx+1]
	sfmt.idx += 2
	return int2long(r2, r1)
}
