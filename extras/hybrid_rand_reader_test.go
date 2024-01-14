package extras

import (
	"encoding/binary"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type RandSource64Mock struct {
	mock.Mock
	counter uint8
}

func (m *RandSource64Mock) Seed(seed int64) {
	m.Called(seed)
	m.counter = 0
}

func (m *RandSource64Mock) Uint64() uint64 {
	args := m.Mock.Called()
	return args.Get(0).(uint64)
	// uint64(m.counter) * 0x0101010101010101
}

func (m *RandSource64Mock) Int63() int64 {
	return int64(m.Uint64() >> 1)
}

func TestMockSource64(t *testing.T) {
	src := new(RandSource64Mock)
	var values = []uint64{
		0x0101010101010101, 0x0202020202020202,
		0x0303030303030303, 0x0404040404040404,
		0x0505050505050505}
	src.On("Seed", mock.Anything).Return()
	for i := 0; i < 5; i++ {
		src.On("Uint64").Return(values[i]).Once()
	}
	src.Seed(time.Now().Unix())
	for i := 0; i < 5; i++ {
		v := src.Uint64()
		if v != values[i] {
			t.Errorf("Uint64() = %016x, want %016x", v, values[i])
			return
		}
	}
	src.AssertExpectations(t)
}

func TestNewHybridRandReader(t *testing.T) {
	totalMockValues := 3
	src := new(RandSource64Mock)
	for nn := 0; nn < totalMockValues; nn++ {
		src.On("Uint64").Return(uint64(nn+1) * 0x0101010101010101).Once()
	}
	r, err := NewHybridRandReaderWithSizeAndSource(32, src)
	if err != nil {
		t.Errorf("NewHybridRandReaderWithSizeAndSource() error = %v", err)
		return
	}
	size := 16
	b := make([]byte, size)
	n, err := r.Read(b)
	if err != nil {
		t.Errorf("Read() error = %v", err)
		return
	}
	if n != size {
		t.Errorf("Read() n = %v, want %v", n, 5)
		return
	}
	v := byte(1)
	for i := 0; i < size; i++ {
		if i != 0 && i%7 == 0 {
			v += 1
		}
		if b[i] != v {
			t.Errorf("Read() b[%v] = %v, want %v", i, b[i], v)
			return
		}
	}
	src.AssertExpectations(t)
}

func TestNewHybridRandReader_SeedAfterLimit(t *testing.T) {
	shieldSize := uint64(16)
	var values = []uint64{
		0x0807060504030201, 0x100f0e0d0c0b0a09,
		0x2827262524232221, 0x302f2e2d2c2b2a29,
		0x4847464544434241}
	var res = []struct {
		idx, pos uint8
	}{
		{0, 1}, {0, 2}, {0, 3}, {0, 4},
		{0, 5}, {0, 6}, {0, 7}, {1, 1},
		{1, 2}, {1, 3}, {1, 4}, {1, 5},
		{1, 6}, {1, 7}, {2, 1}, {2, 2},
		{2, 3}, {2, 4}, {2, 5}, {2, 6},
		{2, 7}, {3, 1}, {3, 2}, {3, 3},
		{3, 4}, {3, 5}, {3, 6}, {3, 7},
		{4, 1}, {4, 2}, {4, 3}, {4, 4},
		{4, 5}, {4, 6}, {4, 7}, {5, 1},
		{5, 2}, {5, 3}, {5, 4}, {5, 5},
		{5, 6}, {5, 7},
	}
	src := new(RandSource64Mock)
	src.On("Seed", mock.Anything).Return()
	for i := 0; i < len(values); i++ {
		src.On("Uint64").Return(values[i]).Once()
	}
	r, err := NewHybridRandReaderWithSizeAndSource(shieldSize, src)
	if err != nil {
		t.Errorf("NewHybridRandReaderWithSizeAndSource() error = %v", err)
		return
	}
	size := 32
	b := make([]byte, size)
	n, err := r.Read(b)
	if err != nil {
		t.Errorf("Read() error = %v", err)
		return
	}
	if n != size {
		t.Errorf("Read() n = %v, want %v", n, 5)
		return
	}

	byteLst := make([]byte, 8)
	for i := 0; i < size; i++ {
		r := res[i]
		binary.LittleEndian.PutUint64(byteLst, values[r.idx])
		if b[i] != byteLst[r.pos] {
			t.Errorf("Read() b[%v] = %v, want %v", i, b[i], byteLst[r.pos])
			return
		}
	}
	src.AssertExpectations(t)
}
