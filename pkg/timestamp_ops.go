package pkg

import "time"

// ToAdjustedTimestamp converts a time.Time to a uint32 timestamp with custom epoch
func ToAdjustedTimestamp(t time.Time) uint32 {
	return uint32(t.Unix() - epochStamp)
}

// ToStandardizedTime converts a uint32 timestamp with custom epoch to a UTC time.Time
func ToStandardizedTime(ts uint32) time.Time {
	return time.Unix(int64(ts)+epochStamp, 0).UTC()
}
