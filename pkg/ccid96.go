package pkg

import (
	"encoding/binary"
	"fmt"
	"time"
)

var NilCcId96 CcId96

type CcId96 struct {
	fingerprintSize byte
	data            [ByteSliceSize96]byte
}

func (id CcId96) Size() byte {
	return ByteSliceSize96
}

func (id CcId96) Time() time.Time {
	return ToStandardizedTime(id.Timestamp())
}

func (id CcId96) Timestamp() uint32 {
	return binary.BigEndian.Uint32(id.data[:TimestampSize])
}

func (id CcId96) Fingerprint() []byte {
	return id.data[TimestampSize : TimestampSize+id.fingerprintSize]
}

func (id CcId96) Payload() []byte {
	return id.data[TimestampSize+id.fingerprintSize:]
}

func (id CcId96) String() string {
	return id.AsBase62()
}

func (id CcId96) GoString() string {
	if id.fingerprintSize > 0 {
		return fmt.Sprintf("CcId{size: %d, timestamp: %d (%s), fingerprint: 0x%x, payload: 0x%x}",
			id.Size(), id.Timestamp(), id.Time().Format(time.RFC3339), id.Fingerprint(), id.Payload())
	}
	return fmt.Sprintf("CcId{size: %d, timestamp: %d (%s), payload: 0x%x}",
		id.Size(), id.Timestamp(), id.Time().Format(time.RFC3339), id.Payload())
}

func (id CcId96) Bytes() []byte {
	return id.data[:]
}

func (id CcId96) AsBase62() string {
	v, _ := EncodeToBase62(id.data[:])
	return v
}

func (id CcId96) AsBase32() string {
	v, _ := EncodeToBase32(id.data[:])
	return v
}

func (id CcId96) AsBase16() string {
	v, _ := EncodeToBase16(id.data[:])
	return v
}

func NewCcId96WithFingerprint(timestamp uint32, fingerprint []byte, payload []byte) (CcId, error) {
	var id CcId96
	payloadSize := byte(len(payload))
	fingerprintSize := byte(len(fingerprint))
	if fingerprintSize > MaxFingerprintSize {
		return NilCcId96, InvalidFingerprintSizeError{
			ProvidedSize: fingerprintSize,
			RequiredSize: MaxFingerprintSize,
		}
	}
	minPayloadSize := ByteSliceSize96 - TimestampSize - fingerprintSize
	if payloadSize < minPayloadSize {
		return NilCcId96, InvalidPayloadSizeError{
			ProvidedSize: payloadSize,
			RequiredSize: minPayloadSize,
		}
	}
	binary.BigEndian.PutUint32(id.data[:TimestampSize], timestamp)
	if fingerprintSize > 0 {
		id.fingerprintSize = fingerprintSize
		copy(id.data[TimestampSize:TimestampSize+fingerprintSize], fingerprint[:])
	}
	copy(id.data[TimestampSize+id.fingerprintSize:], payload[:])
	return id, nil
}
