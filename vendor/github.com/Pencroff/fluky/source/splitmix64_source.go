package source

import "math/rand"

type SplitMix64Source struct {
	state uint64
}

func (s *SplitMix64Source) Seed(seed int64) {
	s.state = uint64(seed)
}

func (s *SplitMix64Source) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

func (s *SplitMix64Source) Uint64() uint64 {
	s.state += 0x9e3779b97f4a7c15
	z := s.state
	z ^= z >> 30
	z *= 0xbf58476d1ce4e5b9
	z ^= z >> 27
	z *= 0x94d049bb133111eb
	z ^= z >> 31
	return z
}

func NewSplitMix64Source(seed int64) rand.Source64 {
	r := &SplitMix64Source{}
	r.Seed(seed)
	return r
}
