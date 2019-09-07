package common

import "time"

func TimeToMillis(t time.Time) int64 {
	return t.UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

func TimeDurationToMillis(d time.Duration) int64 {
	return int64(d) / (int64(time.Millisecond))
}
