package ccid_go

import (
	"bytes"
	"fmt"
	e "github.com/Pencroff/ccid_go/extras"
	p "github.com/Pencroff/ccid_go/pkg"
	"strings"
	"testing"
	"time"
)

// Unexported
type mockClock struct {
	Val time.Time
}

func (m *mockClock) Now() time.Time {
	v := m.Val
	m.Val = v.Add(time.Second)
	return v
}

// Unexported
type mockStaticClock struct {
	Val time.Time
}

func (m *mockStaticClock) Now() time.Time {
	return m.Val
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

// Unexported
type mockStaticReader struct {
	Val byte
}

func (m *mockStaticReader) Read(p []byte) (n int, err error) {
	for i := range p {
		p[i] = m.Val
	}
	return len(p), nil
}

func TestFromStringByte_Error(t *testing.T) {
	t.Run("ccid64", func(t *testing.T) {
		IterateMapAndTestFromStringBytesError(t, p.TestCaseCcId64Map)
	})
	t.Run("ccid96", func(t *testing.T) {
		IterateMapAndTestFromStringBytesError(t, p.TestCaseCcId96Map)
	})
	t.Run("ccid128", func(t *testing.T) {
		IterateMapAndTestFromStringBytesError(t, p.TestCaseCcId128Map)
	})
	t.Run("ccid160", func(t *testing.T) {
		IterateMapAndTestFromStringBytesError(t, p.TestCaseCcId160Map)
	})
}

func IterateMapAndTestFromStringBytesError(t *testing.T, m map[string]p.CcIdTestCases) {
	keys := p.SortKeys(m)
	for _, key := range keys {
		tc := m[key]
		t.Run(key+"_from_string_62_larger", func(t *testing.T) {
			v := tc.Base62 + "0"
			_, err := FromString(v, byte(len(tc.Fingerprint)), 62)
			if err == nil || strings.Index(err.Error(), "CCID: invalid length") != 0 {
				t.Errorf("FromString(%s) =\n%s, want\n%s",
					tc.Base62, err, "invalid byte length")
			}
		})
		t.Run(key+"_from_string_62_smaller", func(t *testing.T) {
			l := len(tc.Base62)
			v := tc.Base62[:l-1]
			_, err := FromString(v, byte(len(tc.Fingerprint)), 62)
			if err == nil || strings.Index(err.Error(), "CCID: invalid length") != 0 {
				t.Errorf("FromString(%s) =\n%s, want\n%s",
					tc.Base62, err, "invalid byte length")
			}
		})
		t.Run(key+"_from_string_32_larger", func(t *testing.T) {
			v := tc.Base32 + "0"
			_, err := FromString(v, byte(len(tc.Fingerprint)), 32)
			if err == nil || strings.Index(err.Error(), "CCID: invalid length") != 0 {
				t.Errorf("FromString(%s) =\n%s, want\n%s",
					tc.Base32, err, "CCID: invalid length")
			}
		})
		t.Run(key+"_from_string_32_smaller", func(t *testing.T) {
			l := len(tc.Base32)
			v := tc.Base32[:l-1]
			_, err := FromString(v, byte(len(tc.Fingerprint)), 32)
			if err == nil || strings.Index(err.Error(), "CCID: invalid length") != 0 {
				t.Errorf("FromString(%s) =\n%s, want\n%s",
					tc.Base32, err, "CCID: invalid length")
			}
		})
		t.Run(key+"_from_string_16_larger", func(t *testing.T) {
			v := tc.Base16 + "0"
			_, err := FromString(v, byte(len(tc.Fingerprint)), 16)
			if err == nil || strings.Index(err.Error(), "CCID: invalid length") != 0 {
				t.Errorf("FromString(%s) =\n%s, want\n%s",
					tc.Base16, err, "CCID: invalid length")
			}
		})
		t.Run(key+"_from_string_16_smaller", func(t *testing.T) {
			l := len(tc.Base16)
			v := tc.Base16[:l-1]
			_, err := FromString(v, byte(len(tc.Fingerprint)), 16)
			if err == nil || strings.Index(err.Error(), "CCID: invalid length") != 0 {
				t.Errorf("FromString(%s) =\n%s, want\n%s",
					tc.Base16, err, "CCID: invalid length")
			}
		})
		t.Run(key+"_from_bytes_larger", func(t *testing.T) {
			v := append(tc.Bytes, 0)
			_, err := FromBytes(v, byte(len(tc.Fingerprint)))
			if err == nil || strings.Index(err.Error(), "CCID: invalid length") != 0 {
				t.Errorf("FromBytes(%x) =\n%s, want\n%s",
					v, err, "CCID: invalid length")
			}
		})
		t.Run(key+"_from_bytes_smaller", func(t *testing.T) {
			l := len(tc.Bytes)
			v := tc.Bytes[:l-1]
			_, err := FromBytes(v, byte(len(tc.Fingerprint)))
			if err == nil || strings.Index(err.Error(), "CCID: invalid length") != 0 {
				t.Errorf("FromBytes(%x) =\n%s, want\n%s",
					v, err, "CCID: invalid length")
			}
		})
	}
}

func TestFromStringByte(t *testing.T) {
	t.Run("ccid64", func(t *testing.T) {
		IterateMapAndTestFromStringBytes(t, p.TestCaseCcId64Map)
	})
	t.Run("ccid96", func(t *testing.T) {
		IterateMapAndTestFromStringBytes(t, p.TestCaseCcId96Map)
	})
	t.Run("ccid128", func(t *testing.T) {
		IterateMapAndTestFromStringBytes(t, p.TestCaseCcId128Map)
	})
	t.Run("ccid160", func(t *testing.T) {
		IterateMapAndTestFromStringBytes(t, p.TestCaseCcId160Map)
	})
}

func IterateMapAndTestFromStringBytes(t *testing.T, m map[string]p.CcIdTestCases) {
	keys := p.SortKeys(m)
	for _, key := range keys {
		tc := m[key]
		t.Run(key+"_from_string_62", func(t *testing.T) {
			got, _ := FromString(tc.Base62, byte(len(tc.Fingerprint)), 62)
			v := fmt.Sprintf("%#v", got)
			if v != tc.GoString {
				t.Errorf("FromString(%s) =\n%s, want\n%s",
					tc.Base62, v, tc.GoString)
			}
		})
		t.Run(key+"_from_string_32", func(t *testing.T) {
			got, _ := FromString(tc.Base32, byte(len(tc.Fingerprint)), 32)
			v := fmt.Sprintf("%#v", got)
			if v != tc.GoString {
				t.Errorf("FromString(%s) =\n%s, want\n%s",
					tc.Base32, v, tc.GoString)
			}
		})
		t.Run(key+"_from_string_16", func(t *testing.T) {
			got, _ := FromString(tc.Base16, byte(len(tc.Fingerprint)), 16)
			v := fmt.Sprintf("%#v", got)
			if v != tc.GoString {
				t.Errorf("FromString(%s) =\n%s, want\n%s",
					tc.Base62, v, tc.GoString)
			}
		})
		t.Run(key+"_from_bytes", func(t *testing.T) {
			got, _ := FromBytes(tc.Bytes, byte(len(tc.Fingerprint)))
			v := fmt.Sprintf("%#v", got)
			if v != tc.GoString {
				t.Errorf("FromBytes(%x) =\n%s, want\n%s",
					tc.Bytes, v, tc.GoString)
			}
		})
	}
}

func TestNonMonotonicCcIdGen(t *testing.T) {
	keys := p.SortKeys(testCaseCcIdGenMap)

	t.Run("Next", func(t *testing.T) {
		for _, key := range keys {
			tc := testCaseCcIdGenMap[key]
			r := &mockReader{Val: 0xA5}
			c := &mockClock{Val: tc.mockTime}
			gen, _ := newCcIdGenWithClock(tc.size, tc.fingerprint, r, nil, c)
			t.Run(key, func(t *testing.T) {
				idA, _ := gen.Next()
				idB, _ := gen.Next()
				if idA.Size() != tc.size {
					t.Errorf("A: Incorrect size. Get:\n%d, want\n%d", idA.Size(), tc.size)
				}
				if idB.Size() != tc.size {
					t.Errorf("B: Incorrect size. Get:\n%d, want\n%d", idB.Size(), tc.size)
				}
				bA := idA.Bytes()
				if !bytes.Equal(idA.Bytes(), tc.nextA) {
					t.Errorf("A: Incorrect bytes. Get:\n%#v, want\n%#v", bA, tc.nextA)
				}
				bB := idB.Bytes()
				if !bytes.Equal(idB.Bytes(), tc.nextB) {
					t.Errorf("B: Incorrect bytes. Get:\n%#v, want\n%#v", bB, tc.nextB)
				}
			})
		}
	})
	t.Run("NextWithTime", func(t *testing.T) {
		for _, key := range keys {
			tc := testCaseCcIdGenMap[key]
			r := &mockReader{Val: 0xA5}
			gen, _ := NewMonotonicCcIdGenWithFingerprint(tc.size, tc.fingerprint, r, nil) // Used real clock under the hood
			t.Run(key, func(t *testing.T) {
				idA, _ := gen.NextWithTime(tc.mockTime)
				idB, _ := gen.NextWithTime(tc.mockTime.Add(time.Second))
				if idA.Size() != tc.size {
					t.Errorf("A: Incorrect size. Get:\n%d, want\n%d", idA.Size(), tc.size)
				}
				if idB.Size() != tc.size {
					t.Errorf("B: Incorrect size. Get:\n%d, want\n%d", idB.Size(), tc.size)
				}
				bA := idA.Bytes()
				if !bytes.Equal(idA.Bytes(), tc.nextA) {
					t.Errorf("A: Incorrect bytes. Get:\n%#v, want\n%#v", bA, tc.nextA)
				}
				bB := idB.Bytes()
				if !bytes.Equal(idB.Bytes(), tc.nextB) {
					t.Errorf("B: Incorrect bytes. Get:\n%#v, want\n%#v", bB, tc.nextB)
				}
			})
		}
	})
}

func TestMonotonicCcIdGen(t *testing.T) {
	keys := p.SortKeys(testCaseMonotonicCcIdGenMap)

	t.Run("Next", func(t *testing.T) {
		for _, key := range keys {
			tc := testCaseMonotonicCcIdGenMap[key]
			r := &mockStaticReader{Val: 0xA5}
			c := &mockStaticClock{Val: tc.mockTime}
			s := p.NewFiftyPercentMonotonicStrategy(r)
			gen, _ := newCcIdGenWithClock(tc.size, tc.fingerprint, r, s, c)
			t.Run(key, func(t *testing.T) {
				idA, _ := gen.Next()
				idB, _ := gen.Next()
				if idA.Size() != tc.size {
					t.Errorf("A: Incorrect size. Get:\n%d, want\n%d", idA.Size(), tc.size)
				}
				if idB.Size() != tc.size {
					t.Errorf("B: Incorrect size. Get:\n%d, want\n%d", idB.Size(), tc.size)
				}
				bA := idA.Bytes()
				if !bytes.Equal(idA.Bytes(), tc.nextA) {
					t.Errorf("A: Incorrect bytes. Get:\n%#v, want\n%#v", bA, tc.nextA)
				}
				bB := idB.Bytes()
				if !bytes.Equal(idB.Bytes(), tc.nextB) {
					t.Errorf("B: Incorrect bytes. Get:\n%#v, want\n%#v", bB, tc.nextB)
				}
			})
		}
	})
	t.Run("NextWithTime", func(t *testing.T) {
		for _, key := range keys {
			tc := testCaseMonotonicCcIdGenMap[key]
			r := &mockStaticReader{Val: 0xA5}
			s := p.NewFiftyPercentMonotonicStrategy(r)
			gen, _ := NewMonotonicCcIdGenWithFingerprint(tc.size, tc.fingerprint, r, s) // Used real clock under the hood
			t.Run(key, func(t *testing.T) {
				idA, _ := gen.NextWithTime(tc.mockTime)
				idB, _ := gen.NextWithTime(tc.mockTime)
				if idA.Size() != tc.size {
					t.Errorf("A: Incorrect size. Get:\n%d, want\n%d", idA.Size(), tc.size)
				}
				if idB.Size() != tc.size {
					t.Errorf("B: Incorrect size. Get:\n%d, want\n%d", idB.Size(), tc.size)
				}
				bA := idA.Bytes()
				if !bytes.Equal(idA.Bytes(), tc.nextA) {
					t.Errorf("A: Incorrect bytes. Get:\n%#v, want\n%#v", bA, tc.nextA)
				}
				bB := idB.Bytes()
				if !bytes.Equal(idB.Bytes(), tc.nextB) {
					t.Errorf("B: Incorrect bytes. Get:\n%#v, want\n%#v", bB, tc.nextB)
				}
			})
		}
	})
}

func TestMonotonicOverloadCcIdGen(t *testing.T) {
	keys := p.SortKeys(testCaseMonotonicOverloadCcIdGenMap)

	t.Run("Next", func(t *testing.T) {
		for _, key := range keys {
			tc := testCaseMonotonicOverloadCcIdGenMap[key]
			r := &mockStaticReader{Val: 0xFF}
			c := &mockStaticClock{Val: tc.mockTime}
			s := p.NewFiftyPercentMonotonicStrategy(r)
			gen, _ := newCcIdGenWithClock(tc.size, tc.fingerprint, r, s, c)
			t.Run(key, func(t *testing.T) {
				idA, _ := gen.Next()
				idB, _ := gen.Next()
				if idA.Size() != tc.size {
					t.Errorf("A: Incorrect size. Get:\n%d, want\n%d", idA.Size(), tc.size)
				}
				if idB.Size() != tc.size {
					t.Errorf("B: Incorrect size. Get:\n%d, want\n%d", idB.Size(), tc.size)
				}
				bA := idA.Bytes()
				if !bytes.Equal(idA.Bytes(), tc.nextA) {
					t.Errorf("A: Incorrect bytes. Get:\n%x, want\n%x", bA, tc.nextA)
				}
				bB := idB.Bytes()
				if !bytes.Equal(idB.Bytes(), tc.nextB) {
					t.Errorf("B: Incorrect bytes. Get:\n%x, want\n%x", bB, tc.nextB)
				}
			})
		}
	})
	t.Run("NextWithTime", func(t *testing.T) {
		for _, key := range keys {
			tc := testCaseMonotonicOverloadCcIdGenMap[key]
			r := &mockStaticReader{Val: 0xFF}
			s := p.NewFiftyPercentMonotonicStrategy(r)
			gen, _ := NewMonotonicCcIdGenWithFingerprint(tc.size, tc.fingerprint, r, s) // Used real clock under the hood
			t.Run(key, func(t *testing.T) {
				idA, _ := gen.NextWithTime(tc.mockTime)
				idB, _ := gen.NextWithTime(tc.mockTime)
				if idA.Size() != tc.size {
					t.Errorf("A: Incorrect size. Get:\n%d, want\n%d", idA.Size(), tc.size)
				}
				if idB.Size() != tc.size {
					t.Errorf("B: Incorrect size. Get:\n%d, want\n%d", idB.Size(), tc.size)
				}
				bA := idA.Bytes()
				if !bytes.Equal(idA.Bytes(), tc.nextA) {
					t.Errorf("A: Incorrect bytes. Get:\n%x, want\n%x", bA, tc.nextA)
				}
				bB := idB.Bytes()
				if !bytes.Equal(idB.Bytes(), tc.nextB) {
					t.Errorf("B: Incorrect bytes. Get:\n%x, want\n%x", bB, tc.nextB)
				}
			})
		}
	})
}

func TestMonotonicRandCcIdGen(t *testing.T) {
	sizes := []byte{
		p.ByteSliceSize64,
		p.ByteSliceSize96,
		p.ByteSliceSize128,
		p.ByteSliceSize160,
	}
	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			r := &e.SecureRandReader{}
			s := p.NewFiftyPercentMonotonicStrategy(r)
			gen, _ := NewMonotonicCcIdGen(size, r, s) // Used real clock under the hood
			t.Run("Next", func(t *testing.T) {
				idA, _ := gen.Next()
				idB, _ := gen.Next()
				if bytes.Compare(idA.Bytes(), idB.Bytes()) >= 0 {
					t.Errorf("Not monotonic.\nidA: %x\nmore then\nidB: %x", idA.Bytes(), idB.Bytes())
				}
			})
			t.Run("NextWithTime", func(t *testing.T) {
				idA, _ := gen.NextWithTime(p.RealClock{}.Now())
				idB, _ := gen.NextWithTime(p.RealClock{}.Now())
				if bytes.Compare(idA.Bytes(), idB.Bytes()) >= 0 {
					t.Errorf("Not monotonic.\nidA: %x\nmore then\nidB: %x", idA.Bytes(), idB.Bytes())
				}
			})
		})
	}
}

