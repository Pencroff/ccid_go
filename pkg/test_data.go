//go:build !pkgsite
// +build !pkgsite

package pkg

import (
	"sort"
	"time"
)

type testCaseEncodeDecode struct {
	data   []byte
	base62 string
	base32 string
	base16 string
}

var testCaseDecodeErrorMap = map[string]testCaseEncodeDecode{
	"64 bit overflow next": {
		[]byte{},
		"LygHa16AHYG",
		"G000000000000",
		"",
	},
	"96 bit overflow next": {
		[]byte{},
		"1f2SI9UJPXvb7vdJ2",
		"20000000000000000000",
		"",
	},
	"128 bit overflow next": {
		[]byte{},
		"7n42DGM5Tflk9n8mt7Fhc8",
		"80000000000000000000000000",
		"",
	},
	"160 bit overflow next": {
		[]byte{},
		"aWgEPTl1tmebfsQzFP4bxwgy80W",
		"",
		"",
	},
	"64 bit overflow": {
		[]byte{},
		"MygHa16AHYF",
		"GZZZZZZZZZZZZ",
		"",
	},
	"96 bit overflow": {
		[]byte{},
		"2f2SI9UJPXvb7vdJ1",
		"2ZZZZZZZZZZZZZZZZZZZ",
		"",
	},
	"128 bit overflow": {
		[]byte{},
		"8n42DGM5Tflk9n8mt7Fhc7",
		"8ZZZZZZZZZZZZZZZZZZZZZZZZZ",
		"",
	},
	"160 bit overflow": {
		[]byte{},
		"bWgEPTl1tmebfsQzFP4bxwgy80V",
		"",
		"",
	},
	"64 bit overflow max": {
		[]byte{},
		"zzzzzzzzzzz",
		"ZZZZZZZZZZZZZ",
		"",
	},
	"96 bit overflow max": {
		[]byte{},
		"zzzzzzzzzzzzzzzzz",
		"ZZZZZZZZZZZZZZZZZZZZ",
		"",
	},
	"128 bit overflow max": {
		[]byte{},
		"zzzzzzzzzzzzzzzzzzzzzz",
		"ZZZZZZZZZZZZZZZZZZZZZZZZZZ",
		"",
	},
	"160 bit overflow max": {
		[]byte{},
		"zzzzzzzzzzzzzzzzzzzzzzzzzzz",
		"",
		"",
	},
}

