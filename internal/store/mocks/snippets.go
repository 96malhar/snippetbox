package mocks

import (
	"github.com/96malhar/snippetbox/internal/store"
	"time"
)

type MockSnippetStore struct {
	snippets []*store.Snippet
}

func (m *MockSnippetStore) Insert(title string, content string, expirationDays int) (int, error) {
	currentTime := time.Date(2000, time.January, 1, 12, 0, 0, 0, time.UTC)
	snippet := store.Snippet{
		ID:      m.generateId(),
		Title:   title,
		Content: content,
		Expires: currentTime.Add(time.Hour * 24 * time.Duration(expirationDays)),
	}
	m.snippets = append(m.snippets, &snippet)
	return snippet.ID, nil
}

func (m *MockSnippetStore) Get(id int) (*store.Snippet, error) {
	for _, sn := range m.snippets {
		if sn.ID == id {
			return sn, nil
		}
	}
	return nil, store.ErrNoRecord
}

func (m *MockSnippetStore) Latest() ([]*store.Snippet, error) {
	return m.snippets, nil
}

func (m *MockSnippetStore) generateId() int {
	return len(m.snippets) + 1
}

func NewMockSnippetStore(seed ...*store.Snippet) *MockSnippetStore {
	return &MockSnippetStore{
		snippets: seed,
	}
}
