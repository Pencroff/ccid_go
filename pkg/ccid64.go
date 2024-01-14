package pkg

import (
	"encoding/binary"
	"fmt"
	"time"
)

var NilCcId64 CcId64

type CcId64 struct {
	fingerprintSize byte
	data            [ByteSliceSize64]byte
}

func (id CcId64) Size() byte {
	return ByteSliceSize64
}

func (id CcId64) Time() time.Time {
	return ToStandardizedTime(id.Timestamp())
}

func (id CcId64) Timestamp() uint32 {
	return binary.BigEndian.Uint32(id.data[:TimestampSize])
}

func (id CcId64) Fingerprint() []byte {
	return id.data[TimestampSize : TimestampSize+id.fingerprintSize]
}

func (id CcId64) Payload() []byte {
	return id.data[TimestampSize+id.fingerprintSize:]
}

func (id CcId64) String() string {
	return id.AsBase62()
}

func (id CcId64) GoString() string {
	if id.fingerprintSize > 0 {
		return fmt.Sprintf("CcId{size: %d, timestamp: %d (%s), fingerprint: 0x%x, payload: 0x%x}",
			id.Size(), id.Timestamp(), id.Time().Format(time.RFC3339), id.Fingerprint(), id.Payload())
	}
	return fmt.Sprintf("CcId{size: %d, timestamp: %d (%s), payload: 0x%x}",
		id.Size(), id.Timestamp(), id.Time().Format(time.RFC3339), id.Payload())
}

func (id CcId64) Bytes() []byte {
	return id.data[:]
}

func (id CcId64) AsBase62() string {
	v, _ := EncodeToBase62(id.data[:])
	return v
}

func (id CcId64) AsBase32() string {
	v, _ := EncodeToBase32(id.data[:])
	return v
}

func (id CcId64) AsBase16() string {
	v, _ := EncodeToBase16(id.data[:])
	return v
}

func (id CcId64) Uint64() uint64 {
	return binary.BigEndian.Uint64(id.data[:])
}

func NewCcId64WithFingerprint(timestamp uint32, fingerprint []byte, payload []byte) (CcId, error) {
	var id CcId64
	payloadSize := byte(len(payload))
	fingerprintSize := byte(len(fingerprint))
	if fingerprintSize > MaxFingerprintSize64 {
		return NilCcId64, InvalidFingerprintSizeError{
			ProvidedSize: fingerprintSize,
			RequiredSize: MaxFingerprintSize64,
		}
	}
	minPayloadSize := ByteSliceSize64 - TimestampSize - fingerprintSize
	if payloadSize < minPayloadSize {
		return NilCcId64, InvalidPayloadSizeError{
			ProvidedSize: payloadSize,
			RequiredSize: minPayloadSize,
		}
	}
	binary.BigEndian.PutUint32(id.data[:TimestampSize], timestamp)
	if fingerprintSize > 0 {
		id.fingerprintSize = MaxFingerprintSize64
		id.data[TimestampSize] = fingerprint[0]
	}
	copy(id.data[TimestampSize+id.fingerprintSize:], payload[:])
	return id, nil
}
