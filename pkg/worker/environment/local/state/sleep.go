package state

import (
	"time"

	"github.com/senseyeio/duration"
)

func calcTimeDuration(durationIso string) (time.Duration, error) {
	if durationIso == "" {
		return 0, nil
	}

	sleep, err := duration.ParseISO8601(durationIso)
	if err != nil {
		return 0, err
	}

	d1 := time.Now()
	d2 := sleep.Shift(d1)
	return d2.Sub(d1), nil
}
