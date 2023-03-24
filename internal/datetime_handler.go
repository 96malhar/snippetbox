package internal

import "time"

type DateTimeHandler struct{}

func (d *DateTimeHandler) GetCurrentTimeUTC() time.Time {
	return time.Now().UTC().Round(time.Second)
}
