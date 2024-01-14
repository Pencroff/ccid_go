package ccid_go

import (
	"encoding/binary"
	p "github.com/Pencroff/ccid_go/pkg"
	"io"
	"time"
)

// CcIdGen is an interface for generating CcIds, or CcId  generator.
type CcIdGen interface {
	// Next generates the next CcId using the current time.
	Next() (p.CcId, error)
	// NextWithTime generates the next CcId using the provided time.
	// 't' is time.Time for generated CcId
	NextWithTime(t time.Time) (p.CcId, error)
}

// CcIdGenImplementation is an implementation of CcIdGen.
type CcIdGenImplementation struct {
	size          byte
	nilCcId       p.CcId
	fingerprint   []byte
	payload       []byte
	ctor          p.CcIdCtor
	rndRd         io.Reader
	strategy      p.CcIdMonotonicStrategy
	clock         p.Clock
	lastTimestamp uint32
	lastPayload   []byte
}

func (g *CcIdGenImplementation) Next() (p.CcId, error) {
	return g.NextWithTime(g.clock.Now())
}

func (g *CcIdGenImplementation) NextWithTime(t time.Time) (p.CcId, error) {
	timestamp := p.ToAdjustedTimestamp(t)
	var carry byte
	if g.strategy != nil && timestamp <= g.lastTimestamp {
		timestamp = g.lastTimestamp
		g.payload, carry = g.strategy.Mutate(g.lastPayload)
		if carry > 0 {
			timestamp += 1
			_, err := g.rndRd.Read(g.payload)
			if err != nil {
				return g.nilCcId, err
			}
		}
	} else {
		_, err := g.rndRd.Read(g.payload)
		if err != nil {
			return g.nilCcId, err
		}
	}
	g.lastTimestamp = timestamp
	copy(g.lastPayload, g.payload[:])
	return g.ctor(timestamp, g.fingerprint, g.payload)
}

// NewCcIdGen creates a new CcId Generator (no fingerprint, no monotonic strategy).
// 'size' must be the size of the CcId in bytes.
// 'rndRd' must be a reader for providing random bytes.
func NewCcIdGen(size byte, rndRd io.Reader) (CcIdGen, error) {
	return newCcIdGenWithClock(size, nil, rndRd, nil, p.RealClock{})
}

// NewCcIdGenWithFingerprint creates a new CcId Generator with fingerprint (no monotonic strategy).
// 'size' must be the size of the CcId in bytes.
// 'fingerprint' must be a byte slice of the correct size for the CcId.
// 'rndRd' must be a reader for providing random bytes.
func NewCcIdGenWithFingerprint(size byte, fingerprint []byte, rndRd io.Reader) (CcIdGen, error) {
	return newCcIdGenWithClock(size, fingerprint, rndRd, nil, p.RealClock{})
}

// NewMonotonicCcIdGen creates a new CcId Generator with monotonic strategy (no fingerprint).
// 'size' must be the size of the CcId in bytes.
// 'rndRd' must be a reader for providing random bytes.
// 'strategy' must be a monotonic strategy, to set logic for payload mutation.
func NewMonotonicCcIdGen(size byte, rndRd io.Reader, strategy p.CcIdMonotonicStrategy) (CcIdGen, error) {
	return newCcIdGenWithClock(size, nil, rndRd, strategy, p.RealClock{})
}

// NewMonotonicCcIdGenWithFingerprint creates a new CcId Generator with fingerprint and monotonic strategy.
// 'size' must be the size of the CcId in bytes.
// 'fingerprint' must be a byte slice of the correct size for the CcId.
// 'rndRd' must be a reader for providing random bytes.
// 'strategy' must be a monotonic strategy, to set logic for payload mutation.
func NewMonotonicCcIdGenWithFingerprint(size byte, fingerprint []byte, rndRd io.Reader, strategy p.CcIdMonotonicStrategy) (CcIdGen, error) {
	return newCcIdGenWithClock(size, fingerprint, rndRd, strategy, p.RealClock{})
}

func newCcIdGenWithClock(size byte, fingerprint []byte, rndRd io.Reader, s p.CcIdMonotonicStrategy, c p.Clock) (CcIdGen, error) {
	var ctor p.CcIdCtor
	var nilCcId p.CcId
	switch size {
	case p.ByteSliceSize64:
		ctor = p.NewCcId64WithFingerprint
		nilCcId = p.NilCcId64
	case p.ByteSliceSize96:
		ctor = p.NewCcId96WithFingerprint
		nilCcId = p.NilCcId96
	case p.ByteSliceSize128:
		ctor = p.NewCcId128WithFingerprint
		nilCcId = p.NilCcId128
	case p.ByteSliceSize160:
		ctor = p.NewCcId160WithFingerprint
		nilCcId = p.NilCcId160
	}
	payloadSize := size - p.TimestampSize - byte(len(fingerprint))
	return &CcIdGenImplementation{
		size:          size,
		nilCcId:       nilCcId,
		fingerprint:   fingerprint,
		payload:       make([]byte, payloadSize),
		ctor:          ctor,
		rndRd:         rndRd,
		strategy:      s,
		clock:         c,
		lastTimestamp: 0,
		lastPayload:   make([]byte, payloadSize),
	}, nil
}

// FromString creates a CcId from a string. It requires the fingerprint size and the base of the string.
// 's' must be a string of the correct size for the CcId.
// 'fingerprintSize' must be the size of the fingerprint in bytes.
// 'base' must be the base of the string. It can be 16, 32 or 62.
func FromString(s string, fingerprintSize byte, base byte) (p.CcId, error) {
	var b []byte
	var err error
	switch base {
	case p.BASE62:
		b, err = p.DecodeFromBase62(s)
		if err != nil {
			return nil, err
		}
	case p.BASE32:
		b, err = p.DecodeFromBase32(s)
		if err != nil {
			return nil, err
		}
	case p.BASE16:
		b, err = p.DecodeFromBase16(s)
		if err != nil {
			return nil, err
		}
	}
	return FromBytes(b, fingerprintSize)
}

// FromBytes creates a CcId from a byte slice. It requires the fingerprint size.
// 'b' must be a byte slice of the correct size for the CcId.
// The size can be determined by calling CcId.Size().
// 'fingerprintSize' must be the size of the fingerprint in bytes.
// It's 0 or 1 for CcId64, 0 to 5 for CcId96, CcId128, CcId160.
func FromBytes(b []byte, fingerprintSize byte) (p.CcId, error) {
	l := len(b)
	var c p.CcIdCtor
	switch l {
	case p.ByteSliceSize64:
		c = p.NewCcId64WithFingerprint
	case p.ByteSliceSize96:
		c = p.NewCcId96WithFingerprint
	case p.ByteSliceSize128:
		c = p.NewCcId128WithFingerprint
	case p.ByteSliceSize160:
		c = p.NewCcId160WithFingerprint
	default:
		return nil, p.InvalidLengthError(byte(l))
	}
	timestamp := binary.BigEndian.Uint32(b[:p.TimestampSize])
	fingerprintEndIdx := p.TimestampSize + fingerprintSize
	return c(timestamp, b[p.TimestampSize:fingerprintEndIdx], b[fingerprintEndIdx:])
}
