package datetime

import (
	"testing"
	"time"
)

func TestDateTimeHandler_GetCurrentTimeUTC(t *testing.T) {
	tests := []struct {
		name    string
		setTime time.Time
		want    time.Time
	}{
		{
			name:    "Round up to seconds",
			setTime: time.Date(2000, 12, 25, 12, 0, 2, 500000001, time.UTC),
			want:    time.Date(2000, 12, 25, 12, 0, 3, 0, time.UTC),
		},
		{
			name:    "Round down to seconds",
			setTime: time.Date(2000, 12, 25, 12, 0, 2, 499999999, time.UTC),
			want:    time.Date(2000, 12, 25, 12, 0, 2, 0, time.UTC),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now = func() time.Time {
				return tt.setTime
			}
			h := &Handler{}
			if got := h.GetCurrentTimeUTC(); !got.Equal(tt.want) {
				t.Errorf("GetCurrentTimeUTC() = %v, want %v", got, tt.want)
			}
		})
	}
}