var testCaseEncodeDecodeMap = map[string]testCaseEncodeDecode{
	"64 bit min": {
		data:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		base62: "00000000000",
		base32: "0000000000000",
		base16: "0000000000000000",
	},
	"96 bit min": {
		data:   []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		base62: "00000000000000000",
		base32: "00000000000000000000",
		base16: "000000000000000000000000",
	},
	"128 bit min": {
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		"0000000000000000000000",
		"00000000000000000000000000",
		"00000000000000000000000000000000",
	},
	"160 bit min": {
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		"000000000000000000000000000",
		"00000000000000000000000000000000",
		"0000000000000000000000000000000000000000",
	},
	"64 bit max": {
		[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		"LygHa16AHYF",
		"FZZZZZZZZZZZZ",
		"FFFFFFFFFFFFFFFF",
	},
	"96 bit max": {
		[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		"1f2SI9UJPXvb7vdJ1",
		"1ZZZZZZZZZZZZZZZZZZZ",
		"FFFFFFFFFFFFFFFFFFFFFFFF",
	},
	"128 bit max": {
		[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		"7n42DGM5Tflk9n8mt7Fhc7",
		"7ZZZZZZZZZZZZZZZZZZZZZZZZZ",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
	},
	"160 bit max": {
		[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		"aWgEPTl1tmebfsQzFP4bxwgy80V",
		"ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
		"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
	},
	"64 bit 0xaa": {
		[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
		"EeSBig46r2A",
		"ANANANANANANA",
		"AAAAAAAAAAAAAAAA",
	},
	"96 bit 0xaa": {
		[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
		"16gyC6KCwMcOkcQCg",
		"1ANANANANANANANANANA",
		"AAAAAAAAAAAAAAAAAAAAAAAA",
	},
	"128 bit 0xaa": {
		[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
		"5C2goAu3eRqUlrQWakAT4k",
		"5ANANANANANANANANANANANANA",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
	},
	"160 bit 0xaa": {
		[]byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
		"OLmowJq1GWR4RuxKAGiPJISe5L0",
		"NANANANANANANANANANANANANANANANA",
		"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
	},
	"64 bit 0x55": {
		[]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55},
		"7KE5rL23QW5",
		"5ANANANANANAN",
		"5555555555555555",
	},
	"96 bit 0x55": {
		[]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55},
		"0YLU63A6TBJCNJD6L",
		"0NANANANANANANANANAN",
		"555555555555555555555555",
	},
	"128 bit 0x55": {
		[]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55},
		"2b1LP5S1pDvFNviGIN5EXN",
		"2NANANANANANANANANANANANAN",
		"55555555555555555555555555555555",
	},
	"160 bit 0x55": {
		[]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55},
		"CAtPT9v0dGDXDxTf58MCeeEK2fV",
		"ANANANANANANANANANANANANANANANAN",
		"5555555555555555555555555555555555555555",
	},
	"64 bit 0xaa half": {
		[]byte{0x00, 0x00, 0x00, 0x00, 0xaa, 0xaa, 0xaa, 0xaa},
		"0000037mA82",
		"0000002NANANA",
		"00000000AAAAAAAA",
	},
	"96 bit 0xaa half": {
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
		"000000000rHgMFRik",
		"00000000005ANANANANA",
		"000000000000AAAAAAAAAAAA",
	},
	"128 bit 0xaa half": {
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
		"00000000000EeSBig46r2A",
		"0000000000000ANANANANANANA",
		"0000000000000000AAAAAAAAAAAAAAAA",
	},
	"160 bit 0xaa half": {
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
		"000000000000041o9sj4bR0Njv0",
		"0000000000000000NANANANANANANANA",
		"00000000000000000000AAAAAAAAAAAAAAAAAAAA",
	},
	"64 bit 0x55 half": {
		[]byte{0x00, 0x00, 0x00, 0x00, 0x55, 0x55, 0x55, 0x55},
		"000001Yt541",
		"0000001ANANAN",
		"0000000055555555",
	},
	"96 bit 0x55 half": {
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55},
		"000000000QdqB7irN",
		"00000000002NANANANAN",
		"000000000000555555555555",
	},
	"128 bit 0x55 half": {
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55},
		"000000000007KE5rL23QW5",
		"00000000000005ANANANANANAN",
		"00000000000000005555555555555555",
	},
	"160 bit 0x55 half": {
		[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55},
		"000000000000020u4wMXIiVBrxV",
		"0000000000000000ANANANANANANANAN",
		"0000000000000000000055555555555555555555",
	},
	"64 bit 0x55 half inverse": {
		[]byte{0x55, 0x55, 0x55, 0x55, 0x00, 0x00, 0x00, 0x00},
		"7KE5rJTALS4",
		"5ANANAM000000",
		"5555555500000000",
	},
	"96 bit 0x55 half inverse": {
		[]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		"0YLU63A6SkfMCBUEy",
		"0NANANANAN8000000000",
		"555555555555000000000000",
	},
	"128 bit 0x55 half inverse": {
		[]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		"2b1LP5S1pDv83hcOxL1o1I",
		"2NANANANANANAG000000000000",
		"55555555555555550000000000000000",
	},
	"160 bit 0x55 half inverse": {
		[]byte{0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x55, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		"CAtPT9v0dGDXDvSl0BzfLvj8Ai0",
		"ANANANANANANANAN0000000000000000",
		"5555555555555555555500000000000000000000",
	},
	"64 bit grow": {
		[]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xff},
		"1YtudU73D5z",
		"14D2PF2DBSQQZ",
		"123456789ABCDEFF",
	},
	"96 bit grow": {
		[]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xff, 0x12, 0x34, 0x56, 0x78},
		"07KHzdWk1GSjrCGrI",
		"04HMASW9NF6YZW938NKR",
		"123456789ABCDEFF12345678",
	},
	"128 bit grow": {
		[]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xff, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xff},
		"0YLmNXOzYtrdhNIZg9Xb5j",
		"0J6HB7H6NWVVZH4D2PF2DBSQQZ",
		"123456789ABCDEFF123456789ABCDEFF",
	},
	"160 bit grow": {
		[]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xff, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xff, 0x12, 0x34, 0x56, 0x78},
		"2b2j72wSEzGoPBFvlZElU0h9tkG",
		"28T5CY4TQKFFY4HMASW9NF6YZW938NKR",
		"123456789ABCDEFF123456789ABCDEFF12345678",
	},
	"64 bit grow inverse": {
		[]byte{0xff, 0xed, 0xcb, 0xa9, 0x87, 0x65, 0x43, 0x21},
		"LyIoXRNZhPl",
		"FZVEBN63PAGS1",
		"FFEDCBA987654321",
	},
	"96 bit grow inverse": {
		[]byte{0xff, 0xed, 0xcb, 0xa9, 0x87, 0x65, 0x43, 0x21, 0xff, 0xed, 0xcb, 0xa9},
		"1f0gGhTwSvqEUbOaP",
		"1ZZDSEMRESA347ZYVJX9",
		"FFEDCBA987654321FFEDCBA9",
	},
	"128 bit grow inverse": {
		[]byte{0xff, 0xed, 0xcb, 0xa9, 0x87, 0x65, 0x43, 0x21, 0xff, 0xed, 0xcb, 0xa9, 0x87, 0x65, 0x43, 0x21},
		"7mviPHR4kjQohQVJTeEupt",
		"7ZXQ5TK1V58CGZZVEBN63PAGS1",
		"FFEDCBA987654321FFEDCBA987654321",
	},
	"160 bit grow inverse": {
		[]byte{0xff, 0xed, 0xcb, 0xa9, 0x87, 0x65, 0x43, 0x21, 0xff, 0xed, 0xcb, 0xa9, 0x87, 0x65, 0x43, 0x21, 0xff, 0xed, 0xcb, 0xa9},
		"aW3EEYLc7HvSDx52u4JLRfUTFov",
		"ZZPWQAC7CN1J3ZZDSEMRESA347ZYVJX9",
		"FFEDCBA987654321FFEDCBA987654321FFEDCBA9",
	},
}

