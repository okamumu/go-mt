package mt

type MTRandomIF interface {
	GenRand64() uint64
}

type MTRand struct {
	Source MTRandomIF
}

func NewMTRand64() MTRand {
	mt := NewMT64()
	return MTRand{Source: mt}
}

// The method to generate a non-negative 63bit integer
func (rng *MTRand) Int63() int64 {
	return int64(rng.mt.GenRand64() >> 1)
}

// The method to generate a non-negative 63bit integer
func (rng *MTRand) UInt64() uint64 {
	return rng.mt.GenRand64()
}

// The method to generate a float64 on [0, 1)
func (rng *MTRand) Float64() float64 {
	return float64(rng.mt.GenRand64()>>11) * (1.0 / 9007199254740992.0)
}

// The method to generate a float64 on [0, 1)
func (rng *MTRand) RealC0O1() float64 {
	return float64(rng.mt.GenRand64()>>11) * (1.0 / 9007199254740992.0)
}

// The method to generate a float64 on (0, 1)
func (rng *MTRand) RealO0O1() float64 {
	return (float64(rng.mt.GenRand64()>>12) + 0.5) * (1.0 / 4503599627370496.0)
}

// The method to generate a float64 on [0, 1]
func (rng *MTRand) RealC0C1() float64 {
	return float64(rng.mt.GenRand64()>>11) * (1.0 / 9007199254740991.0)
}
