package pkg

import (
	"encoding/binary"
	"fmt"
	"time"
)

var NilCcId160 CcId160

type CcId160 struct {
	fingerprintSize byte
	data            [ByteSliceSize160]byte
}

func (id CcId160) Size() byte {
	return ByteSliceSize160
}

func (id CcId160) Time() time.Time {
	return ToStandardizedTime(id.Timestamp())
}

func (id CcId160) Timestamp() uint32 {
	return binary.BigEndian.Uint32(id.data[:TimestampSize])
}

func (id CcId160) Fingerprint() []byte {
	return id.data[TimestampSize : TimestampSize+id.fingerprintSize]
}

func (id CcId160) Payload() []byte {
	return id.data[TimestampSize+id.fingerprintSize:]
}

func (id CcId160) String() string {
	return id.AsBase62()
}

func (id CcId160) GoString() string {
	if id.fingerprintSize > 0 {
		return fmt.Sprintf("CcId{size: %d, timestamp: %d (%s), fingerprint: 0x%x, payload: 0x%x}",
			id.Size(), id.Timestamp(), id.Time().Format(time.RFC3339), id.Fingerprint(), id.Payload())
	}
	return fmt.Sprintf("CcId{size: %d, timestamp: %d (%s), payload: 0x%x}",
		id.Size(), id.Timestamp(), id.Time().Format(time.RFC3339), id.Payload())
}

func (id CcId160) Bytes() []byte {
	return id.data[:]
}

func (id CcId160) AsBase62() string {
	v, _ := EncodeToBase62(id.data[:])
	return v
}

func (id CcId160) AsBase32() string {
	v, _ := EncodeToBase32(id.data[:])
	return v
}

func (id CcId160) AsBase16() string {
	v, _ := EncodeToBase16(id.data[:])
	return v
}

func NewCcId160WithFingerprint(timestamp uint32, fingerprint []byte, payload []byte) (CcId, error) {
	var id CcId160
	payloadSize := byte(len(payload))
	fingerprintSize := byte(len(fingerprint))
	if fingerprintSize > MaxFingerprintSize {
		return NilCcId160, InvalidFingerprintSizeError{
			ProvidedSize: fingerprintSize,
			RequiredSize: MaxFingerprintSize,
		}
	}
	minPayloadSize := ByteSliceSize160 - TimestampSize - fingerprintSize
	if payloadSize < minPayloadSize {
		return NilCcId160, InvalidPayloadSizeError{
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
