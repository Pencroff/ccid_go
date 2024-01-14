package pkg

import (
	"fmt"
	"time"
)

const (
	TimestampSize = 4

	MaxFingerprintSize64 = 1
	MaxFingerprintSize   = 5

	ByteSliceSize64  = 8
	ByteSliceSize96  = 12
	ByteSliceSize128 = 16
	ByteSliceSize160 = 20

	zeroString = "0000000000000000000000000000000000000000"

	epochStamp int64 = 1400000000

	BASE62 = 62
	BASE32 = 32
	BASE16 = 16
)

// CcId represents a unique identifier with various properties.
// It provides methods for retrieving these properties and for exporting the identifier in different formats.
type CcId interface {
	// Size returns the size of the CcId.
	Size() byte
	// Time returns the time associated with the CcId.
	Time() time.Time
	// Timestamp returns the timestamp of the CcId. Unix epoch in seconds.
	Timestamp() uint32
	// Fingerprint returns the fingerprint of the CcId.
	Fingerprint() []byte
	// Payload returns random payload of the CcId.
	Payload() []byte
	// String returns the CcId as a base62 string.
	String() string
	// Bytes returns the CcId as a raw byte slice.
	Bytes() []byte
	// AsBase62 returns the CcId as a base62 string.
	AsBase62() string // (Encode: 95 MB/s / Decode: 155 MB/s)
	// AsBase32 returns the CcId as a base32 string.
	AsBase32() string // (Encode: 425 MB/s / Decode: 365 MB/s)
	// AsBase16 returns the CcId as a base16 string.
	AsBase16() string // (Encode: 680 MB/s / Decode: 490 MB/s)
}

type CcIdCtor func(timestamp uint32, fingerprint []byte, payload []byte) (CcId, error)

type InvalidPayloadSizeError struct {
	ProvidedSize byte
	RequiredSize byte
}

func (e InvalidPayloadSizeError) Error() string {
	return fmt.Sprintf("CCID: invalid payload length %d bytes, required min %d bytes", e.ProvidedSize, e.RequiredSize)
}

type InvalidFingerprintSizeError struct {
	ProvidedSize byte
	RequiredSize byte
}

func (e InvalidFingerprintSizeError) Error() string {
	return fmt.Sprintf("CCID: invalid fingerprint length %d bytes, required max %d bytes", e.ProvidedSize, e.RequiredSize)

}

type OverflowError byte

func (e OverflowError) Error() string {
	return fmt.Sprintf("CCID: decode overflow, puvided value not representable by %d bytes slice", byte(e))
}

type InvalidLengthError byte

func (e InvalidLengthError) Error() string {
	return fmt.Sprintf("CCID: invalid length %d bytes", byte(e))
}

type InvalidCharacterError struct {
	Character byte
	Pos       uint8
}

func (e InvalidCharacterError) Error() string {
	return fmt.Sprintf("CCID: invalid character %q at position %d", e.Character, e.Pos)
}

// Clock is an interface for getting current time.
// It is used for mocking time in tests.
type Clock interface {
	Now() time.Time
}

// RealClock is a real time clock.
type RealClock struct{}

// Now returns current time.
func (RealClock) Now() time.Time {
	return time.Now()
}