type tcAdd8 struct {
	a, b, carryIn, sum, carryOut uint8
}

type tcAdd8Slice struct {
	a, b, sum         []uint8
	carryIn, carryOut uint8
}

var testCaseByteMathMap = map[string]tcAdd8{
	"0+0+0":     {0, 0, 0, 0, 0},
	"0+0+1":     {0, 0, 1, 1, 0},
	"0+1+0":     {0, 1, 0, 1, 0},
	"0+1+1":     {0, 1, 1, 2, 0},
	"1+0+0":     {1, 0, 0, 1, 0},
	"1+0+1":     {1, 0, 1, 2, 0},
	"1+1+0":     {1, 1, 0, 2, 0},
	"1+1+1":     {1, 1, 1, 3, 0},
	"255+0+0":   {255, 0, 0, 255, 0},
	"255+0+1":   {255, 0, 1, 0, 1},
	"255+1+0":   {255, 1, 0, 0, 1},
	"255+1+1":   {255, 1, 1, 1, 1},
	"255+255+0": {255, 255, 0, 254, 1},
	"255+255+1": {255, 255, 1, 255, 1},
}

var testCaseAdd8SliceMap = map[string]tcAdd8Slice{
	"empty+empty+0":      {[]uint8{}, []uint8{}, []uint8{}, 0, 0},
	"empty+empty+1":      {[]uint8{}, []uint8{}, []uint8{1}, 1, 0},
	"empty+1+0":          {[]uint8{}, []uint8{1}, []uint8{1}, 0, 0},
	"empty+1+1":          {[]uint8{}, []uint8{1}, []uint8{2}, 1, 0},
	"1+empty+0":          {[]uint8{1}, []uint8{}, []uint8{1}, 0, 0},
	"1+empty+1":          {[]uint8{1}, []uint8{}, []uint8{2}, 1, 0},
	"0+0+0":              {[]uint8{0}, []uint8{0}, []uint8{0}, 0, 0},
	"0+0+1":              {[]uint8{0}, []uint8{0}, []uint8{1}, 1, 0},
	"0+1+0":              {[]uint8{0}, []uint8{1}, []uint8{1}, 0, 0},
	"0+1+1":              {[]uint8{0}, []uint8{1}, []uint8{2}, 1, 0},
	"1+0+0":              {[]uint8{1}, []uint8{0}, []uint8{1}, 0, 0},
	"1+0+1":              {[]uint8{1}, []uint8{0}, []uint8{2}, 1, 0},
	"1+1+0":              {[]uint8{1}, []uint8{1}, []uint8{2}, 0, 0},
	"1+1+1":              {[]uint8{1}, []uint8{1}, []uint8{3}, 1, 0},
	"255+0+0":            {[]uint8{255}, []uint8{0}, []uint8{255}, 0, 0},
	"255+0+1":            {[]uint8{255}, []uint8{0}, []uint8{0}, 1, 1},
	"255+1+0":            {[]uint8{255}, []uint8{1}, []uint8{0}, 0, 1},
	"255+1+1":            {[]uint8{255}, []uint8{1}, []uint8{1}, 1, 1},
	"255+255+0":          {[]uint8{255}, []uint8{255}, []uint8{254}, 0, 1},
	"255+255+1":          {[]uint8{255}, []uint8{255}, []uint8{255}, 1, 1},
	"sliceGrow+0":        {[]uint8{2, 1}, []uint8{20, 10}, []uint8{22, 11}, 0, 0},
	"sliceGrow+1":        {[]uint8{2, 1}, []uint8{20, 10}, []uint8{22, 12}, 1, 0},
	"slice3x255+0":       {[]uint8{255, 255, 255}, []uint8{255, 255, 255}, []uint8{255, 255, 254}, 0, 1},
	"slice3x255+1":       {[]uint8{255, 255, 255}, []uint8{255, 255, 255}, []uint8{255, 255, 255}, 1, 1},
	"slice3-2x255+0":     {[]uint8{255, 255, 255}, []uint8{255, 255}, []uint8{0, 255, 254}, 0, 1},
	"slice3-2x255+1":     {[]uint8{255, 255, 255}, []uint8{255, 255}, []uint8{0, 255, 255}, 1, 1},
	"slice1-3x255+0":     {[]uint8{255}, []uint8{255, 255, 255}, []uint8{0, 0, 254}, 0, 1},
	"slice1-3x255+1":     {[]uint8{255}, []uint8{255, 255, 255}, []uint8{0, 0, 255}, 1, 1},
	"empty+slice3+0":     {[]uint8{}, []uint8{32, 64, 128}, []uint8{32, 64, 128}, 0, 0},
	"empty+slice3+1":     {[]uint8{}, []uint8{32, 64, 128}, []uint8{32, 64, 129}, 1, 0},
	"empty+slice3x255+1": {[]uint8{}, []uint8{255, 255, 255}, []uint8{0, 0, 0}, 1, 1},
}

