package datetime

import "time"

var now = time.Now

type Handler struct{}

func (h *Handler) GetCurrentTimeUTC() time.Time {
	return now().UTC().Round(time.Second)
}
