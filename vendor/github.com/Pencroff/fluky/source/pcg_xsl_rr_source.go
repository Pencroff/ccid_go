package source

import (
	et "github.com/Pencroff/go-toolkit/extended_type"
	"math/bits"
	"math/rand"
)

// PCG-XSL-RR 128 bit state, 64 bit output
// https://www.pcg-random.org/pdf/hmc-cs-2014-0905.pdf - 6.3.3 PCG-XSL-RR

type PcgXslRrSource struct {
	state et.Uint128
	inc   et.Uint128
	mul   et.Uint128
}

func (s *PcgXslRrSource) Seed(seed int64) {
	initseq := uint64(54)
	s.state = et.ZeroUint128
	s.inc = et.From64(initseq).Lsh(1).Or64(1)
	s.step()
	s.state = s.state.Add64(uint64(seed))
	s.step()
}

func (s *PcgXslRrSource) Uint64() uint64 {
	s.step()
	return s.stateToValue()
}

func (s *PcgXslRrSource) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

func (s *PcgXslRrSource) step() {
	s.state = s.state.Mul(s.mul).Add(s.inc)
}

func (s *PcgXslRrSource) stateToValue() (v uint64) {
	v = s.state.Lo ^ s.state.Hi
	rot := int(s.state.Rsh(122).Lo&0x3F) * -1
	v = bits.RotateLeft64(v, rot)
	return
}

func NewPcgXslRrSource(seed int64) rand.Source64 {
	r := &PcgXslRrSource{
		state: et.Uint128{},
		inc:   et.Uint128{},
		mul:   et.New(0x4385df649fccf645, 0x2360ed051fc65da4),
	}
	r.Seed(seed)
	return r
}
