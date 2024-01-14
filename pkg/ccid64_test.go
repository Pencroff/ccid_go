package pkg

import (
	"encoding/binary"
	"fmt"
	"testing"
	"time"
)

func TestNewCcId64(t *testing.T) {
	keys := SortKeys(TestCaseCcId64Map)
	for _, key := range keys {
		tc := TestCaseCcId64Map[key]
		t.Run(key, func(t *testing.T) {
			var got CcId
			got, _ = NewCcId64WithFingerprint(tc.timestamp, tc.Fingerprint, tc.payload)
			size := byte(len(tc.Bytes))
			if got.Size() != size {
				t.Errorf("Size:\nNewCcId64(%x, %x, %x) =\n%d, want\n%d",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Size(), size)
			}
			if got.Timestamp() != tc.timestamp {
				t.Errorf("Timestamp:\nNewCcId64(%x, %x, %x) =\n%d, want\n%d",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Timestamp(), tc.timestamp)
			}
			if got.Time() != tc.time {
				t.Errorf("Time:\nNewCcId64(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Time().Format(time.RFC3339), tc.time.Format(time.RFC3339))
			}
			byteIdx := 4 + len(tc.Fingerprint)
			if !SliceEqual(got.Payload(), tc.Bytes[byteIdx:]) {
				t.Errorf("Payload:\nNewCcId64(%x, %x, %x) =\n%x, want\n%x",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Payload(), tc.Bytes[byteIdx:])
			}
			if !SliceEqual(got.Fingerprint(), tc.Fingerprint) {
				t.Errorf("Fingerprint:\nNewCcId64(%x, %x, %x) =\n%v, want\n%v",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Fingerprint(), tc.Fingerprint)
			}
			if !SliceEqual(got.Bytes(), tc.Bytes) {
				t.Errorf("Bytes:\nNewCcId64(%x, %x, %x) =\n%x, want\n%x",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Bytes(), tc.Bytes)
			}
		})
		t.Run(key+"_as_base", func(t *testing.T) {
			var got CcId
			got, _ = NewCcId64WithFingerprint(tc.timestamp, tc.Fingerprint, tc.payload)
			if got.String() != tc.Base62 {
				t.Errorf("String-AsBase62:\nNewCcId64(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, got.String(), tc.Base62)
			}
			if got.AsBase62() != tc.Base62 {
				t.Errorf("AsBase62:\nNewCcId64(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, got.AsBase62(), tc.Base62)
			}
			if got.AsBase32() != tc.Base32 {
				t.Errorf("AsBase32:\nNewCcId64(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, got.AsBase32(), tc.Base32)
			}
			if got.AsBase16() != tc.Base16 {
				t.Errorf("AsBase16:\nNewCcId64(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, got.AsBase16(), tc.Base16)
			}
		})
		t.Run(key+"_go_string", func(t *testing.T) {
			var got CcId
			got, _ = NewCcId64WithFingerprint(tc.timestamp, tc.Fingerprint, tc.payload)
			v := fmt.Sprintf("%#v", got)
			if v != tc.GoString {
				t.Errorf("GoString:\nNewCcId64(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, v, tc.GoString)
			}
		})
		t.Run(key+"_uint64", func(t *testing.T) {
			var got CcId
			got, _ = NewCcId64WithFingerprint(tc.timestamp, tc.Fingerprint, tc.payload)
			v, ok := got.(CcId64)
			if !ok {
				t.Errorf("Cast to CcId64 failed")
			}
			if v.Uint64() != binary.BigEndian.Uint64(tc.Bytes) {
				t.Errorf("Uint64:\nNewCcId64(%x, %x, %x) =\n%d, want\n%d",
					tc.timestamp, tc.Fingerprint, tc.payload, v.Uint64(), binary.BigEndian.Uint64(tc.Bytes))
			}
		})
	}
}

func TestNewCcId64_Error(t *testing.T) {
	keys := SortKeys(testCaseCcId64ErrorMap)
	for _, key := range keys {
		tc := testCaseCcId64ErrorMap[key]
		t.Run(key, func(t *testing.T) {
			var got CcId
			var err error
			got, err = NewCcId64WithFingerprint(tc.timestamp, tc.fingerprint, tc.payload)
			if err == nil || err.Error() != tc.errMsg {
				t.Errorf("NewCcId64(%d, %x) error = %v,\nwantErr %v", tc.timestamp, tc.payload, err, tc.errMsg)
				return
			}
			if !SliceEqual(got.Bytes(), tc.bytes) {
				t.Errorf("NewCcId64(%d, %x) =\n%x, want\n%x", tc.timestamp, tc.payload, got.Bytes(), tc.bytes)
			}
		})
	}
}
