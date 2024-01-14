package source

import "math/rand"

// Based on https://www.youtube.com/watch?v=LWFzPP8ZbdU
// https://github.com/Descrout/squirrel3-rs
// https://github.com/archer884/squirrel-rng

type Squirrel3Source struct {
	position uint64
	seed     uint64
}

func (s *Squirrel3Source) Position() int64 {
	return int64(s.position)
}

func (s *Squirrel3Source) SetPosition(pos int64) {
	s.position = uint64(pos)
}

func (s *Squirrel3Source) Seed(seed int64) {
	s.position = 0
	s.seed = uint64(seed)
}

func (s *Squirrel3Source) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

func (s *Squirrel3Source) Uint64() uint64 {
	r := squirrel3(s.position, s.seed)
	s.position += 1
	return r
}

type Squirrel3x2Source struct {
	Squirrel3Source
}

func (s *Squirrel3x2Source) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

func (s *Squirrel3x2Source) Uint64() uint64 {
	const mask32 = 1<<32 - 1
	x := (squirrel3(s.position, s.seed) >> 16) & mask32 // used high 32 bits
	s.position += 1
	y := (squirrel3(s.position, s.seed) >> 16) & mask32 // used high 32 bits
	s.position += 1
	return x<<32 | y
}

type Squirrel3Prime64Source struct {
	Squirrel3Source
}

func (s *Squirrel3Prime64Source) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

func (s *Squirrel3Prime64Source) Uint64() uint64 {
	r := squirrel3uint64Noise(s.position, s.seed)
	s.position += 1
	return r
}

func squirrel3(pos, seed uint64) uint64 {
	/*
		const unsigned int BIT_NOISE1 = 0x6831DA4; // 0b0110 1000 1110'0011 0001 1101 1010'0100;
		const unsigned int BIT_NOISE2 = 0xB5297A4D; // 0b1011 ' 0101 ' 0010 ' 1001 ' 0111'1010' 0100'1101;
		const unsigned int BIT_NOISE3 - 0x1856C4E9; // 0b0001 1011'0101 0110'1100'0100'1110 1001;
		unsigned int mangledBits = (unsigned int) positionx;
		mangledBits *- BIT_NOISE1;
		mangledBits += seed;
		mangledits ^= (mangledBits » 8);
		mangledBits += BIT_NOISE2;
		mangledits ^= (mangledBits « 8);
		mangledBits *= BIT_NOISE3;
		mangledits ^= (mangledBits » 8);
		return mangledBits;
	*/
	const BIT_NOISE1 = 0x6831DA4  // 0b0110 1000 1110 0011 0001 1101 1010 0100
	const BIT_NOISE2 = 0xB5297A4D // 0b1011 0101 0010 1001 0111 1010 0100 1101
	const BIT_NOISE3 = 0x1856C4E9 // 0b0001 1011'0101 0110'1100'0100'1110 1001
	mangledBits := pos
	mangledBits *= BIT_NOISE1
	mangledBits += seed
	mangledBits ^= mangledBits >> 8
	mangledBits += BIT_NOISE2
	mangledBits ^= mangledBits << 8
	mangledBits *= BIT_NOISE3
	mangledBits ^= mangledBits >> 8
	return mangledBits
}

func squirrel3uint64Noise(pos, seed uint64) uint64 {
	const BIT_NOISE1 = 0x761a4cd0f83414b7
	const BIT_NOISE2 = 0xbb1dadd338b9dda9
	const BIT_NOISE3 = 0x8e71b6b573ab7755
	mangledBits := pos
	mangledBits *= BIT_NOISE1
	mangledBits += seed
	mangledBits ^= mangledBits >> 8
	mangledBits += BIT_NOISE2
	mangledBits ^= mangledBits << 8
	mangledBits *= BIT_NOISE3
	mangledBits ^= mangledBits >> 8
	return mangledBits
}

func NewSquirrel3Source(seed int64) rand.Source64 {
	r := &Squirrel3Source{}
	r.Seed(seed)
	return r
}

func NewSquirrel3x2Source(seed int64) rand.Source64 {
	r := &Squirrel3x2Source{}
	r.Seed(seed)
	return r
}

func NewSquirrel3Prime64Source(seed int64) rand.Source64 {
	r := &Squirrel3Prime64Source{}
	r.Seed(seed)
	return r
}