type CcIdTestCases struct {
	name        string
	timestamp   uint32
	time        time.Time
	Fingerprint []byte
	payload     []byte
	Bytes       []byte
	Base62      string
	Base32      string
	Base16      string
	GoString    string
}

type CcIdTestErrorCases struct {
	name        string
	timestamp   uint32
	fingerprint []byte
	payload     []byte
	bytes       []byte
	errMsg      string
}

var testCaseCcId64ErrorMap = map[string]CcIdTestErrorCases{
	"small payload error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: nil,
		payload:     []byte{1, 2, 3},
		bytes:       NilCcId64.Bytes(),
		errMsg:      "CCID: invalid payload length 3 bytes, required min 4 bytes",
	},
	"finger print error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: []byte{0x12, 0x34, 0x56, 0x78, 0x9a},
		payload:     []byte{1, 2, 3, 4},
		bytes:       NilCcId64.Bytes(),
		errMsg:      "CCID: invalid fingerprint length 5 bytes, required max 1 bytes",
	},
	"finger print small payload error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: []byte{0x12},
		payload:     []byte{1, 2},
		bytes:       NilCcId64.Bytes(),
		errMsg:      "CCID: invalid payload length 2 bytes, required min 3 bytes",
	},
}

var TestCaseCcId64Map = map[string]CcIdTestCases{
	"min id": {
		timestamp:   0,
		time:        time.Date(2014, 05, 13, 16, 53, 20, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{1, 2, 3, 4},
		Bytes:       []byte{0, 0, 0, 0, 1, 2, 3, 4},
		Base62:      "00000018wom",
		Base32:      "00000000G40R4",
		Base16:      "0000000001020304",
		GoString:    "CcId{size: 8, timestamp: 0 (2014-05-13T16:53:20Z), payload: 0x01020304}",
	},
	"some id": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{0x12, 0x34, 0x56, 0x78},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x12, 0x34, 0x56, 0x78},
		Base62:      "1YtudRc1sam",
		Base32:      "14D2PF0938NKR",
		Base16:      "1234567812345678",
		GoString:    "CcId{size: 8, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x12345678}",
	},
	"some id fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0x99},
		payload:     []byte{0x12, 0x34, 0x56},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x99, 0x12, 0x34, 0x56},
		Base62:      "1YtudU59sta",
		Base32:      "14D2PF2CH4D2P",
		Base16:      "1234567899123456",
		GoString:    "CcId{size: 8, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0x99, payload: 0x123456}",
	},
	"max id": {
		timestamp:   0xFFFFFFFF,
		time:        time.Date(2150, 06, 19, 23, 21, 35, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{0xFF, 0xFF, 0xFF, 0xFF},
		Bytes:       []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		Base62:      "LygHa16AHYF",
		Base32:      "FZZZZZZZZZZZZ",
		Base16:      "FFFFFFFFFFFFFFFF",
		GoString:    "CcId{size: 8, timestamp: 4294967295 (2150-06-19T23:21:35Z), payload: 0xffffffff}",
	},
	"large payload": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 1, 2, 3, 4},
		Base62:      "1YtudRIVJkC",
		Base32:      "14D2PF00G40R4",
		Base16:      "1234567801020304",
		GoString:    "CcId{size: 8, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x01020304}",
	},
	"empty fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{},
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 1, 2, 3, 4},
		Base62:      "1YtudRIVJkC",
		Base32:      "14D2PF00G40R4",
		Base16:      "1234567801020304",
		GoString:    "CcId{size: 8, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x01020304}",
	},
	"fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0x99},
		payload:     []byte{1, 2, 3},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x99, 1, 2, 3},
		Base62:      "1YtudU559iF",
		Base32:      "14D2PF2CG20G3",
		Base16:      "1234567899010203",
		GoString:    "CcId{size: 8, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0x99, payload: 0x010203}",
	},
	"fingerprint large payload": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0x99},
		payload:     []byte{1, 2, 3, 4, 5},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x99, 1, 2, 3},
		Base62:      "1YtudU559iF",
		Base32:      "14D2PF2CG20G3",
		Base16:      "1234567899010203",
		GoString:    "CcId{size: 8, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0x99, payload: 0x010203}",
	},
}

