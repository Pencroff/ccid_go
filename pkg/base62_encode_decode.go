package pkg

import (
	"encoding/binary"
)

const (
	base62Alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	Base62strSize64  = 11
	Base62strSize96  = 17
	Base62strSize128 = 22
	Base62strSize160 = 27
)

var (
	reverseBase62Table = "" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\xff\xff\xff\xff\xff\xff" +
		"\xff\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18" +
		"\x19\x1a\x1b\x1c\x1d\x1e\x1f\x20\x21\x22\x23\xff\xff\xff\xff\xff" +
		"\xff\x24\x25\x26\x27\x28\x29\x2a\x2b\x2c\x2d\x2e\x2f\x30\x31\x32" +
		"\x33\x34\x35\x36\x37\x38\x39\x3a\x3b\x3c\x3d\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"
)

// EncodeToBase62 encodes a byte slice to a base62 string.
// The byte order is big endian.
// Implementation is based on ksuid quick algorithm.
func EncodeToBase62(b []byte) (string, error) {
	l := len(b)
	size, err := getBase62strSize(byte(l))
	if err != nil {
		return "", err
	}
	res := make([]byte, size)
	asBase62(b, res, base62Alphabet)
	return string(res), nil
}

// DecodeFromBase62 decodes a base62 string to a byte slice.
// The byte order is big endian.
// Implementation is based on ksuid quick algorithm.
func DecodeFromBase62(str string) ([]byte, error) {
	l := len(str)
	size, err := getBase62byteSliceSize(byte(l))
	if err != nil {
		return []byte{}, err
	}
	res := make([]byte, size)
	err = fromBase62([]byte(str), res)
	if err != nil {
		return []byte{}, err
	}
	return res, nil

}

func asBase62(src, dst []byte, alphabet string) {
	const dstBase = 62 // len(alphabet)
	bytePointer := len(src)
	pointer := len(dst)

	partsSize := bytePointer >> 2
	parts := [5]uint64{}
	for i := 0; i < partsSize; i += 1 {
		idx := i << 2
		parts[i] = uint64(binary.BigEndian.Uint32(src[idx : idx+4]))
	}

	bp := parts[:partsSize]
	bq := [5]uint64{}

	for len(bp) != 0 {
		quotient := bq[:0]
		digit := uint64(0)
		remainder := uint64(0)
		for _, c := range bp {
			value := c | remainder<<32 // uint64(c) + uint64(remainder)*srcBase
			digit = value / dstBase
			remainder = value % dstBase

			if len(quotient) != 0 || digit != 0 {
				quotient = append(quotient, digit)
			}
		}

		pointer -= 1
		dst[pointer] = alphabet[remainder]
		bp = quotient
	}

	copy(dst[:pointer], zeroString)
}

func fromBase62(src, dst []byte) error {
	const srcBase = 62
	const dstBase = 1 << 32 // 4294967296 // 2^32
	const dstMask = dstBase - 1
	pointer := byte(len(src))
	bytePointer := byte(len(dst))

	parts := [Base62strSize160]byte{}

	// This line helps BCE (Bounds Check Elimination).
	// It may be safely removed.
	_ = src[pointer-1]

	for i := pointer - 1; i < 255; i -= 1 {
		v := reverseBase62Table[src[i]]
		if v == 0xff {
			return InvalidCharacterError{src[i], i}
		}
		parts[i] = v
	}

	bp := parts[:pointer]
	bq := [Base62strSize160]byte{}

	for len(bp) > 0 {
		quotient := bq[:0]
		remainder := uint64(0)

		for _, c := range bp {
			value := uint64(c) + uint64(remainder)*srcBase
			digit := value >> 32        // value / dstBase
			remainder = value & dstMask // value % dstBase

			if len(quotient) != 0 || digit != 0 {
				quotient = append(quotient, byte(digit))
			}
		}

		if bytePointer < 4 {
			return OverflowError(byte(len(dst)))
		}

		dst[bytePointer-4] = byte(remainder >> 24)
		dst[bytePointer-3] = byte(remainder >> 16)
		dst[bytePointer-2] = byte(remainder >> 8)
		dst[bytePointer-1] = byte(remainder)
		bytePointer -= 4
		bp = quotient
	}

	var zero [20]byte
	copy(dst[:bytePointer], zero[:])
	return nil
}

func getBase62strSize(l byte) (byte, error) {
	var size byte
	switch l {
	case ByteSliceSize64:
		size = Base62strSize64
	case ByteSliceSize96:
		size = Base62strSize96
	case ByteSliceSize128:
		size = Base62strSize128
	case ByteSliceSize160:
		size = Base62strSize160
	default:
		return 0, InvalidLengthError(l)
	}
	return size, nil
}

func getBase62byteSliceSize(l byte) (byte, error) {
	var size byte
	switch l {
	case Base62strSize64:
		size = ByteSliceSize64
	case Base62strSize96:
		size = ByteSliceSize96
	case Base62strSize128:
		size = ByteSliceSize128
	case Base62strSize160:
		size = ByteSliceSize160
	default:
		return 0, InvalidLengthError(l)
	}
	return size, nil
}
