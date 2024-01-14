package pkg

import (
	"fmt"
	"testing"
)

func TestAsBaseBigInt(t *testing.T) {
	keys := SortKeys(testCaseEncodeDecodeMap)

	t.Run("base62", func(t *testing.T) {
		for _, name := range keys {
			tc := testCaseEncodeDecodeMap[name]
			t.Run(name, func(t *testing.T) {
				got, _ := AsBase62BigInt(tc.data)
				if got != tc.base62 {
					t.Errorf("AsBase62BigInt(%v) =\n'%s' (%d), want\n'%s' (%d)",
						tc.data, got, len(got), tc.base62, len(tc.base62))
				}
			})
		}
	})
	t.Run("base32", func(t *testing.T) {
		for _, name := range keys {
			tc := testCaseEncodeDecodeMap[name]
			t.Run(name, func(t *testing.T) {
				got, _ := AsBase32BigInt(tc.data)
				if got != tc.base32 {
					t.Errorf("AsBase32BigInt(%v) =\n'%s' (%d), want\n'%s' (%d)",
						tc.data, got, len(got), tc.base32, len(tc.base32))
				}
			})
		}
	})
	t.Run("base16", func(t *testing.T) {
		for _, name := range keys {
			tc := testCaseEncodeDecodeMap[name]
			t.Run(name, func(t *testing.T) {
				got, _ := AsBase16BigInt(tc.data)
				if got != tc.base16 {
					t.Errorf("AsBase16BigInt(%v) =\n'%s' (%d), want\n'%s' (%d)",
						tc.data, got, len(got), tc.base16, len(tc.base16))
				}
			})
			t.Run(fmt.Sprintf("sprintf_%s", name), func(t *testing.T) {
				got, _ := AsBase16Sprintf(tc.data)
				if got != tc.base16 {
					t.Errorf("AsBase16Sprintf(%v) =\n'%s' (%d), want\n'%s' (%d)",
						tc.data, got, len(got), tc.base16, len(tc.base16))
				}
			})
		}
	})
}

func TestAsBaseBigInt_Error(t *testing.T) {
	for i := 0; i < 32; i++ {
		if i == ByteSliceSize64 || i == ByteSliceSize96 || i == ByteSliceSize128 || i == ByteSliceSize160 {
			continue
		}
		a := make([]byte, i)
		_, err := AsBase62BigInt(a)
		if err == nil {
			t.Errorf("AsBase62BigInt(%v) error = %v, want %v", a, err, "invalid byte length")
		}
		_, err = AsBase32BigInt(a)
		if err == nil {
			t.Errorf("AsBase32BigInt(%v) error = %v, want %v", a, err, "invalid byte length")
		}
		_, err = AsBase16BigInt(a)
		if err == nil {
			t.Errorf("AsBase16BigInt(%v) error = %v, want %v", a, err, "invalid byte length")
		}
	}
}