func TestMonotonicRandCcIdGen_LargeCycle(t *testing.T) {
	sizes := []byte{
		p.ByteSliceSize64,
		p.ByteSliceSize96,
		p.ByteSliceSize128,
		p.ByteSliceSize160,
	}
	fingerPrints := [][]byte{
		make([]byte, 1),
		make([]byte, 5),
		make([]byte, 5),
		make([]byte, 5),
	}
	for idx, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			r, _ := e.NewHybridRandReader()
			s := p.NewFiftyPercentMonotonicStrategy(r)
			gen, _ := NewMonotonicCcIdGen(size, r, s) // Used real clock under the hood
			t.Run("Next", func(t *testing.T) {
				prev, _ := gen.Next()
				for i := 0; i < 1_000_000; i++ {
					next, _ := gen.Next()
					if bytes.Compare(prev.Bytes(), next.Bytes()) >= 0 {
						t.Errorf("%d - Not monotonic.\nPrev: %x\nmore then\nNext: %x", i, prev.Bytes(), next.Bytes())
						t.Errorf("Details:\nPrev: %p - %#v\nNext: %p - %#v", &prev, prev, &next, next)
					}
					if bytes.Compare(prev.Bytes()[:4], next.Bytes()[:4]) > 0 {
						t.Errorf("%d - Time decreased. Now: %s.\nPrev: %#v\nNext: %#v", i, time.Now().Format(time.RFC3339), prev, next)
					}
					prev = next
				}
			})
		})
		t.Run(fmt.Sprintf("size_%d_with_fingerprint", size), func(t *testing.T) {
			r, _ := e.NewHybridRandReader()
			s := p.NewFiftyPercentMonotonicStrategy(r)
			gen, _ := NewMonotonicCcIdGenWithFingerprint(size, fingerPrints[idx], r, s) // Used real clock under the hood
			t.Run("Next", func(t *testing.T) {
				prev, _ := gen.Next()
				for i := 0; i < 1_000_000; i++ {
					next, _ := gen.Next()
					if bytes.Compare(prev.Bytes(), next.Bytes()) >= 0 {
						t.Errorf("%d - Not monotonic.\nPrev: %x\nmore then\nNext: %x", i, prev.Bytes(), next.Bytes())
						t.Errorf("Details:\nPrev: %p - %#v\nNext: %p - %#v", &prev, prev, &next, next)
					}
					if bytes.Compare(prev.Bytes()[:4], next.Bytes()[:4]) > 0 {
						t.Errorf("%d - Time decreased. Now: %s.\nPrev: %#v\nNext: %#v", i, time.Now().Format(time.RFC3339), prev, next)
					}
					prev = next
				}
			})
		})
	}
}
