package triple

import "time"

func nextCheck(old time.Time, interval time.Duration, offset time.Duration) time.Time {
	wakeTime := old.Truncate(interval)
	if interval == time.Hour*24 { // daily
		wakeTime = time.Date(wakeTime.Year(), wakeTime.Month(), wakeTime.Day(), 0, 0, 0, 0, wakeTime.Location())
		// gate以8点钟为日线开始
		if offset > 0 {
			wakeTime = wakeTime.Add(offset)
		}
	} else if interval == time.Hour*24*7 { // weekly
	}

	return wakeTime
}
