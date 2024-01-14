package pkg

import (
	"fmt"
	"testing"
	"time"
)

func TestNewCcId160(t *testing.T) {
	keys := SortKeys(TestCaseCcId160Map)
	for _, key := range keys {
		tc := TestCaseCcId160Map[key]
		t.Run(key, func(t *testing.T) {
			var got CcId
			got, _ = NewCcId160WithFingerprint(tc.timestamp, tc.Fingerprint, tc.payload)
			size := byte(len(tc.Bytes))
			if got.Size() != size {
				t.Errorf("Size:\nNewCcId160(%x, %x, %x) =\n%d, want\n%d",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Size(), size)
			}
			if got.Timestamp() != tc.timestamp {
				t.Errorf("Timestamp:\nNewCcId160(%x, %x, %x) =\n%d, want\n%d",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Timestamp(), tc.timestamp)
			}
			if got.Time() != tc.time {
				t.Errorf("Time:\nNewCcId160(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Time().Format(time.RFC3339), tc.time.Format(time.RFC3339))
			}
			byteIdx := 4 + len(tc.Fingerprint)
			if !SliceEqual(got.Payload(), tc.Bytes[byteIdx:]) {
				t.Errorf("Payload:\nNewCcId160(%x, %x, %x) =\n%x, want\n%x",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Payload(), tc.Bytes[byteIdx:])
			}
			if !SliceEqual(got.Fingerprint(), tc.Fingerprint) {
				t.Errorf("Fingerprint:\nNewCcId160(%x, %x, %x) =\n%v, want\n%v",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Fingerprint(), tc.Fingerprint)
			}
			if !SliceEqual(got.Bytes(), tc.Bytes) {
				t.Errorf("Bytes:\nNewCcId160(%x, %x, %x) =\n%x, want\n%x",
					tc.timestamp, tc.Fingerprint, tc.payload, got.Bytes(), tc.Bytes)
			}
		})
		t.Run(key+"_as_base", func(t *testing.T) {
			var got CcId
			got, _ = NewCcId160WithFingerprint(tc.timestamp, tc.Fingerprint, tc.payload)
			if got.String() != tc.Base62 {
				t.Errorf("String-AsBase62:\nNewCcId160(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, got.String(), tc.Base62)
			}
			if got.AsBase62() != tc.Base62 {
				t.Errorf("AsBase62:\nNewCcId160(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, got.AsBase62(), tc.Base62)
			}
			if got.AsBase32() != tc.Base32 {
				t.Errorf("AsBase32:\nNewCcId160(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, got.AsBase32(), tc.Base32)
			}
			if got.AsBase16() != tc.Base16 {
				t.Errorf("AsBase16:\nNewCcId160(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, got.AsBase16(), tc.Base16)
			}
		})
		t.Run(key+"_go_string", func(t *testing.T) {
			var got CcId
			got, _ = NewCcId160WithFingerprint(tc.timestamp, tc.Fingerprint, tc.payload)
			v := fmt.Sprintf("%#v", got)
			if v != tc.GoString {
				t.Errorf("GoString:\nNewCcId160(%x, %x, %x) =\n%s, want\n%s",
					tc.timestamp, tc.Fingerprint, tc.payload, v, tc.GoString)
			}
		})
	}
}

func TestNewCcId160_Error(t *testing.T) {
	keys := SortKeys(testCaseCcId160ErrorMap)
	for _, key := range keys {
		tc := testCaseCcId160ErrorMap[key]
		t.Run(key, func(t *testing.T) {
			var got CcId
			var err error
			got, err = NewCcId160WithFingerprint(tc.timestamp, tc.fingerprint, tc.payload)
			if err == nil || err.Error() != tc.errMsg {
				t.Errorf("NewCcId160(%x, %x, %x) error = %v,\nwantErr %v", tc.timestamp, tc.fingerprint, tc.payload, err, tc.errMsg)
				return
			}
			if !SliceEqual(got.Bytes(), tc.bytes) {
				t.Errorf("NewCcId160(%x, %x, %x) =\n%x, want\n%x", tc.timestamp, tc.fingerprint, tc.payload, got.Bytes(), tc.bytes)
			}
		})
	}
}
