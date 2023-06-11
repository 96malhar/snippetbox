package datetime

import "time"

var utcNow = func() time.Time {
	return time.Now().UTC()
}

type Handler struct{}

func (h *Handler) GetCurrentTimeUTC() time.Time {
	return utcNow().Round(time.Second)
}
