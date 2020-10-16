package gosfmt

import "fmt"

type Params struct {
	mexp    uint32
	pos1    uint32
	sl1     uint32
	sl2     uint32
	sr1     uint32
	sr2     uint32
	msk1    uint32
	msk2    uint32
	msk3    uint32
	msk4    uint32
	parity  [4]uint32
	n       uint32
	n32     uint32
	n64     uint32
	idstr   string
	fSR2    uint32
	fSL2    uint32
	rSR2    uint32
	rSL2    uint32
	pos     uint32
	msk1_64 uint64
	msk2_64 uint64
	msk3_64 uint64
}

func NewParams(mexp uint32, pos1 uint32, sl1 uint32, sl2 uint32, sr1 uint32, sr2 uint32,
	msk1 uint32, msk2 uint32, msk3 uint32, msk4 uint32, parity [4]uint32) Params {
	params := Params{
		mexp:   mexp,
		pos1:   pos1,
		sl1:    sl1,
		sl2:    sl2,
		sr1:    sr1,
		sr2:    sr2,
		msk1:   msk1,
		msk2:   msk2,
		msk3:   msk3,
		msk4:   msk4,
		parity: parity,
	}
	params.n = params.mexp/128 + 1
	params.n32 = params.n * 4
	params.n64 = params.n * 2
	params.idstr = fmt.Sprintf("SFMT-%d:%d-%d-%d-%d-%d:%08x-%08x-%08x-%08x",
		params.mexp,
		params.pos1,
		params.sl1,
		params.sl2,
		params.sr1,
		params.sr2,
		params.msk1,
		params.msk2,
		params.msk3,
		params.msk4)
	return params
}

var params604 Params

func init() {
	params604 = NewParams(607, 2, 15, 3, 13, 3,
		0xfdff37ff, 0xef7f3f7d, 0xff777b7d, 0x7ff7fb2f,
		[4]uint32{0x00000001, 0x00000000, 0x00000000, 0x5986f054})
}
