package store

import (
	"database/sql"
	"errors"
	"github.com/96malhar/snippetbox/internal/datetime"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"strings"
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

func (s *UserStore) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
    VALUES($1, $2, $3, $4)`

	createdAt := s.datetimeHandler.GetCurrentTimeUTC()

	_, err = s.db.Exec(stmt, name, email, string(hashedPassword), createdAt)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code == "23505" && strings.Contains(pqErr.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}

	return nil
}

func (s *UserStore) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = $1"

	err := s.db.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	return id, nil
}

func (s *UserStore) Exists(id int) (bool, error) {
	return false, nil
}
