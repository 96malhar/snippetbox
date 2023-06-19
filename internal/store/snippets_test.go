package store

import (
	dtmock "github.com/96malhar/snippetbox/internal/datetime/mocks"
	"testing"
	"time"
)

func TestSnippetStore_Get(t *testing.T) {
	tests := []struct {
		name            string
		id              int
		mockCurrentTime string
		wantSnippet     *Snippet
		wantErr         error
	}{
		{
			name:            "Exists ID=1",
			id:              1,
			mockCurrentTime: "2022-12-01T10:00:00Z",
			wantSnippet: &Snippet{
				ID:      1,
				Title:   "Snippet 1 Title",
				Content: "Snippet 1 content.",
				Created: parseTime(t, time.RFC3339, "2022-01-01T10:00:00Z"),
				Expires: parseTime(t, time.RFC3339, "2023-01-01T10:00:00Z"),
			},
			wantErr: nil,
		},
		{
			name:            "Exists ID=2",
			id:              2,
			mockCurrentTime: "2022-12-01T10:00:00Z",
			wantSnippet: &Snippet{
				ID:      2,
				Title:   "Snippet 2 Title",
				Content: "Snippet 2 content.",
				Created: parseTime(t, time.RFC3339, "2022-02-01T10:00:00Z"),
				Expires: parseTime(t, time.RFC3339, "2023-02-01T10:00:00Z"),
			},
			wantErr: nil,
		},
		{
			name:            "Exists but expired",
			id:              1,
			mockCurrentTime: "2024-01-01T10:00:00Z",
			wantErr:         ErrNoRecord,
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
			if tt.mockCurrentTime != "" {
				s.datetimeHandler = dtmock.NewMockDateTimeHandler(parseTime(t, time.RFC3339, tt.mockCurrentTime))
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

func TestSnippetStore_Latest(t *testing.T) {
	testcases := []struct {
		name            string
		mockCurrentTime string
		wantSnippets    []*Snippet
	}{
		{
			name:            "Unexpired snippets",
			mockCurrentTime: "2022-12-01T10:00:00Z",
			wantSnippets: []*Snippet{
				{
					ID:      2,
					Title:   "Snippet 2 Title",
					Content: "Snippet 2 content.",
					Created: parseTime(t, time.RFC3339, "2022-02-01T10:00:00Z"),
					Expires: parseTime(t, time.RFC3339, "2023-02-01T10:00:00Z"),
				},
				{
					ID:      1,
					Title:   "Snippet 1 Title",
					Content: "Snippet 1 content.",
					Created: parseTime(t, time.RFC3339, "2022-01-01T10:00:00Z"),
					Expires: parseTime(t, time.RFC3339, "2023-01-01T10:00:00Z"),
				},
			},
		},
		{
			name:            "Expired snippets",
			mockCurrentTime: "2024-12-01T10:00:00Z",
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			db, testDbName := newTestDB(t)
			setupDB(t, db)
			t.Cleanup(func() {
				db.Close()
				dropDB(t, testDbName)
			})

			s := NewSnippetStore(db)
			if tt.mockCurrentTime != "" {
				s.datetimeHandler = dtmock.NewMockDateTimeHandler(parseTime(t, time.RFC3339, tt.mockCurrentTime))
			}

			gotSnippets, err := s.Latest()
			if err != nil {
				t.Fatalf("Unexpected error = %v", err)
			}

			if len(gotSnippets) != len(tt.wantSnippets) {
				t.Fatalf("len(gotSnippets)=%d; len(wantSnippets)=%d", len(gotSnippets), len(tt.wantSnippets))
			}

			for i := range tt.wantSnippets {
				checkSnippet(t, gotSnippets[i], tt.wantSnippets[i])
			}
		})
	}
}

func TestSnippetStore_Insert(t *testing.T) {
	db, testDbName := newTestDB(t)
	setupDB(t, db)
	t.Cleanup(func() {
		db.Close()
		dropDB(t, testDbName)
	})

	s := NewSnippetStore(db)
	mockCurrTime := time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC)
	s.datetimeHandler = dtmock.NewMockDateTimeHandler(mockCurrTime)

	id, err := s.Insert("Snippet 3 Title", "Snippet 3 content.", 10)

	if err != nil {
		t.Fatalf("Unexpected err = %v", err)
	}

	if id != 3 {
		t.Errorf("got ID = %d; want ID = 3", id)
	}

	wantSnippet := &Snippet{
		ID:      3,
		Title:   "Snippet 3 Title",
		Content: "Snippet 3 content.",
		Created: mockCurrTime,
		Expires: mockCurrTime.Add(time.Hour * 24 * 10),
	}
	gotSnippet, _ := s.Get(3)
	checkSnippet(t, gotSnippet, wantSnippet)
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
