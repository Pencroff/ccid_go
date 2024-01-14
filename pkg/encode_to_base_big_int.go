package pkg

import (
	"fmt"
	"math/big"
)

// All code in this file just for benchmarking and testing purposes

// Deprecated: AsBase62BigInt transforms a byte slice to a base62 string
// Implemented for benchmarking and testing purposes
func AsBase62BigInt(b []byte) (string, error) {
	l := len(b)
	if l != 8 && l != 12 && l != 16 && l != 20 {
		return "", fmt.Errorf("invalid byte length %d", l)
	}
	size := Base62strSize64
	if l == 12 {
		size = Base62strSize96
	}
	if l == 16 {
		size = Base62strSize128
	}
	if l == 20 {
		size = Base62strSize160
	}
	return asBaseBigInt(b, base62Alphabet, size), nil
}

// Deprecated: AsBase32BigInt transforms a byte slice to a base32 string
// Implemented for benchmarking and testing purposes
func AsBase32BigInt(b []byte) (string, error) {
	l := len(b)
	if l != 8 && l != 12 && l != 16 && l != 20 {
		return "", fmt.Errorf("invalid byte length %d", l)
	}
	size := Base32strSize64
	if l == 12 {
		size = Base32strSize96
	}
	if l == 16 {
		size = Base32strSize128
	}
	if l == 20 {
		size = Base32strSize160
	}
	return asBaseBigInt(b, base32Alphabet, size), nil
}

// Deprecated: AsBase16Sprintf has same speed as AsBase62BigInt
func AsBase16Sprintf(b []byte) (string, error) {
	l := len(b)
	if l != 8 && l != 12 && l != 16 && l != 20 {
		return "", fmt.Errorf("invalid byte length %d", l)
	}
	r := fmt.Sprintf("%X", b)
	return r, nil
}

// Deprecated: AsBase16BigInt transforms a byte slice to a base16 string
// Implemented for benchmarking and testing purposes
func AsBase16BigInt(b []byte) (string, error) {
	l := len(b)
	if l != 8 && l != 12 && l != 16 && l != 20 {
		return "", fmt.Errorf("invalid byte length %d", l)
	}
	size := Base16strSize64
	if l == 12 {
		size = Base16strSize96
	}
	if l == 16 {
		size = Base16strSize128
	}
	if l == 20 {
		size = Base16strSize160
	}
	return asBaseBigInt(b, base16Alphabet, size), nil
}

// Deprecated: asBaseBigInt transforms a byte slice to a string using the given alphabet
// Implemented for benchmarking and testing purposes
func asBaseBigInt(b []byte, alphabet string, size int) string {
	alphabetSize := len(alphabet)
	base := big.NewInt(int64(alphabetSize))
	n := new(big.Int).SetBytes(b)
	mod := new(big.Int)
	r := make([]byte, size)
	pointer := size - 1
	for {
		if n.Cmp(base) < 0 {
			r[pointer] = alphabet[n.Int64()]
			break
		}
		n.DivMod(n, base, mod)
		r[pointer] = alphabet[mod.Int64()]
		pointer -= 1
		if pointer < 0 {
			break
		}
	}
	if pointer > 0 {
		copy(r[:pointer], zeroString[:pointer])
	}
	return string(r)
}
