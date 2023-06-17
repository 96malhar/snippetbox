package mocks

import "time"

type MockDateTimeHandler struct {
	MockCurrentTime time.Time
}

func (h *MockDateTimeHandler) GetCurrentTimeUTC() time.Time {
	return h.MockCurrentTime
}
