package source

import (
	"math/bits"
	"math/rand"
)

type SmallPrngSource struct {
	a uint64
	b uint64
	c uint64
	d uint64
}

func (s *SmallPrngSource) Seed(seed int64) {
	s.a = 0xf1ea5eed
	s.b = uint64(seed)
	s.c = uint64(seed)
	s.d = uint64(seed)
	for i := 0; i < 20; i++ {
		s.Uint64()
	}
}
func (s *SmallPrngSource) Uint64() uint64 {
	e := s.a - bits.RotateLeft64(s.b, 7)
	s.a = s.b ^ bits.RotateLeft64(s.c, 13)
	s.b = s.c + bits.RotateLeft64(s.d, 37)
	s.c = s.d + e
	s.d = e + s.a
	return s.d
}
func (s *SmallPrngSource) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

func NewSmallPrngSource(seed int64) rand.Source64 {
	r := &SmallPrngSource{}
	r.Seed(seed)
	return r
}
