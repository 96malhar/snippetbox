package store

import (
	"testing"
)

func TestSnippetStore_Get(t *testing.T) {
	tests := []struct {
		name        string
		id          int
		wantSnippet *Snippet
		wantErr     error
	}{
		{
			name: "Exists",
			id:   1,
			wantSnippet: &Snippet{
				ID:      1,
				Title:   "An old silent pond",
				Content: "An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō",
			},
			wantErr: nil,
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

	s := NewSnippetStore(db)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSnippet, err := s.Get(tt.id)
			if err != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr = %v", err, tt.wantErr)
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
		t.Errorf("got Title = %s; want Title = %d", gotSnippet.Title, wantSnippet.ID)
	}
	if gotSnippet.Content != wantSnippet.Content {
		t.Errorf("got Title = %s; want Title = %d", gotSnippet.Title, wantSnippet.ID)
	}
}
