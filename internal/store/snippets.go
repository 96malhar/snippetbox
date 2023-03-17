package store

import (
	"database/sql"
	"errors"
	"github.com/96malhar/snippetbox/internal"
	"time"
)

// Snippet defines the type to hold the data for an individual snippet. Notice how
// the fields of the struct correspond to the fields in our MySQL snippets
// table?
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetStore is a type which wraps a sql.DB connection pool.
type SnippetStore struct {
	db              *sql.DB
	datetimeHandler interface {
		GetCurrentTimeUTC() time.Time
	}
}

func NewSnippetStore(db *sql.DB) *SnippetStore {
	return &SnippetStore{db: db, datetimeHandler: &internal.DateTimeHandler{}}
}

// Insert will add a new snippet into the database and return the snippet ID.
func (s *SnippetStore) Insert(title string, content string, expirationDays time.Duration) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
    VALUES($1, $2, $3, $4)
	returning id`

	created := s.datetimeHandler.GetCurrentTimeUTC()
	expires := created.Add(time.Hour * 24 * expirationDays)

	var id int
	err := s.db.QueryRow(stmt, title, content, created, expires).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// Get will return a specific snippet based on its id.
func (s *SnippetStore) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > $1 AND id = $2`

	var sn Snippet
	currTime := s.datetimeHandler.GetCurrentTimeUTC()
	err := s.db.QueryRow(stmt, currTime, id).Scan(&sn.ID, &sn.Title, &sn.Content, &sn.Created, &sn.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return &sn, nil
}

// Latest will return the 10 most recently created snippets.
func (s *SnippetStore) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
    		WHERE expires > $1 ORDER BY id DESC LIMIT 10`

	rows, err := s.db.Query(stmt, s.datetimeHandler.GetCurrentTimeUTC())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []*Snippet
	for rows.Next() {
		var sn Snippet
		err = rows.Scan(&sn.ID, &sn.Title, &sn.Content, &sn.Created, &sn.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, &sn)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
