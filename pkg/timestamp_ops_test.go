package pkg

import (
	"testing"
	"time"
)

type testCaseTimestampOps struct {
	time time.Time
	ts   uint32
}

var testCaseTimestampOpsMap = map[string]testCaseTimestampOps{
	"min time": {
		time: time.Date(2014, 05, 13, 16, 53, 20, 0, time.UTC),
		ts:   0,
	},
	"2014-05-13T16:53:21Z": {
		time: time.Date(2014, 05, 13, 16, 53, 21, 0, time.UTC),
		ts:   1,
	},
	"2020-01-01T00:00:00Z": {
		time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		ts:   177836800,
	},
	"0x12345678": {
		time: time.Date(2024, 1, 16, 15, 44, 56, 0, time.UTC),
		ts:   0x12345678, // 305419896
	},
	"2024-01-01T23:59:59Z": {
		time: time.Date(2024, 1, 1, 23, 59, 59, 0, time.UTC),
		ts:   304153599,
	},
	"max time": {
		time: time.Date(2150, 06, 19, 23, 21, 35, 0, time.UTC),
		ts:   0xFFFFFFFF, // 4294967295
	},
}

func TestToAdjustedTimestamp(t *testing.T) {
	keys := SortKeys(testCaseTimestampOpsMap)
	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			tc := testCaseTimestampOpsMap[key]
			got := ToAdjustedTimestamp(tc.time)
			if got != tc.ts {
				t.Errorf("ToAdjustedTimestamp(%v) =\n%d, want\n%d", tc.time, got, tc.ts)
			}
		})
	}
}

func TestToStandardizedTime(t *testing.T) {
	keys := SortKeys(testCaseTimestampOpsMap)
	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			tc := testCaseTimestampOpsMap[key]
			got := ToStandardizedTime(tc.ts)
			if got.UTC() != tc.time {
				t.Errorf("ToStandardizedTime(%d) =\n%s, want\n%s", tc.ts, got.Format(time.RFC3339), tc.time.Format(time.RFC3339))
			}
		})
	}
}
