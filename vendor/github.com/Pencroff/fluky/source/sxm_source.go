package source

import (
	"math/bits"
	"math/rand"
)

const (
	BIT_NOISE_A uint64 = 0xa55d924b5499ede9 // 0b1010010101011101100100100100101101010100100110011110110111101001
	BIT_NOISE_B uint64 = 0x531e8dd1c779c9d1 // 0b0101001100011110100011011101000111000111011110011100100111010001
	BIT_NOISE_C uint64 = 0xf6e7845905e176cb // 0b1111011011100111100001000101100100000101111000010111011011001011
)

// SxmSource is a random number generator based on Squirrel3 with more bits permutation (inspired by xxhash)
// Sxm build from Sums and XORs and Multiplications steps (SXM)
// It's a 64-bit generator with position and seed support
type SxmSource struct {
	position uint64
	seed     uint64
}

// Position returns current position
func (s *SxmSource) Position() int64 {
	return int64(s.position)
}

// SetPosition sets current position
func (s *SxmSource) SetPosition(pos int64) {
	s.position = uint64(pos)
}

// Seed sets seed value and resets position
func (s *SxmSource) Seed(seed int64) {
	s.position = 0
	s.seed = uint64(seed)
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (s *SxmSource) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// Uint64 returns a pseudo-random 64-bit value as a uint64.
func (s *SxmSource) Uint64() uint64 {
	r := sxmFunc(s.position, s.seed)
	s.position += 1
	return r
}

// SxmMixSource is a random number generator based on SxmSource with mix64 (like in xxhash) bits permutation in the end
type SxmMixSource struct {
	SxmSource
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (s *SxmMixSource) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// Uint64 returns a pseudo-random 64-bit value as a uint64.
func (s *SxmMixSource) Uint64() uint64 {
	return mix64(s.SxmSource.Uint64())
}

// NewSxmSource returns a new SxmSource that uses the given seed value.
func NewSxmSource(seed int64) rand.Source64 {
	r := &SxmSource{}
	r.Seed(seed)
	return r
}

// NewSxmMixSource returns a new SxmMixSource that uses the given seed value.
func NewSxmMixSource(seed int64) rand.Source64 {
	r := &SxmMixSource{}
	r.Seed(seed)
	return r
}

// sxmFunc is bits permutation function for SxmSource
func sxmFunc(pos, seed uint64) uint64 {
	v1 := pos + BIT_NOISE_A
	v2 := seed + BIT_NOISE_B
	v3 := pos - BIT_NOISE_C
	v4 := seed - BIT_NOISE_C
	r := bits.RotateLeft64(v1, 1) + bits.RotateLeft64(v2, 7) + bits.RotateLeft64(v3, 11) + bits.RotateLeft64(v4, 13)
	r *= BIT_NOISE_A
	r ^= bits.RotateLeft64(r, 37)
	r *= BIT_NOISE_B
	r ^= bits.RotateLeft64(r, 29)
	r *= BIT_NOISE_C
	r ^= bits.RotateLeft64(r, 19)
	return r
}

// sxmMixFunc is bits permutation function with mix64 round for SxmMixSource
func sxmMixFunc(pos, seed uint64) uint64 {
	return mix64(sxmFunc(pos, seed))
}

// mix64 is bits permutation function from xxhash
func mix64(h uint64) uint64 {
	h ^= h >> 33
	h *= BIT_NOISE_B
	h ^= h >> 29
	h *= BIT_NOISE_C
	h ^= h >> 32
	return h
}
