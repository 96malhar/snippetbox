package mocks

import (
	"github.com/96malhar/snippetbox/internal/store"
	"time"
)

type MockUserStore struct {
	users []*store.User
}

func (m *MockUserStore) Insert(name, email, password string) error {
	if email == "dupe@example.com" {
		return store.ErrDuplicateEmail
	}
	user := store.User{
		ID:             m.generateId(),
		Name:           name,
		Email:          email,
		HashedPassword: []byte(password),
		Created:        time.Date(1996, time.April, 28, 3, 0, 0, 0, time.UTC),
	}
	m.users = append(m.users, &user)
	return nil
}

func (m *MockUserStore) Authenticate(email, password string) (int, error) {
	for _, usr := range m.users {
		if usr.Email == email && string(usr.HashedPassword) == password {
			return usr.ID, nil
		}
	}
	return 0, store.ErrInvalidCredentials
}

func (m *MockUserStore) Exists(id int) (bool, error) {
	for _, usr := range m.users {
		if usr.ID == id {
			return true, nil
		}
	}
	return false, nil
}

func (m *MockUserStore) Get(id int) (*store.User, error) {
	for _, usr := range m.users {
		if usr.ID == id {
			return usr, nil
		}
	}
	return nil, store.ErrNoRecord
}

func (m *MockUserStore) generateId() int {
	return len(m.users) + 1
}

func NewMockUserStore(users ...*store.User) *MockUserStore {
	return &MockUserStore{
		users: users,
	}
}
