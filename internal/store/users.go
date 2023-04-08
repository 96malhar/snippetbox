package store

import (
	"database/sql"
	"github.com/96malhar/snippetbox/internal/datetime"
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserStore struct {
	db              *sql.DB
	datetimeHandler interface {
		GetCurrentTimeUTC() time.Time
	}
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{db: db, datetimeHandler: &datetime.DateTimeHandler{}}
}

func (m *UserStore) Insert(name, email, password string) error {
	return nil
}

func (m *UserStore) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserStore) Exists(id int) (bool, error) {
	return false, nil
}