var testCaseCcId96ErrorMap = map[string]CcIdTestErrorCases{
	"small payload error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: nil,
		payload:     []byte{1, 2, 3, 4, 5, 6, 7},
		bytes:       NilCcId96.Bytes(),
		errMsg:      "CCID: invalid payload length 7 bytes, required min 8 bytes",
	},
	"finger print error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc},
		payload:     []byte{1, 2, 3, 4},
		bytes:       NilCcId96.Bytes(),
		errMsg:      "CCID: invalid fingerprint length 6 bytes, required max 5 bytes",
	},
	"finger print small payload error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: []byte{0x12, 0x34, 0x56, 0x78, 0x9a},
		payload:     []byte{1, 2},
		bytes:       NilCcId96.Bytes(),
		errMsg:      "CCID: invalid payload length 2 bytes, required min 3 bytes",
	},
}

var TestCaseCcId96Map = map[string]CcIdTestCases{
	"min id": {
		timestamp:   0,
		time:        time.Date(2014, 05, 13, 16, 53, 20, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
		Bytes:       []byte{0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8},
		Base62:      "00000005McJmDgrvc",
		Base32:      "0000000020G30G2GC1R8",
		Base16:      "000000000102030405060708",
		GoString:    "CcId{size: 12, timestamp: 0 (2014-05-13T16:53:20Z), payload: 0x0102030405060708}",
	},
	"some id": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88},
		Base62:      "07KHzdKvYVTUKifX6",
		Base32:      "04HMASW128HK8HAPCXW8",
		Base16:      "123456781122334455667788",
		GoString:    "CcId{size: 12, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x1122334455667788}",
	},
	"some id fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0x99, 0xaa, 0xbb},
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x99, 0xaa, 0xbb, 0x11, 0x22, 0x33, 0x44, 0x55},
		Base62:      "07KHzdWeJr1R43LQj",
		Base32:      "04HMASW9KANV24H36H2N",
		Base16:      "1234567899AABB1122334455",
		GoString:    "CcId{size: 12, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0x99aabb, payload: 0x1122334455}",
	},
	"max id": {
		timestamp:   0xFFFFFFFF,
		time:        time.Date(2150, 06, 19, 23, 21, 35, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		Bytes:       []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		Base62:      "1f2SI9UJPXvb7vdJ1",
		Base32:      "1ZZZZZZZZZZZZZZZZZZZ",
		Base16:      "FFFFFFFFFFFFFFFFFFFFFFFF",
		GoString:    "CcId{size: 12, timestamp: 4294967295 (2150-06-19T23:21:35Z), payload: 0xffffffffffffffff}",
	},
	"large payload": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 1, 2, 3, 4, 5, 6, 7, 8},
		Base62:      "07KHzdJXicN2nekfI",
		Base32:      "04HMASW020G30G2GC1R8",
		Base16:      "123456780102030405060708",
		GoString:    "CcId{size: 12, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x0102030405060708}",
	},
	"empty fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{},
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 1, 2, 3, 4, 5, 6, 7, 8},
		Base62:      "07KHzdJXicN2nekfI",
		Base32:      "04HMASW020G30G2GC1R8",
		Base16:      "123456780102030405060708",
		GoString:    "CcId{size: 12, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x0102030405060708}",
	},
	"min fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0x99},
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x99, 1, 2, 3, 4, 5, 6, 7},
		Base62:      "07KHzdWan3RXAOmU3",
		Base32:      "04HMASW9J0820C20A1G7",
		Base16:      "123456789901020304050607",
		GoString:    "CcId{size: 12, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0x99, payload: 0x01020304050607}",
	},
	"max fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0x99, 0x88, 0x77, 0x66, 0x55},
		payload:     []byte{1, 2, 3, 4, 5},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x99, 0x88, 0x77, 0x66, 0x55, 1, 2, 3},
		Base62:      "07KHzdWdbgLit1TQR",
		Base32:      "04HMASW9K23QCSAG20G3",
		Base16:      "123456789988776655010203",
		GoString:    "CcId{size: 12, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0x9988776655, payload: 0x010203}",
	},
	"fingerprint large payload": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0x99, 0x88, 0x77},
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x99, 0x88, 0x77, 1, 2, 3, 4, 5},
		Base62:      "07KHzdWdbgE3rbRor",
		Base32:      "04HMASW9K23Q04106105",
		Base16:      "123456789988770102030405",
		GoString:    "CcId{size: 12, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0x998877, payload: 0x0102030405}",
	},
}

