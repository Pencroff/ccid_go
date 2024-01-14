package pkg

import (
	"sort"
	"strings"
	"testing"
)

func TestBase16Alphabet(t *testing.T) {
	size := 16
	alphabet := base16Alphabet
	v := strings.Split(alphabet, "")
	sort.Strings(v)
	if strings.Join(v, "") != alphabet {
		t.Errorf("base%d alphabet is not Lexicographically sortable", size)
	}
}

func TestBase16Encode(t *testing.T) {
	size := len(base16Alphabet)
	keys := SortKeys(testCaseEncodeDecodeMap)
	for _, name := range keys {
		tc := testCaseEncodeDecodeMap[name]
		t.Run(name, func(t *testing.T) {
			got, _ := EncodeToBase16(tc.data)
			if got != tc.base16 {
				t.Errorf("EncodeToBase%d(%v) =\n'%s' (%d), want\n'%s' (%d)",
					size, tc.data, got, len(got), tc.base16, len(tc.base16))
			}
		})
	}
}

func TestBase16Encode_InvalidLength(t *testing.T) {
	size := len(base16Alphabet)
	for i := 0; i < 256; i++ {
		if i == ByteSliceSize64 || i == ByteSliceSize96 || i == ByteSliceSize128 || i == ByteSliceSize160 {
			continue
		}
		a := make([]byte, i)
		_, err := EncodeToBase16(a)
		if err == nil || strings.Index(err.Error(), "CCID: invalid length") == -1 {
			t.Errorf("EncodeToBase%d(%v) error = %v, want %v", size, a, err, "invalid byte length")
		}
	}
}

func TestBase16Decode(t *testing.T) {
	size := len(base16Alphabet)
	keys := SortKeys(testCaseEncodeDecodeMap)
	for _, name := range keys {
		tc := testCaseEncodeDecodeMap[name]
		t.Run(name, func(t *testing.T) {
			got, _ := DecodeFromBase16(tc.base16)
			if !SliceEqual(got, tc.data) {
				t.Errorf("DecodeFromBase%d(%v) =\n%x, want\n%x",
					size, tc.base16, got, tc.data)
			}
		})
	}
}

func TestBase16Decode_InvalidLength(t *testing.T) {
	size := len(base16Alphabet)
	for i := 0; i < 256; i++ {
		if i == Base16strSize64 || i == Base16strSize96 || i == Base16strSize128 || i == Base16strSize160 {
			continue
		}
		a := make([]byte, i)
		for j := 0; j < i; j++ {
			a[j] = "A"[0]
		}
		_, err := DecodeFromBase16(string(a))
		if err == nil || strings.Index(err.Error(), "CCID: invalid length") == -1 {
			t.Errorf("DecodeFromBase%d(%s) error =\n%v, want\n%v\nlength: %d", size, a, err, "CCID: invalid length", len(a))
		}
	}
}

func TestBase16Decode_InvalidCharacterError(t *testing.T) {
	size := len(base16Alphabet)
	lst := []byte{
		Base16strSize64,
		Base16strSize96,
		Base16strSize128,
		Base16strSize160,
	}
	for _, v := range lst {
		a := make([]byte, v)
		for j := byte(0); j < v; j++ {
			a[j] = "!"[0]
		}
		_, err := DecodeFromBase16(string(a))
		if err == nil || strings.Index(err.Error(), "CCID: invalid character") == -1 {
			t.Errorf("DecodeFromBase%d(%v) error = %v, want %v", size, a, err, "CCID: invalid character")
		}
	}
}
