package pkg

import (
	"io"
)

// CcIdMonotonicStrategy is an interface for monotonic strategy
// It should return mutated payload and carry byte
type CcIdMonotonicStrategy interface {
	Mutate(v []byte) (res []byte, carry byte)
}

// IncreaseMonotonicStrategy is a monotonic strategy that increase payload by 1
type IncreaseMonotonicStrategy struct{}

// Mutate implements CcIdMonotonicStrategy interface
// It increases payload by 1
func (s *IncreaseMonotonicStrategy) Mutate(v []byte) (res []byte, carry byte) {
	res, carry = Add8BigEndian(v, []byte{1}, 0)
	return
}

// NewIncreaseMonotonicStrategy creates a new IncreaseMonotonicStrategy
func NewIncreaseMonotonicStrategy() CcIdMonotonicStrategy {
	return &IncreaseMonotonicStrategy{}
}

// FiftyPercentMonotonicStrategy is a monotonic strategy that increase payload by random value between 0 and 50% of payload bits size
// For example, if payload is 32 bits, it will increase payload by random value of 16 bits
type FiftyPercentMonotonicStrategy struct {
	rd io.Reader
}

// Mutate implements CcIdMonotonicStrategy interface
// It increases payload by random value between 0 and 50% of payload bits size
// If no data read, it will return 1 as carry to trigger next time tick
func (s *FiftyPercentMonotonicStrategy) Mutate(v []byte) (res []byte, carry byte) {
	l := len(v)
	isOdd := l%2 == 1
	size := l >> 1
	if isOdd {
		size += 1
	}
	rndData := make([]byte, size)
fillWithRandom:
	_, err := s.rd.Read(rndData)
	// if no data read, return 1 as carry to trigger next time tick
	if err != nil {
		rndData = nil
		carry = 1
		return
	}
	if isOdd {
		rndData[0] &= 0x0F
	}
	if isZeroFilled(rndData) {
		goto fillWithRandom
	}
	res, carry = Add8BigEndian(v, rndData, 0)
	return
}

// NewFiftyPercentMonotonicStrategy creates a new FiftyPercentMonotonicStrategy
// It will use the given io.Reader to generate random bytes
func NewFiftyPercentMonotonicStrategy(rd io.Reader) CcIdMonotonicStrategy {
	return &FiftyPercentMonotonicStrategy{rd}
}

func isZeroFilled(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return false
		}
	}
	return true
}
