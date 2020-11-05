package mt

const (
	nn       int    = 312
	mm       int    = 156
	matrix_A uint64 = 0xB5026F5AA96619E9
	um       uint64 = 0xffffffff80000000
	lm       uint64 = 0x000000007fffffff
)

var mag01 [2]uint64

func init() {
	mag01 = [2]uint64{0, matrix_A}
}

type MT64 struct {
	mt  []uint64
	mti int
}

func NewMT64() *MT64 {
	return &MT64{
		mt:  make([]uint64, nn, nn),
		mti: nn + 1,
	}
}

func (m *MT64) InitGenRand(seed uint64) {
	m.mt[0] = seed
	for m.mti = 1; m.mti < nn; m.mti++ {
		m.mt[m.mti] = (6364136223846793005*(m.mt[m.mti-1]^(m.mt[m.mti-1]>>62)) + uint64(m.mti))
	}
}

func (m *MT64) Seed(seed int64) {
	m.InitGenRand(uint64(seed))
}

func (m *MT64) InitByArray(initKey []uint64) {
	mt := m.mt
	keyLength := len(initKey)
	m.InitGenRand(19650218)
	i, j := 1, 0
	var k int
	if nn > keyLength {
		k = nn
	} else {
		k = keyLength
	}

	for ; k > 0; k-- {
		mt[i] = (mt[i] ^ ((mt[i-1] ^ (mt[i-1] >> 62)) * 3935559000370003845)) + initKey[j] + uint64(j)
		i++
		j++
		if i >= nn {
			mt[0] = mt[nn-1]
			i = 1
		}
		if j >= keyLength {
			j = 0
		}
	}
	for k = nn - 1; k > 0; k-- {
		mt[i] = (mt[i] ^ ((mt[i-1] ^ (mt[i-1] >> 62)) * 2862933555777941757)) - uint64(i)
		i++
		if i >= nn {
			mt[0] = mt[nn-1]
			i = 1
		}
	}
	mt[0] = 1 << 63
}

func (m *MT64) GenRand64() uint64 {
	if m.mti >= nn {
		mt := m.mt
		if m.mti == nn+1 {
			m.InitGenRand(5489)
		}
		var i int
		for i = 0; i < nn-mm; i++ {
			x := (mt[i] & um) | (mt[i+1] & lm)
			mt[i] = mt[i+mm] ^ (x >> 1) ^ mag01[x&1]
		}
		for ; i < nn-1; i++ {
			x := (mt[i] & um) | (mt[i+1] & lm)
			mt[i] = mt[i+(mm-nn)] ^ (x >> 1) ^ mag01[x&1]
		}
		x := (mt[nn-1] & um) | (mt[0] & lm)
		mt[nn-1] = mt[mm-1] ^ (x >> 1) ^ mag01[x&1]
		m.mti = 0
	}
	x := m.mt[m.mti]
	m.mti++

	x ^= (x >> 29) & 0x5555555555555555
	x ^= (x << 17) & 0x71D67FFFEDA60000
	x ^= (x << 37) & 0xFFF7EEE000000000
	x ^= (x >> 43)
	return x
}
