package gosfmt

type SFMT64 struct {
	params Params
	idx    uint32
	x      []uint64
}

func int2long(u, l uint32) uint64 {
	return uint64(u)<<32 | uint64(l)
}

func NewSFMT64(params Params) *SFMT64 {
	sfmt := &SFMT64{
		params: params,
	}
	sfmt.params.fSR2 = 8 * sfmt.params.sr2
	sfmt.params.fSL2 = 8 * sfmt.params.sl2
	sfmt.params.rSR2 = 64 - sfmt.params.fSR2
	sfmt.params.rSL2 = 64 - sfmt.params.fSL2
	sfmt.params.pos = 2 * sfmt.params.pos1

	var tmp uint32
	tmp = 0xffffffff
	tmp >>= sfmt.params.sr1
	sfmt.params.msk1_64 = int2long(sfmt.params.msk2&tmp, sfmt.params.msk1&tmp)
	sfmt.params.msk2_64 = int2long(sfmt.params.msk4&tmp, sfmt.params.msk3&tmp)
	tmp = 0xffffffff
	tmp >>= sfmt.params.sl1
	sfmt.params.msk3_64 = int2long(tmp, tmp)
	sfmt.x = make([]uint64, sfmt.params.n64, sfmt.params.n64)
	return sfmt
}

func (sfmt *SFMT64) initGenRand(seed uint32) {
	var tmpu, tmpl uint32
	tmpu = seed
	tmpl = 1812433253*(tmpu^(tmpu>>30)) + 1
	sfmt.x[0] = int2long(tmpl, tmpu)
	for i, k := uint32(2), 1; i < sfmt.params.n32; i, k = i+2, k+1 {
		tmpu = 1812433253*(tmpl^(tmpl>>30)) + i
		tmpl = 1812433252*(tmpu^(tmpu>>30)) + i + 1
		sfmt.x[k] = int2long(tmpl, tmpu)
	}
	sfmt.idx = sfmt.params.n64
	sfmt.periodCertification()
}

// func func1(x uint32) uint32 {
// 	return (x ^ (x >> 27)) * 1664525
// }

// func func2(x uint32) uint32 {
// 	return (x ^ (x >> 27)) * 1566083941
// }

// func (sfmt *SFMT32) initByArray(initKey []uint32) {
// 	params := sfmt.params
// 	x := sfmt.x
// 	keyLength := uint32(len(initKey))
// 	size := params.n * 4

// 	var lag uint32
// 	var count uint32
// 	var i, j, r uint32

// 	if size >= 623 {
// 		lag = 11
// 	} else if size >= 68 {
// 		lag = 7
// 	} else if size >= 39 {
// 		lag = 5
// 	} else {
// 		lag = 3
// 	}
// 	mid := (size - lag) / 2

// 	for i = 0; i < params.n32; i++ {
// 		x[i] = 0x8b8b8b8b
// 	}
// 	if keyLength+1 > params.n32 {
// 		count = keyLength + 1
// 	} else {
// 		count = params.n32
// 	}

// 	r = func1(x[0] ^ x[mid] ^ x[params.n32-1])
// 	x[mid] += r
// 	r += keyLength
// 	x[mid+lag] += r
// 	x[0] = r

// 	count--
// 	for i, j = 1, 0; j < count && j < keyLength; j++ {
// 		r = func1(x[i] ^ x[(i+mid)%params.n32] ^ x[(i+params.n32-1)%params.n32])
// 		x[(i+mid)%params.n32] += r
// 		r += initKey[j] + i
// 		x[(i+mid+lag)%params.n32] += r
// 		x[i] = r
// 		i = (i + 1) % params.n32
// 	}
// 	for ; j < count; j++ {
// 		r = func1(x[i] ^ x[(i+mid)%params.n32] ^ x[(i+params.n32-1)%params.n32])
// 		x[(i+mid)%params.n32] += r
// 		r += i
// 		x[(i+mid+lag)%params.n32] += r
// 		x[i] = r
// 		i = (i + 1) % params.n32
// 	}
// 	for j = 0; j < params.n32; j++ {
// 		r = func2(x[i] + x[(i+mid)%params.n32] + x[(i+params.n32-1)%params.n32])
// 		x[(i+mid)%params.n32] ^= r
// 		r -= i
// 		x[(i+mid+lag)%params.n32] ^= r
// 		x[i] = r
// 		i = (i + 1) % params.n32
// 	}
// 	sfmt.idx = params.n32
// 	sfmt.periodCertification()
// }

func (sfmt *SFMT64) periodCertification() {
	parity := []uint64{int2long(sfmt.params.parity[1], sfmt.params.parity[0]), int2long(sfmt.params.parity[3], sfmt.params.parity[2])}
	inner := uint32(0)
	for i := 0; i < 2; i++ {
		inner ^= uint32(sfmt.x[i])
		inner ^= uint32(sfmt.x[i] >> 32)
	}
	for i := uint(16); i > 0; i >>= 1 {
		inner ^= inner >> i
	}
	inner &= 1
	if inner == 1 {
		return
	}
	for i := 0; i < 2; i++ {
		work := uint64(1)
		for j := 0; j < 64; j++ {
			if (work & parity[i]) != 0 {
				sfmt.x[i] ^= work
				return
			}
			work = work << 1
		}
	}
}

func (sfmt *SFMT64) genRandAll() {
	params := sfmt.params
	x := sfmt.x
	r1 := params.n64 - 4
	r2 := params.n64 - 2
	for i, pos := uint32(0), params.pos; i < params.n64; i, pos = i+2, pos+2 {
		if pos >= params.n64 {
			pos = 0
		}
		x[i+1] ^= ((x[i+1] << params.fSL2) | (x[i] >> params.rSL2)) ^
			((x[pos+1] >> params.sr1) & params.msk2_64) ^
			(x[r1+1] >> params.fSR2) ^
			((x[r2+1] << params.sl1) & params.msk3_64)
		x[i] ^= (x[i] << params.fSL2) ^
			((x[pos] >> params.sr1) & params.msk1_64) ^
			((x[r1] >> params.fSR2) | (x[r1+1] << params.rSR2)) ^
			((x[r2] << params.sl1) & params.msk3_64)
		r1, r2 = r2, i
	}
}

func (sfmt *SFMT64) genRand64() uint64 {
	if sfmt.idx >= sfmt.params.n64 {
		sfmt.genRandAll()
		sfmt.idx = 0
	}
	result := sfmt.x[sfmt.idx]
	sfmt.idx++
	return result
}
