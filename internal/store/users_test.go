package store

import (
	"github.com/96malhar/snippetbox/internal/testutils"
	"testing"
)

func TestUserStore_Exists(t *testing.T) {
	testutils.RunAsIntegTest(t)
	testcases := []struct {
		name    string
		id      int
		want    bool
		wantErr error
	}{
		{
			name: "Does exist",
			id:   1,
			want: true,
		},
		{
			name: "Does not exist",
			id:   2,
			want: false,
		},
	}

	db, testDbName := newTestDB(t)
	setupDB(t, db)
	t.Cleanup(func() {
		db.Close()
		dropDB(t, testDbName)
	})

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewUserStore(db)
			got, err := s.Exists(tc.id)
			if got != tc.want {
				t.Errorf("Exists(%d) = %v; want = %v", tc.id, got, tc.want)
			}
			if err != tc.wantErr {
				t.Errorf("err = %v; wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestUserStore_Authenticate(t *testing.T) {
	testutils.RunAsIntegTest(t)
	testcases := []struct {
		name         string
		userEmail    string
		userPassword string
		wantId       int
		wantErr      error
	}{
		{
			name:         "valid user",
			userEmail:    "john@example.com",
			userPassword: "Hello, World!",
			wantId:       1,
		},
		{
			name:         "Invalid email",
			userEmail:    "jane@example.com",
			userPassword: "Hello, World!",
			wantId:       0,
			wantErr:      ErrInvalidCredentials,
		},
		{
			name:         "Invalid password",
			userEmail:    "john@example.com",
			userPassword: "Bye, World!",
			wantId:       0,
			wantErr:      ErrInvalidCredentials,
		},
	}

	db, testDbName := newTestDB(t)
	setupDB(t, db)
	t.Cleanup(func() {
		db.Close()
		dropDB(t, testDbName)
	})

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewUserStore(db)
			gotId, err := s.Authenticate(tc.userEmail, tc.userPassword)
			if gotId != tc.wantId {
				t.Errorf("gotId = %d; wantId = %d", gotId, tc.wantId)
			}
			if err != tc.wantErr {
				t.Errorf("err = %v; wantErr = %v", err, tc.wantErr)
			}
		})
	}
}

func TestUserStore_Insert(t *testing.T) {
	testutils.RunAsIntegTest(t)
	testcases := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		wantId       int
		wantErr      error
	}{
		{
			name:         "Valid new user",
			userName:     "Jane",
			userEmail:    "jane@example.com",
			userPassword: "random-pass-123",
			wantId:       2,
		},
		{
			name:         "duplicate email",
			userName:     "Jane",
			userEmail:    "john@example.com",
			userPassword: "random-pass-123",
			wantId:       0,
			wantErr:      ErrDuplicateEmail,
		},
	}

	db, testDbName := newTestDB(t)
	setupDB(t, db)
	t.Cleanup(func() {
		db.Close()
		dropDB(t, testDbName)
	})

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			s := NewUserStore(db)
			err := s.Insert(tc.userName, tc.userEmail, tc.userPassword)
			if err != tc.wantErr {
				t.Fatalf("Error = %v; Want Error = %v", err, tc.wantErr)
			}

			if tc.wantId != 0 {
				exists, err := s.Exists(tc.wantId)
				if err != nil {
					t.Errorf("Unexpected error = %v", err)
				}
				if exists != true {
					t.Errorf("A user with ID = %d does not exist in the DB", tc.wantId)
				}
			}
		})
	}
}
