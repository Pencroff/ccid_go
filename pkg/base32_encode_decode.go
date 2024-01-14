package pkg

import (
	"math/bits"
)

const (
	base32Alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"

	Base32strSize64  = 13
	Base32strSize96  = 20
	Base32strSize128 = 26
	Base32strSize160 = 32
)

var (
	reverseBase32Table = "" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\xff\xff\xff\xff\xff\xff" +
		"\xff\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\xff\x12\x13\xff\x14\x15\xff" +
		"\x16\x17\x18\x19\x1a\xff\x1b\x1c\x1d\x1e\x1f\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff"
)

// EncodeToBase32 encodes a byte slice to a base32 string.
// The byte order is big endian.
func EncodeToBase32(b []byte) (string, error) {
	l := len(b)
	size, err := getBase32strSize(byte(l))
	if err != nil {
		return "", err
	}
	res := make([]byte, size)
	asBase32(b, res)
	return string(res), nil
}

// DecodeFromBase32 decodes a base32 string to a byte slice.
// The byte order is big endian.
func DecodeFromBase32(str string) ([]byte, error) {
	l := len(str)
	size, err := getBase32byteSliceSize(byte(l))
	if err != nil {
		return []byte{}, err
	}
	res := make([]byte, size)
	err = fromBase32([]byte(str), res)
	if err != nil {
		return []byte{}, err
	}
	return res, nil
}

// byte		0       1       2       3       4
// bit 		1111111111111111111111111111111111111111
// char		|0  ||1  ||2  ||3  ||4  ||5  ||6  ||7  |
//
// byte		0       1       2
// bit 		111111111111111111111111
// char	   |3  ||4  ||5  ||6  ||7  |
//
// 64 bits = 8 bytes = 3 + 5 => 13 = 5 + 8 chars
// 96 bits = 12 bytes = 2 + 2*5 => 20 = 4 + 2*8 chars
// 128 bits = 16 bytes = 1 + 3*5 => 26 = 2 + 3*8 chars
// 160 bits = 20 bytes = 4*5 => 32 = 4*8 chars

func asBase32(src, dst []byte) {
	bytePointer := byte(len(src))
	pointer := byte(len(dst))
	buf := [8]byte{}
	for bytePointer > 3 {
		subSrc := src[bytePointer-5 : bytePointer]
		subDst := dst[pointer-8 : pointer]

		buf[7] = subSrc[4]
		buf[6] = (subSrc[4] >> 5) | (subSrc[3] << 3)
		buf[5] = subSrc[3] >> 2
		buf[4] = (subSrc[3] >> 7) | (subSrc[2] << 1)
		buf[3] = (subSrc[2] >> 4) | (subSrc[1] << 4)
		buf[2] = subSrc[1] >> 1
		buf[1] = (subSrc[1] >> 6) | (subSrc[0] << 2)
		buf[0] = subSrc[0] >> 3

		subDst[7] = base32Alphabet[buf[7]&0x1F]
		subDst[6] = base32Alphabet[buf[6]&0x1F]
		subDst[5] = base32Alphabet[buf[5]&0x1F]
		subDst[4] = base32Alphabet[buf[4]&0x1F]
		subDst[3] = base32Alphabet[buf[3]&0x1F]
		subDst[2] = base32Alphabet[buf[2]&0x1F]
		subDst[1] = base32Alphabet[buf[1]&0x1F]
		subDst[0] = base32Alphabet[buf[0]&0x1F]

		bytePointer -= 5
		pointer -= 8
	}

	subSrc := src[0:bytePointer]
	subDst := dst[0:pointer]

	switch bytePointer {
	case 3:
		buf[4] = subSrc[2]
		buf[3] = (subSrc[2] >> 5) | (subSrc[1] << 3)
		buf[2] = subSrc[1] >> 2
		buf[1] = (subSrc[1] >> 7) | (subSrc[0] << 1)
		buf[0] = subSrc[0] >> 4
	case 2:
		buf[3] = subSrc[1]
		buf[2] = (subSrc[1] >> 5) | (subSrc[0] << 3)
		buf[1] = subSrc[0] >> 2
		buf[0] = subSrc[0] >> 7
	case 1:
		buf[1] = subSrc[0]
		buf[0] = subSrc[0] >> 5
	}
	for i := pointer - 1; i < 255; i -= 1 {
		subDst[i] = base32Alphabet[buf[i]&0x1F]
	}
}

func fromBase32(src, dst []byte) error {
	pointer := byte(len(src))
	bytePointer := byte(len(dst))
	buf := [Base32strSize160]byte{}

	for i := pointer - 1; i < 255; i -= 1 {
		v := reverseBase32Table[src[i]]
		if v == 0xff {
			return InvalidCharacterError{src[i], i}
		}
		buf[i] = v
	}

	if bytePointer != ByteSliceSize160 {
		restBits := (bytePointer << 3) % 5
		requiredBits := byte(bits.Len8(buf[0]))
		if requiredBits > restBits {
			return OverflowError(bytePointer)
		}
	}

	for pointer >= 8 {
		bufSub := buf[pointer-8 : pointer]
		dstSub := dst[bytePointer-5 : bytePointer]

		dstSub[4] = bufSub[7] | (bufSub[6] << 5)
		dstSub[3] = (bufSub[6] >> 3) | (bufSub[5] << 2) | (bufSub[4] << 7)
		dstSub[2] = (bufSub[4] >> 1) | (bufSub[3] << 4)
		dstSub[1] = (bufSub[3] >> 4) | (bufSub[2] << 1) | (bufSub[1] << 6)
		dstSub[0] = (bufSub[1] >> 2) | (bufSub[0] << 3)

		pointer -= 8
		bytePointer -= 5
	}

	bufSub := buf[0:pointer]
	dstSub := dst[0:bytePointer]

	switch bytePointer {
	case 3:
		dstSub[2] = bufSub[4] | (bufSub[3] << 5)
		dstSub[1] = (bufSub[3] >> 3) | (bufSub[2] << 2) | (bufSub[1] << 7)
		dstSub[0] = (bufSub[1] >> 1) | (bufSub[0] << 4)
	case 2:
		dstSub[1] = bufSub[3] | (bufSub[2] << 5)
		dstSub[0] = (bufSub[2] >> 3) | (bufSub[1] << 2) | (bufSub[0] << 7)
	case 1:
		dstSub[0] = bufSub[1] | (bufSub[0] << 5)
	}

	return nil
}

func getBase32strSize(l byte) (byte, error) {
	var size byte
	switch l {
	case ByteSliceSize64:
		size = Base32strSize64
	case ByteSliceSize96:
		size = Base32strSize96
	case ByteSliceSize128:
		size = Base32strSize128
	case ByteSliceSize160:
		size = Base32strSize160
	default:
		return 0, InvalidLengthError(l)
	}
	return size, nil
}

func getBase32byteSliceSize(l byte) (byte, error) {
	var size byte
	switch l {
	case Base32strSize64:
		size = ByteSliceSize64
	case Base32strSize96:
		size = ByteSliceSize96
	case Base32strSize128:
		size = ByteSliceSize128
	case Base32strSize160:
		size = ByteSliceSize160
	default:
		return 0, InvalidLengthError(l)
	}
	return size, nil
}
