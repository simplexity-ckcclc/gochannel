package common

import (
	"testing"
	"time"
)

func TestTimeDurationToMillis(t *testing.T) {
	times := map[time.Duration]int64{
		3 * time.Second:               3000,
		5*time.Minute + 3*time.Second: 5*60*1000 + 3000,
	}

	for d, i := range times {
		o := TimeDurationToMillis(d)
		if o != i {
			t.Errorf("TimeDurationToMillis(%s) = %d. want %d", d, o, i)
		}
	}
}
