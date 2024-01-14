package pkg

import (
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"testing"
)

var testList = []string{
	//"64 bit min",
	//"160 bit min",
	//"64 bit grow",
	"160 bit grow",
}

const allTests = false

func BenchmarkEncodeBase16(b *testing.B) {
	for _, tcName := range testList {
		tc := testCaseEncodeDecodeMap[tcName].data
		size := len(tc)
		if allTests {
			b.Run(tcName+"_BigInt", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					AsBase16BigInt(tc)
				}
				b.SetBytes(int64(size))
			})
			//// same as BigInt
			//b.Run(tcName+"_Sprintf", func(b *testing.B) {
			//	for i := 0; i < b.N; i++ {
			//		AsBase16BigInt(tc)
			//	}
			//	b.SetBytes(int64(size))
			//})
			// incorrect str register
			b.Run(tcName+"_hex", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					hex.EncodeToString(tc)
				}
				b.SetBytes(int64(size))
			})
		}
		b.Run(tcName, func(bb *testing.B) {
			for i := 0; i < bb.N; i++ {
				EncodeToBase16(tc)
			}
			bb.SetBytes(int64(size))
		})
		fmt.Printf("---\n")
	}
}

func BenchmarkDecodeBase16(b *testing.B) {
	for _, tcName := range testList {
		tc := testCaseEncodeDecodeMap[tcName]
		size := len(tc.data)
		if allTests {
			b.Run(tcName+"_hex", func(bb *testing.B) {
				for i := 0; i < bb.N; i++ {
					hex.DecodeString(tc.base16)
				}
				bb.SetBytes(int64(size))

			})
		}
		b.Run(tcName, func(bb *testing.B) {
			for i := 0; i < bb.N; i++ {
				DecodeFromBase16(tc.base16)
			}
			bb.SetBytes(int64(size))
		})
		fmt.Printf("---\n")
	}
}

func BenchmarkEncodeBase32(b *testing.B) {
	std := base32.NewEncoding(base32Alphabet).WithPadding(base32.NoPadding)
	for _, tcName := range testList {
		tc := testCaseEncodeDecodeMap[tcName].data
		size := len(tc)
		if allTests {
			b.Run(tcName+"_BigInt", func(bb *testing.B) {
				for i := 0; i < bb.N; i++ {
					AsBase32BigInt(tc)
				}
				bb.SetBytes(int64(size))
			})
			b.Run(tcName+"_std", func(bb *testing.B) {
				for i := 0; i < bb.N; i++ {
					std.EncodeToString(tc)
				}
				bb.SetBytes(int64(size))
			})
		}
		b.Run(tcName, func(bb *testing.B) {
			for i := 0; i < bb.N; i++ {
				EncodeToBase32(tc)
			}
			bb.SetBytes(int64(size))
		})
		fmt.Printf("---\n")
	}
}

func BenchmarkDecodeBase32(b *testing.B) {
	std := base32.NewEncoding(base32Alphabet).WithPadding(base32.NoPadding)
	for _, tcName := range testList {
		tc := testCaseEncodeDecodeMap[tcName]
		size := len(tc.data)
		if allTests {
			b.Run(tcName+"_std", func(bb *testing.B) {
				for i := 0; i < bb.N; i++ {
					std.DecodeString(tc.base32)
				}
				bb.SetBytes(int64(size))
			})
		}
		b.Run(tcName, func(bb *testing.B) {
			for i := 0; i < bb.N; i++ {
				DecodeFromBase32(tc.base32)
			}
			bb.SetBytes(int64(size))
		})
		fmt.Printf("---\n")
	}
}

func BenchmarkEncodeBase62(b *testing.B) {
	for _, tcName := range testList {
		tc := testCaseEncodeDecodeMap[tcName].data
		size := len(tc)
		if allTests {
			b.Run(tcName+"_BigInt", func(bb *testing.B) {
				for i := 0; i < bb.N; i++ {
					AsBase62BigInt(tc)
				}
				bb.SetBytes(int64(size))
			})
		}
		b.Run(tcName, func(bb *testing.B) {
			for i := 0; i < bb.N; i++ {
				EncodeToBase62(tc)
			}
			bb.SetBytes(int64(size))
		})
		fmt.Printf("---\n")
	}
}

func BenchmarkDecodeBase62(b *testing.B) {
	for _, tcName := range testList {
		tc := testCaseEncodeDecodeMap[tcName]
		size := len(tc.data)
		b.Run(tcName, func(bb *testing.B) {
			for i := 0; i < bb.N; i++ {
				DecodeFromBase62(tc.base62)
			}
			bb.SetBytes(int64(size))
		})
		fmt.Printf("---\n")
	}
}
