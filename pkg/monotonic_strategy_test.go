package pkg

import (
	"bytes"
	"testing"
)

type MonotonicStrategyTestCase struct {
	in    []byte
	out   []byte
	carry byte
}

func TestIncreaseMonotonicStrategy_Mutate(t *testing.T) {
	var testCaseIncreaseMonotonicStrategy = map[string]MonotonicStrategyTestCase{
		"regular": {
			[]byte{1, 2, 3},
			[]byte{1, 2, 4},
			0,
		},
		"first byte upper bound": {
			[]byte{1, 2, 254},
			[]byte{1, 2, 255},
			0,
		},
		"first byte overflow": {
			[]byte{1, 2, 255},
			[]byte{1, 3, 0},
			0,
		},
		"overflow": {
			[]byte{255, 255, 255, 255},
			[]byte{0, 0, 0, 0},
			1,
		},
	}
	s := NewIncreaseMonotonicStrategy()
	keys := SortKeys(testCaseIncreaseMonotonicStrategy)
	for _, k := range keys {
		t.Run(k, func(t *testing.T) {
			tc := testCaseIncreaseMonotonicStrategy[k]
			output, carryOut := s.Mutate(tc.in)
			if !bytes.Equal(output, tc.out) || carryOut != tc.carry {
				t.Errorf("s.Mutate(%v) =\n%v, %d; want\n%v, %d", tc.in, output, carryOut, tc.out, tc.carry)
			}
		})
	}
}

// Unexported
type mockReader struct {
	Val byte
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = m.Val
		m.Val += 1
	}
	return len(p), nil
}

type mockReaderSlice struct {
	Val []byte
	idx int
}

func (m *mockReaderSlice) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = m.Val[m.idx]
		m.idx += 1
		if m.idx >= len(m.Val) {
			m.idx = 0
		}
	}
	return len(p), nil
}

func TestFiftyPercentMonotonicStrategy_Mutate(t *testing.T) {
	var testCaseFiftyPercentMonotonicStrategy = map[string]MonotonicStrategyTestCase{
		"regular odd": {
			[]byte{1, 2, 3, 4, 5}, // 0x05, 0xA6, 0xA7 - 50% of 5 bytes
			[]byte{0x01, 0x02, 0x08, 0xAA, 0xAC},
			0,
		},
		"regular even": {
			[]byte{1, 2, 3, 4, 5, 6}, // 0xA5, 0xA6, 0xA7 - 50% of 6 bytes
			[]byte{0x01, 0x02, 0x03, 0xA9, 0xAB, 0xAD},
			0,
		},
		"overflow odd": {
			[]byte{0xFF, 0xFF, 0xFD, 0xFC}, // 0xA5, 0xA6 - 50% of 4 bytes
			[]byte{0x00, 0x00, 0xA3, 0xA2},
			1,
		},
		"overflow even": {
			[]byte{0xFF, 0xFF, 0xFE, 0xFD, 0xFC}, // 0x05, 0xA6, 0xA7 - 50% of 5 bytes
			[]byte{0x00, 0x00, 0x04, 0xA4, 0xA3},
			1,
		},
	}

	keys := SortKeys(testCaseFiftyPercentMonotonicStrategy)
	for _, k := range keys {
		rd := mockReader{0xA5}
		s := NewFiftyPercentMonotonicStrategy(&rd)
		t.Run(k, func(t *testing.T) {
			tc := testCaseFiftyPercentMonotonicStrategy[k]
			output, carryOut := s.Mutate(tc.in)
			if !bytes.Equal(output, tc.out) || carryOut != tc.carry {
				t.Errorf("s.Mutate(%v) =\n%v, %d; want\n%v, %d", tc.in, output, carryOut, tc.out, tc.carry)
			}
		})
	}
}

func TestFiftyPercentMonotonicStrategy_MutateWithZeroSequence(t *testing.T) {
	in := []byte{0x01, 0x01, 0x01, 0x01}
	out := []byte{0x01, 0x01, 0x02, 0x02}
	rd := mockReaderSlice{[]byte{0x00, 0x00, 0x00, 0x00, 0x01, 0x01, 0x01, 0x01}, 0}
	s := NewFiftyPercentMonotonicStrategy(&rd)
	output, carryOut := s.Mutate(in)
	if !bytes.Equal(output, out) || carryOut != 0 {
		t.Errorf("s.Mutate(%v) =\n%v, %d; want\n%v, %d", in, output, carryOut, out, 0)
	}
}

func TestFiftyPercentMonotonicStrategy_MutateWithZeroSequenceOdd(t *testing.T) {
	in := []byte{0x01, 0x01, 0x01}
	out := []byte{0x01, 0x02, 0x02}
	rd := mockReaderSlice{[]byte{0x10, 0x00, 0x01, 0x01, 0x01, 0x01}, 0}
	s := NewFiftyPercentMonotonicStrategy(&rd)
	output, carryOut := s.Mutate(in)
	if !bytes.Equal(output, out) || carryOut != 0 {
		t.Errorf("s.Mutate(%v) =\n%v, %d; want\n%v, %d", in, output, carryOut, out, 0)
	}
}