var testCaseCcId128ErrorMap = map[string]CcIdTestErrorCases{
	"small payload error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: nil,
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		bytes:       NilCcId128.Bytes(),
		errMsg:      "CCID: invalid payload length 11 bytes, required min 12 bytes",
	},
	"finger print error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc},
		payload:     []byte{1, 2, 3, 4, 5, 6, 7},
		bytes:       NilCcId128.Bytes(),
		errMsg:      "CCID: invalid fingerprint length 6 bytes, required max 5 bytes",
	},
	"finger print small payload error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: []byte{0x12, 0x34, 0x56, 0x78, 0x9a},
		payload:     []byte{1, 2, 3, 4, 5, 6},
		bytes:       NilCcId128.Bytes(),
		errMsg:      "CCID: invalid payload length 6 bytes, required min 7 bytes",
	},
}

var TestCaseCcId128Map = map[string]CcIdTestCases{
	"min id": {
		timestamp:   0,
		time:        time.Date(2014, 05, 13, 16, 53, 20, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		Bytes:       []byte{0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		Base62:      "0000000P9MVMcaXWHEX8yi",
		Base32:      "00000000820C20A1G7104GM2RC",
		Base16:      "000000000102030405060708090A0B0C",
		GoString:    "CcId{size: 16, timestamp: 0 (2014-05-13T16:53:20Z), payload: 0x0102030405060708090a0b0c}",
	},
	"some id": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc},
		Base62:      "0YLmNWVbf7YaOQib5nMzim",
		Base32:      "0J6HB7G4926D25ASKQH2CTNEYC",
		Base16:      "12345678112233445566778899AABBCC",
		GoString:    "CcId{size: 16, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x112233445566778899aabbcc}",
	},
	"some id fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0xdd, 0xee, 0xff},
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
		Base62:      "0YLmNXq2QX11gZMcdEr1iL",
		Base32:      "0J6HB7HQFEZW8J4CT4ANK7F24S",
		Base16:      "12345678DDEEFF112233445566778899",
		GoString:    "CcId{size: 16, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0xddeeff, payload: 0x112233445566778899}",
	},
	"max id": {
		timestamp:   0xFFFFFFFF,
		time:        time.Date(2150, 06, 19, 23, 21, 35, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		Bytes:       []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		Base62:      "7n42DGM5Tflk9n8mt7Fhc7",
		Base32:      "7ZZZZZZZZZZZZZZZZZZZZZZZZZ",
		Base16:      "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		GoString:    "CcId{size: 16, timestamp: 4294967295 (2150-06-19T23:21:35Z), payload: 0xffffffffffffffffffffffff}",
	},
	"large payload": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc},
		Base62:      "0YLmNWVbf7YaOQib5nMzim",
		Base32:      "0J6HB7G4926D25ASKQH2CTNEYC",
		Base16:      "12345678112233445566778899AABBCC",
		GoString:    "CcId{size: 16, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x112233445566778899aabbcc}",
	},
	"empty fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{},
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc},
		Base62:      "0YLmNWVbf7YaOQib5nMzim",
		Base32:      "0J6HB7G4926D25ASKQH2CTNEYC",
		Base16:      "12345678112233445566778899AABBCC",
		GoString:    "CcId{size: 16, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x112233445566778899aabbcc}",
	},
	"min fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0xdd},
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0xdd, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb},
		Base62:      "0YLmNXpgne2KU0qObqHPwp",
		Base32:      "0J6HB7HQ8H48SM8NB6EY49KANV",
		Base16:      "12345678DD112233445566778899AABB",
		GoString:    "CcId{size: 16, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0xdd, payload: 0x112233445566778899aabb}",
	},
	"max fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0xbb, 0xcc, 0xdd, 0xee, 0xff},
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77},
		Base62:      "0YLmNXcIdJeuW33h8kJikR",
		Base32:      "0J6HB7HEYCVQQFY4926D25ASKQ",
		Base16:      "12345678BBCCDDEEFF11223344556677",
		GoString:    "CcId{size: 16, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0xbbccddeeff, payload: 0x11223344556677}",
	},
	"fingerprint large payload": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0xdd, 0xee, 0xff},
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99},
		Base62:      "0YLmNXq2QX11gZMcdEr1iL",
		Base32:      "0J6HB7HQFEZW8J4CT4ANK7F24S",
		Base16:      "12345678DDEEFF112233445566778899",
		GoString:    "CcId{size: 16, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0xddeeff, payload: 0x112233445566778899}",
	},
}

