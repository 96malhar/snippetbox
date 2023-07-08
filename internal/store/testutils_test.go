package store

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func newTestDB(t *testing.T) (*sql.DB, string) {
	randomSuffix := strings.Split(uuid.New().String(), "-")[0]
	testDBName := fmt.Sprintf("snippetbox_test_%s", randomSuffix)

	db := getDBConn(t, "postgres", "postgres", "postgres")
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		t.Fatalf("Failed to create database %s. Err = %s", testDBName, err)
	}
	db.Close()

	db = getDBConn(t, "postgres", "postgres", testDBName)
	t.Logf("Connected to database %s", testDBName)
	return db, testDBName
}

func setupDB(t *testing.T, db *sql.DB) {
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}
}

func cleanupDB(t *testing.T, db *sql.DB) {
	script, err := os.ReadFile("./testdata/teardown.sql")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}
}

func dropDB(t *testing.T, dbName string) {
	db := getDBConn(t, "postgres", "postgres", "postgres")
	defer db.Close()
	_, err := db.Exec(fmt.Sprintf("DROP DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("Failed to drop database %s. Err = %s", dbName, err)
	}
}

func getDBConn(t *testing.T, user, password, dbname string) *sql.DB {
	dsn := fmt.Sprintf("host=localhost port=5432 user=%s password=%s sslmode=disable dbname=%s", user, password, dbname)
	db, _ := sql.Open("postgres", dsn)
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to connect to postgres with DSN = %s\nError = %s", dsn, err)
	}
	return db
}

func parseTime(t *testing.T, layout, value string) time.Time {
	t.Helper()
	tm, err := time.Parse(layout, value)
	if err != nil {
		t.Fatalf("An error occurred while parsing the time. Err = %s", err)
	}
	return tm
}

type userCheck func(t *testing.T, user *User)

func userHasId(id int) userCheck {
	return func(t *testing.T, user *User) {
		if user.ID != id {
			t.Errorf("got ID = %d; want ID = %d", user.ID, id)
		}
	}
}

func userHasName(name string) userCheck {
	return func(t *testing.T, user *User) {
		if user.Name != name {
			t.Errorf("got name = %s; want name = %s", user.Name, name)
		}
	}
}

func userHasEmail(email string) userCheck {
	return func(t *testing.T, user *User) {
		if user.Email != email {
			t.Errorf("got email = %s; want email = %s", user.Email, email)
		}
	}
}
