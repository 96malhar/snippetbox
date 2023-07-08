package main

import "github.com/96malhar/snippetbox/internal/store"

type snippetStoreInterface interface {
	Insert(title string, content string, expirationDays int) (int, error)
	Get(id int) (*store.Snippet, error)
	Latest() ([]*store.Snippet, error)
}

type userStoreInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	Get(id int) (*store.User, error)
}
