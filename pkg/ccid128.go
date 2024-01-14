package pkg

import (
	"encoding/binary"
	"fmt"
	"time"
)

var NilCcId128 CcId128

type CcId128 struct {
	fingerprintSize byte
	data            [ByteSliceSize128]byte
}

func (id CcId128) Size() byte {
	return ByteSliceSize128
}

func (id CcId128) Time() time.Time {
	return ToStandardizedTime(id.Timestamp())
}

func (id CcId128) Timestamp() uint32 {
	return binary.BigEndian.Uint32(id.data[:TimestampSize])
}

func (id CcId128) Fingerprint() []byte {
	return id.data[TimestampSize : TimestampSize+id.fingerprintSize]
}

func (id CcId128) Payload() []byte {
	return id.data[TimestampSize+id.fingerprintSize:]
}

func (id CcId128) String() string {
	return id.AsBase62()
}

func (id CcId128) GoString() string {
	if id.fingerprintSize > 0 {
		return fmt.Sprintf("CcId{size: %d, timestamp: %d (%s), fingerprint: 0x%x, payload: 0x%x}",
			id.Size(), id.Timestamp(), id.Time().Format(time.RFC3339), id.Fingerprint(), id.Payload())
	}
	return fmt.Sprintf("CcId{size: %d, timestamp: %d (%s), payload: 0x%x}",
		id.Size(), id.Timestamp(), id.Time().Format(time.RFC3339), id.Payload())
}

func (id CcId128) Bytes() []byte {
	return id.data[:]
}

func (id CcId128) AsBase62() string {
	v, _ := EncodeToBase62(id.data[:])
	return v
}

func (id CcId128) AsBase32() string {
	v, _ := EncodeToBase32(id.data[:])
	return v
}

func (id CcId128) AsBase16() string {
	v, _ := EncodeToBase16(id.data[:])
	return v
}

func NewCcId128WithFingerprint(timestamp uint32, fingerprint []byte, payload []byte) (CcId, error) {
	var id CcId128
	payloadSize := byte(len(payload))
	fingerprintSize := byte(len(fingerprint))
	if fingerprintSize > MaxFingerprintSize {
		return NilCcId128, InvalidFingerprintSizeError{
			ProvidedSize: fingerprintSize,
			RequiredSize: MaxFingerprintSize,
		}
	}
	minPayloadSize := ByteSliceSize128 - TimestampSize - fingerprintSize
	if payloadSize < minPayloadSize {
		return NilCcId128, InvalidPayloadSizeError{
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
