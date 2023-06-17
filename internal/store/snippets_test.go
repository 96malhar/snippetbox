package store

import (
	dtmock "github.com/96malhar/snippetbox/internal/datetime/mocks"
	"testing"
	"time"
)

func TestSnippetStore_Get(t *testing.T) {
	tests := []struct {
		name         string
		id           int
		mockCurrTime time.Time
		wantSnippet  *Snippet
		wantErr      error
	}{
		{
			name:         "Exists",
			id:           1,
			mockCurrTime: parseTime(t, time.RFC3339, "2022-12-01T10:00:00Z"),
			wantSnippet: &Snippet{
				ID:      1,
				Title:   "An old silent pond",
				Content: "An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō",
				Created: parseTime(t, time.RFC3339, "2022-01-01T10:00:00Z"),
				Expires: parseTime(t, time.RFC3339, "2023-01-01T10:00:00Z"),
			},
			wantErr: nil,
		},
		{
			name:         "Exists but expired",
			id:           1,
			mockCurrTime: parseTime(t, time.RFC3339, "2024-01-01T10:00:00Z"),
			wantErr:      ErrNoRecord,
		},
		{
			name:    "Does not exist",
			id:      3,
			wantErr: ErrNoRecord,
		},
		{
			name:    "Negative ID",
			id:      -2,
			wantErr: ErrNoRecord,
		},
	}

	db, testDbName := newTestDB(t)
	setupDB(t, db)
	t.Cleanup(func() {
		db.Close()
		dropDB(t, testDbName)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSnippetStore(db)
			if !tt.mockCurrTime.IsZero() {
				s.datetimeHandler = &dtmock.MockDateTimeHandler{
					MockCurrentTime: tt.mockCurrTime,
				}
			}

			gotSnippet, err := s.Get(tt.id)
			if err != tt.wantErr {
				t.Fatalf("Get() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantSnippet == nil && gotSnippet != nil {
				t.Errorf("Expected nil snippet, got = %v", gotSnippet)
			}
			if tt.wantSnippet != nil {
				checkSnippet(t, gotSnippet, tt.wantSnippet)
			}
		})
	}
}

func checkSnippet(t *testing.T, gotSnippet, wantSnippet *Snippet) {
	if gotSnippet.ID != wantSnippet.ID {
		t.Errorf("got ID = %d; want ID = %d", gotSnippet.ID, wantSnippet.ID)
	}
	if gotSnippet.Title != wantSnippet.Title {
		t.Errorf("got Title = %s; want Title = %s", gotSnippet.Title, wantSnippet.Title)
	}
	if gotSnippet.Content != wantSnippet.Content {
		t.Errorf("got Title = %s; want Title = %s", gotSnippet.Content, wantSnippet.Content)
	}
	if gotSnippet.Created != wantSnippet.Created {
		t.Errorf("got Created = %v; want Created = %v", gotSnippet.Created, wantSnippet.Created)
	}
	if gotSnippet.Expires != wantSnippet.Expires {
		t.Errorf("got Expires = %v; want Expires = %v", gotSnippet.Expires, wantSnippet.Expires)
	}
}