var testCaseCcId160ErrorMap = map[string]CcIdTestErrorCases{
	"small payload error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: nil,
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		bytes:       NilCcId160.Bytes(),
		errMsg:      "CCID: invalid payload length 15 bytes, required min 16 bytes",
	},
	"finger print error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc},
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		bytes:       NilCcId160.Bytes(),
		errMsg:      "CCID: invalid fingerprint length 6 bytes, required max 5 bytes",
	},
	"finger print small payload error": {
		timestamp:   0xf0e0d0c0,
		fingerprint: []byte{0x12, 0x34, 0x56, 0x78, 0x9a},
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		bytes:       NilCcId160.Bytes(),
		errMsg:      "CCID: invalid payload length 10 bytes, required min 11 bytes",
	},
}

var TestCaseCcId160Map = map[string]CcIdTestCases{
	"min id": {
		timestamp:   0,
		time:        time.Date(2014, 05, 13, 16, 53, 20, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		Bytes:       []byte{0, 0, 0, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		Base62:      "0000001tuWckR0Qgud2DqqiTysq",
		Base32:      "00000001081G81860W40J2GB1G6GW3RG",
		Base16:      "000000000102030405060708090A0B0C0D0E0F10",
		GoString:    "CcId{size: 20, timestamp: 0 (2014-05-13T16:53:20Z), payload: 0x0102030405060708090a0b0c0d0e0f10}",
	},
	"some id": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x01},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x01},
		Base62:      "2b2j6yknbhs3Jp7SXB5fMnkzNmj",
		Base32:      "28T5CY0H48SM8NB6EY49KANVSKEYXZR1",
		Base16:      "12345678112233445566778899AABBCCDDEEFF01",
		GoString:    "CcId{size: 20, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x112233445566778899aabbccddeeff01}",
	},
	"some id fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0xdd, 0xee, 0xff},
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd},
		Base62:      "2b2j74zFdMosqaItZ9RHloJDeA9",
		Base32:      "28T5CY6XXVZH28HK8HAPCXW8K6NBQK6X",
		Base16:      "12345678DDEEFF112233445566778899AABBCCDD",
		GoString:    "CcId{size: 20, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0xddeeff, payload: 0x112233445566778899aabbccdd}",
	},
	"max id": {
		timestamp:   0xFFFFFFFF,
		time:        time.Date(2150, 06, 19, 23, 21, 35, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		Bytes:       []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		Base62:      "aWgEPTl1tmebfsQzFP4bxwgy80V",
		Base32:      "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ",
		Base16:      "FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		GoString:    "CcId{size: 20, timestamp: 4294967295 (2150-06-19T23:21:35Z), payload: 0xffffffffffffffffffffffffffffffff}",
	},
	"large payload": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: nil,
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x01, 0x02},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x01},
		Base62:      "2b2j6yknbhs3Jp7SXB5fMnkzNmj",
		Base32:      "28T5CY0H48SM8NB6EY49KANVSKEYXZR1",
		Base16:      "12345678112233445566778899AABBCCDDEEFF01",
		GoString:    "CcId{size: 20, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x112233445566778899aabbccddeeff01}",
	},
	"empty fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{},
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x01},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x01},
		Base62:      "2b2j6yknbhs3Jp7SXB5fMnkzNmj",
		Base32:      "28T5CY0H48SM8NB6EY49KANVSKEYXZR1",
		Base16:      "12345678112233445566778899AABBCCDDEEFF01",
		GoString:    "CcId{size: 20, timestamp: 305419896 (2024-01-16T15:44:56Z), payload: 0x112233445566778899aabbccddeeff01}",
	},
	"min fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0xff},
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0xff, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		Base62:      "2b2j75zlyuoPcliOOIZr8T6RIcR",
		Base32:      "28T5CY7Z24H36H2NCSVRH6DAQF6DVVQZ",
		Base16:      "12345678FF112233445566778899AABBCCDDEEFF",
		GoString:    "CcId{size: 20, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0xff, payload: 0x112233445566778899aabbccddeeff}",
	},
	"max fingerprint": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0xbb, 0xcc, 0xdd, 0xee, 0xff},
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb},
		Base62:      "2b2j73wqIW2wyWMdozma2887Rpj",
		Base32:      "28T5CY5VSKEYXZRH48SM8NB6EY49KANV",
		Base16:      "12345678BBCCDDEEFF112233445566778899AABB",
		GoString:    "CcId{size: 20, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0xbbccddeeff, payload: 0x112233445566778899aabb}",
	},
	"fingerprint large payload": {
		timestamp:   0x12345678,
		time:        time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		Fingerprint: []byte{0xdd, 0xee, 0xff},
		payload:     []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff},
		Bytes:       []byte{0x12, 0x34, 0x56, 0x78, 0xdd, 0xee, 0xff, 0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb, 0xcc, 0xdd},
		Base62:      "2b2j74zFdMosqaItZ9RHloJDeA9",
		Base32:      "28T5CY6XXVZH28HK8HAPCXW8K6NBQK6X",
		Base16:      "12345678DDEEFF112233445566778899AABBCCDD",
		GoString:    "CcId{size: 20, timestamp: 305419896 (2024-01-16T15:44:56Z), fingerprint: 0xddeeff, payload: 0x112233445566778899aabbccdd}",
	},
}

func SortKeys[T any](m map[string]T) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func SliceEqual(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
