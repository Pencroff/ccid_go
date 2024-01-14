package source

import (
	"math/bits"
	"math/rand"
)

type Xoshiro256ssSource struct {
	s0 uint64
	s1 uint64
	s2 uint64
	s3 uint64
}

func (x *Xoshiro256ssSource) Int63() int64 {
	return int64(x.Uint64() >> 1)
}

func (x *Xoshiro256ssSource) Seed(seed int64) {
	sm := NewSplitMix64Source(seed)
	x.s0 = sm.Uint64()
	x.s1 = sm.Uint64()
	x.s2 = sm.Uint64()
	x.s3 = sm.Uint64()
}

func (x *Xoshiro256ssSource) Uint64() uint64 {
	res := bits.RotateLeft64(x.s1*5, 7) * 9
	t := x.s1 << 17
	x.s2 ^= x.s0
	x.s3 ^= x.s1
	x.s1 ^= x.s2
	x.s0 ^= x.s3
	x.s2 ^= t
	x.s3 = bits.RotateLeft64(x.s3, 45)
	return res
}

func (x *Xoshiro256ssSource) Jump() rand.Source64 {
	s0, s1, s2, s3 := x.s0, x.s1, x.s2, x.s3
	r := &Xoshiro256ssSource{s0: s0, s1: s1, s2: s2, s3: s3}
	s0, s1, s2, s3 = 0, 0, 0, 0
	for _, jump := range jump128 {
		var b uint
		for b = 0; b < 64; b++ {
			if jump&(1<<b) != 0 {
				s0 ^= x.s0
				s1 ^= x.s1
				s2 ^= x.s2
				s3 ^= x.s3
			}
			x.Uint64()
		}
	}
	x.s0 = s0
	x.s1 = s1
	x.s2 = s2
	x.s3 = s3
	return r
}

func (x *Xoshiro256ssSource) LongJump() rand.Source64 {
	s0, s1, s2, s3 := x.s0, x.s1, x.s2, x.s3
	r := &Xoshiro256ssSource{s0, s1, s2, s3}
	s0, s1, s2, s3 = 0, 0, 0, 0
	for _, jump := range jump192 {
		var b uint
		for b = 0; b < 64; b++ {
			if jump&(1<<b) != 0 {
				s0 ^= x.s0
				s1 ^= x.s1
				s2 ^= x.s2
				s3 ^= x.s3
			}
			x.Uint64()
		}
	}
	x.s0 = s0
	x.s1 = s1
	x.s2 = s2
	x.s3 = s3
	return r
}

func NewXoshiro256ssSource(seed int64) rand.Source64 {
	r := &Xoshiro256ssSource{}
	r.Seed(seed)
	return r
}
