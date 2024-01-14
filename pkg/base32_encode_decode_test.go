package pkg

import (
	"sort"
	"strings"
	"testing"
)

func TestBase32Alphabet(t *testing.T) {
	size := 32
	alphabet := base32Alphabet
	v := strings.Split(alphabet, "")
	sort.Strings(v)
	if strings.Join(v, "") != alphabet {
		t.Errorf("base%d alphabet is not Lexicographically sortable", size)
	}
}

func TestBase32Encode(t *testing.T) {
	size := len(base32Alphabet)
	keys := SortKeys(testCaseEncodeDecodeMap)
	for _, name := range keys {
		tc := testCaseEncodeDecodeMap[name]
		t.Run(name, func(t *testing.T) {
			got, _ := EncodeToBase32(tc.data)
			if got != tc.base32 {
				t.Errorf("\nEncodeToBase%d(%v) =\n'%s' (%d), want\n'%s' (%d)",
					size, tc.data, got, len(got), tc.base32, len(tc.base32))
			}
		})
		//implemented litlle endian byte order
		//std := base32.NewEncoding(base32Alphabet).WithPadding(base32.NoPadding)
		//t.Run(name+"_std", func(t *testing.T) {
		//	got := std.EncodeToString(tc.data)
		//	if got != tc.base32 {
		//		t.Errorf("\nEncodeToBase%d(%v) =\n'%s' (%d), want\n'%s' (%d)",
		//			size, tc.data, got, len(got), tc.base32, len(tc.base32))
		//	}
		//})
	}
}

func TestBase32Encode_InvalidLengthError(t *testing.T) {
	size := len(base32Alphabet)
	for i := 0; i < 256; i++ {
		if i == ByteSliceSize64 || i == ByteSliceSize96 || i == ByteSliceSize128 || i == ByteSliceSize160 {
			continue
		}
		a := make([]byte, i)
		_, err := EncodeToBase32(a)
		if err == nil || strings.Index(err.Error(), "CCID: invalid length") == -1 {
			t.Errorf("EncodeToBase%d(%v) error = %v, want %v", size, a, err, "invalid length")
		}
	}
}

func TestBase32Decode(t *testing.T) {
	size := len(base32Alphabet)
	keys := SortKeys(testCaseEncodeDecodeMap)
	for _, name := range keys {
		tc := testCaseEncodeDecodeMap[name]
		t.Run(name, func(t *testing.T) {
			got, err := DecodeFromBase32(tc.base32)
			if err != nil {
				t.Errorf("\nDecodeFromBase%d(%v) error =\n%v", size, tc.base32, err)
			}
			if !SliceEqual(got, tc.data) {
				t.Errorf("\nDecodeFromBase%d(%v) =\n%x, want\n%x",
					size, tc.base32, got, tc.data)
			}
		})
		//implemented litlle endian byte order
		//std := base32.NewEncoding(base32Alphabet).WithPadding(base32.NoPadding)
		//t.Run(name+"_std", func(t *testing.T) {
		//	got, _ := std.DecodeString(tc.base32)
		//	if !SliceEqual(got, tc.data) {
		//		t.Errorf("\nDecodeFromBase%d(%v) =\n%x, want\n%x",
		//			size, tc.base32, got, tc.data)
		//	}
		//})
	}
}

func TestBase32Decode_InvalidLengthError(t *testing.T) {
	size := len(base32Alphabet)
	for i := 0; i < 256; i++ {
		if i == Base32strSize64 || i == Base32strSize96 || i == Base32strSize128 || i == Base32strSize160 {
			continue
		}
		a := make([]byte, i)
		for j := 0; j < i; j++ {
			a[j] = "A"[0]
		}
		_, err := DecodeFromBase32(string(a))
		if err == nil || strings.Index(err.Error(), "CCID: invalid length") == -1 {
			t.Errorf("DecodeFromBase%d(%s) error =\n%v, want\n%v\nlength: %d", size, a, err, "CCID: invalid length", len(a))
		}
	}
}

func TestBase32Decode_InvalidCharacterError(t *testing.T) {
	size := len(base32Alphabet)
	lst := []byte{
		Base32strSize64,
		Base32strSize96,
		Base32strSize128,
		Base32strSize160,
	}
	for _, v := range lst {
		a := make([]byte, v)
		for j := byte(0); j < v; j++ {
			a[j] = "*"[0]
		}
		_, err := DecodeFromBase32(string(a))
		if err == nil || strings.Index(err.Error(), "CCID: invalid character") == -1 {
			t.Errorf("DecodeFromBase%d(%v) error = %v, want %v", size, a, err, "CCID: invalid character")
		}
	}
}

func TestDecodeBase32_OverflowError(t *testing.T) {
	keys := SortKeys(testCaseDecodeErrorMap)
	for _, name := range keys {
		tc := testCaseDecodeErrorMap[name]
		// 160 bit doesn't have overflow error
		if len(tc.base32) == 0 {
			continue
		}
		t.Run(name, func(t *testing.T) {
			_, err := DecodeFromBase32(tc.base32)
			if err == nil || strings.Index(err.Error(), "CCID: decode overflow") == -1 {
				t.Errorf("DecodeFromBase32(%v) error = %v, want %v", tc.base32, err, "CCID: decode overflow")
			}
		})
	}
}
