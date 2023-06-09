package mocks

import "time"

type MockDateTimeHandler struct {
	MockCurrentTime time.Time
}

func NewMockDateTimeHandler(mockCurrentTime time.Time) *MockDateTimeHandler {
	return &MockDateTimeHandler{
		MockCurrentTime: mockCurrentTime,
	}
}

func (h *MockDateTimeHandler) GetCurrentTimeUTC() time.Time {
	return h.MockCurrentTime
}
