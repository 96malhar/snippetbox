package store

import (
	"github.com/96malhar/snippetbox/internal/datetime/mocks"
	"github.com/96malhar/snippetbox/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSnippetStore_Get(t *testing.T) {
	testutils.RunAsIntegTest(t)
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
				s.datetimeHandler = mocks.NewMockDateTimeHandler(parseTime(t, time.RFC3339, tt.mockCurrentTime))
			}

			gotSnippet, err := s.Get(tt.id)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantSnippet, gotSnippet)
		})
	}
}

func TestSnippetStore_Latest(t *testing.T) {
	testutils.RunAsIntegTest(t)
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
				s.datetimeHandler = mocks.NewMockDateTimeHandler(parseTime(t, time.RFC3339, tt.mockCurrentTime))
			}

			gotSnippets, err := s.Latest()

			require.NoError(t, err)
			assert.Equal(t, len(tt.wantSnippets), len(gotSnippets))
			for i := range tt.wantSnippets {
				assert.Equal(t, tt.wantSnippets[i], gotSnippets[i])
			}
		})
	}
}

func TestSnippetStore_Insert(t *testing.T) {
	testutils.RunAsIntegTest(t)
	db, testDbName := newTestDB(t)
	setupDB(t, db)
	t.Cleanup(func() {
		db.Close()
		dropDB(t, testDbName)
	})

	s := NewSnippetStore(db)
	mockCurrTime := time.Date(2022, 1, 1, 10, 0, 0, 0, time.UTC)
	s.datetimeHandler = mocks.NewMockDateTimeHandler(mockCurrTime)

	id, err := s.Insert("Snippet 3 Title", "Snippet 3 content.", 10)

	require.NoError(t, err)
	assert.Equal(t, 3, id)

	wantSnippet := &Snippet{
		ID:      3,
		Title:   "Snippet 3 Title",
		Content: "Snippet 3 content.",
		Created: mockCurrTime,
		Expires: mockCurrTime.Add(time.Hour * 24 * 10),
	}
	gotSnippet, _ := s.Get(3)
	assert.Equal(t, wantSnippet, gotSnippet)
}
