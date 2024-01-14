package pkg

const (
	base16Alphabet = "0123456789ABCDEF"

	Base16strSize64  = 16
	Base16strSize96  = 24
	Base16strSize128 = 32
	Base16strSize160 = 40
)

var (
	reverseBase16Table = "" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\xff\xff\xff\xff\xff\xff" +
		"\xff\x0a\x0b\x0c\x0d\x0e\x0f\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
		"\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff\xff" +
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

// EncodeToBase16 encodes a byte slice to a base16 string.
// The byte order is big endian.
func EncodeToBase16(b []byte) (string, error) {
	l := len(b)
	size, err := getBase16strSize(byte(l))
	if err != nil {
		return "", err
	}
	r := [Base16strSize160]byte{}
	idx := 0
	for _, v := range b {
		r[idx] = base16Alphabet[v>>4]
		r[idx+1] = base16Alphabet[v&0x0f]
		idx += 2
	}
	return string(r[:size]), nil
}

// DecodeFromBase16 decodes a base16 string to a byte slice.
// The byte order is big endian.
func DecodeFromBase16(str string) ([]byte, error) {
	pointer := byte(len(str))
	size, err := getBase16byteSliceSize(pointer)
	if err != nil {
		return []byte{}, err
	}
	res := [ByteSliceSize160]byte{}
	bytePointer := size - 1
	for i := pointer - 2; i < 250; i -= 2 {
		vLow := reverseBase16Table[str[i+1]]
		vHigh := reverseBase16Table[str[i]]
		if vLow == 0xff || vHigh == 0xff {
			return []byte{}, InvalidCharacterError{str[i], i}
		}
		res[bytePointer] = vHigh<<4 | vLow
		bytePointer -= 1
	}
	return res[:size], nil
}

func getBase16strSize(l byte) (byte, error) {
	var size byte
	switch l {
	case ByteSliceSize64:
		size = Base16strSize64
	case ByteSliceSize96:
		size = Base16strSize96
	case ByteSliceSize128:
		size = Base16strSize128
	case ByteSliceSize160:
		size = Base16strSize160
	default:
		return 0, InvalidLengthError(l)
	}
	return size, nil
}

func getBase16byteSliceSize(l byte) (byte, error) {
	var size byte
	switch l {
	case Base16strSize64:
		size = ByteSliceSize64
	case Base16strSize96:
		size = ByteSliceSize96
	case Base16strSize128:
		size = ByteSliceSize128
	case Base16strSize160:
		size = ByteSliceSize160
	default:
		return 0, InvalidLengthError(l)
	}
	return size, nil
}
