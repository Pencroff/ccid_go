package ccid_go

import (
	"bytes"
	"fmt"
	e "github.com/Pencroff/ccid_go/extras"
	p "github.com/Pencroff/ccid_go/pkg"
	"sync"
	"testing"
)

// SplitMix64 is a 64-bit pseudo-random number generator.
// It is a fast, high-quality generator. For generating just floating-point
// numbers, SplitMix64 is even faster than xoroshiro128+.
// The state must be seeded so that it is not everywhere zero. If you have
// a 64-bit seed, we suggest to seed a SplitMix64 generator and use its
// output to fill s.
type SplitMix64 struct {
	state uint64
}

// NewSplitMix64 returns a new SplitMix64 generator to generate
// non-negative pseudo-random 64-bit integers.
func NewSplitMix64(seed int64) *SplitMix64 {
	rng := &SplitMix64{}
	rng.Seed(seed)
	return rng
}

// Uint64 returns a pseudo-random 64-bit integer as a uint64.
func (s *SplitMix64) Uint64() uint64 {
	s.state += 0x9e3779b97f4a7c15
	z := s.state
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

// Int63 returns a non-negative pseudo-random 63-bit integer as an int64.
func (s *SplitMix64) Int63() int64 {
	return int64(s.Uint64() >> 1)
}

// Seed uses the provided seed value to initialize the generator to a
// deterministic state.
func (s *SplitMix64) Seed(seed int64) {
	s.state = uint64(seed)
	s.state = s.Uint64()
}

func TestCcIdGen_Race(t *testing.T) {

	const (
		numRoutines = 10
		numCycles   = 100
	)
	sizes := []byte{
		p.ByteSliceSize64,
		p.ByteSliceSize96,
		p.ByteSliceSize128,
		p.ByteSliceSize160,
	}
	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			r, _ := e.NewHybridRandReader(NewSplitMix64(123))
			s := p.NewFiftyPercentMonotonicStrategy(r)
			notLockedGen, _ := NewMonotonicCcIdGen(size, r, s)
			gen := NewCcIdGenLocked(notLockedGen)
			prev, _ := gen.Next()
			queue := make(chan p.CcId)
			var wg sync.WaitGroup
			defer wg.Wait()
			wg.Add(numRoutines)
			for i := 0; i < numRoutines; i++ {
				go func(i int, g CcIdGen, q chan p.CcId) {
					defer wg.Done()
					for j := 0; j < numCycles; j++ {
						id, _ := g.Next()
						q <- id
					}
				}(i, gen, queue)
			}

			for i := 0; i < numRoutines*numCycles; i++ {
				next := <-queue
				if bytes.Compare(prev.Bytes(), next.Bytes()) >= 0 {
					fmt.Printf("%d - Not sequential in channel.\nPrev: %x\nmore then\nNext: %x\nDetails:\nPrev: %p - %#v\nNext: %p - %#v\n",
						i, prev.Bytes(), next.Bytes(), &prev, prev, &next, next)
				}
				prev = next
			}
		})
	}

}
