package pkg

import (
	"sort"
	"strings"
	"testing"
)

func TestBase62Alphabet(t *testing.T) {
	size := 62
	alphabet := base62Alphabet
	v := strings.Split(alphabet, "")
	sort.Strings(v)
	if strings.Join(v, "") != alphabet {
		t.Errorf("base%d alphabet is not Lexicographically sortable", size)
	}
}

func TestBase62Encode(t *testing.T) {
	size := len(base62Alphabet)
	keys := SortKeys(testCaseEncodeDecodeMap)
	for _, name := range keys {
		tc := testCaseEncodeDecodeMap[name]
		t.Run(name, func(t *testing.T) {
			got, _ := EncodeToBase62(tc.data)
			if got != tc.base62 {
				t.Errorf("EncodeToBase%d(%v) =\n'%s' (%d), want\n'%s' (%d)",
					size, tc.data, got, len(got), tc.base62, len(tc.base62))
			}
		})
	}
}

func TestBase62Encode_InvalidLengthError(t *testing.T) {
	size := len(base62Alphabet)
	for i := 0; i < 256; i++ {
		if i == ByteSliceSize64 || i == ByteSliceSize96 || i == ByteSliceSize128 || i == ByteSliceSize160 {
			continue
		}
		a := make([]byte, i)
		_, err := EncodeToBase62(a)
		if err == nil || strings.Index(err.Error(), "CCID: invalid length") == -1 {
			t.Errorf("EncodeToBase%d(%v) error = %v, want %v", size, a, err, "invalid byte length")
		}
	}
}

func TestBase62Decode(t *testing.T) {
	size := len(base62Alphabet)
	keys := SortKeys(testCaseEncodeDecodeMap)
	for _, name := range keys {
		tc := testCaseEncodeDecodeMap[name]
		t.Run(name, func(t *testing.T) {
			got, _ := DecodeFromBase62(tc.base62)
			if !SliceEqual(got, tc.data) {
				t.Errorf("DecodeFromBase%d(%v) =\n%x, want\n%x",
					size, tc.base62, got, tc.data)
			}
		})
	}
}

func TestBase62Decode_InvalidLengthError(t *testing.T) {
	size := len(base62Alphabet)
	for i := 0; i < 256; i++ {
		if i == Base62strSize64 || i == Base62strSize96 || i == Base62strSize128 || i == Base62strSize160 {
			continue
		}
		a := make([]byte, i)
		for j := 0; j < i; j++ {
			a[j] = "A"[0]
		}
		_, err := DecodeFromBase62(string(a))
		if err == nil || strings.Index(err.Error(), "CCID: invalid length") == -1 {
			t.Errorf("DecodeFromBase%d(%v) error = %v, want %v", size, a, err, "invalid byte length")
		}
	}
}

func TestDecodeFromBase62_InvalidCharacterError(t *testing.T) {
	size := len(base16Alphabet)
	lst := []byte{
		Base62strSize64,
		Base62strSize96,
		Base62strSize128,
		Base62strSize160,
	}
	for _, v := range lst {
		a := make([]byte, v)
		for j := byte(0); j < v; j++ {
			a[j] = "!"[0]
		}
		_, err := DecodeFromBase62(string(a))
		if err == nil || strings.Index(err.Error(), "CCID: invalid character") == -1 {
			t.Errorf("DecodeFromBase%d(%v) error = %v, want %v", size, a, err, "CCID: invalid character")
		}
	}
}

func TestDecodeBase62_OverflowError(t *testing.T) {
	keys := SortKeys(testCaseDecodeErrorMap)
	for _, name := range keys {
		tc := testCaseDecodeErrorMap[name]
		t.Run(name, func(t *testing.T) {
			_, err := DecodeFromBase62(tc.base62)
			if err == nil || strings.Index(err.Error(), "CCID: decode overflow") == -1 {
				t.Errorf("DecodeFromBase62(%v) error = %v, want %v", tc.base62, err, "CCID: decode overflow")
			}
		})
	}
}
