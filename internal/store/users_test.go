package store

import (
	"github.com/96malhar/snippetbox/internal/testutils"
	"github.com/stretchr/testify/assert"
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
			name: "Exists",
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
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestUserStore_Get(t *testing.T) {
	testutils.RunAsIntegTest(t)
	testcases := []struct {
		name    string
		id      int
		checks  []userCheck
		wantErr error
	}{
		{
			name:   "Exists",
			id:     1,
			checks: []userCheck{userHasId(1), userHasName("John"), userHasEmail("john@example.com")},
		},
		{
			name:    "Does not exist",
			id:      2,
			wantErr: ErrNoRecord,
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
			gotUser, err := s.Get(tc.id)
			assert.ErrorIs(t, err, tc.wantErr)

			for _, check := range tc.checks {
				check(t, gotUser)
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
			assert.ErrorIs(t, err, tc.wantErr)
			assert.Equal(t, tc.wantId, gotId)
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

			assert.ErrorIs(t, err, tc.wantErr)

			if tc.wantId != 0 {
				exists, err := s.Exists(tc.wantId)
				assert.True(t, exists)
				assert.NoError(t, err)
			}
		})
	}
}
