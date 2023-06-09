package datetime

import "time"

var utcNow = func() time.Time {
	return time.Now().UTC()
}

type DateTimeHandler struct{}

func (d *DateTimeHandler) GetCurrentTimeUTC() time.Time {
	return utcNow().Round(time.Second)
}
