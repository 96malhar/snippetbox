package store

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

const superUserConnStr = "host=localhost port=5432 user=postgres sslmode=disable dbname=postgres"

func newTestDB(t *testing.T) (*sql.DB, string) {
	randomSuffix := strings.Split(uuid.New().String(), "-")[0]
	testDBName := fmt.Sprintf("snippetbox_test_%s", randomSuffix)

	db, err := sql.Open("postgres", superUserConnStr)
	if err != nil {
		t.Fatalf("Failed to connect to database postgres. Err = %s", err)
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil {
		t.Fatalf("Failed to create database %s. Err = %s", testDBName, err)
	}

	db, err = sql.Open("postgres", fmt.Sprintf("host=localhost port=5432 user=postgres sslmode=disable dbname=%s", testDBName))
	if err != nil || db.Ping() != nil {
		t.Fatalf("Failed to connect to database %s", testDBName)
	}

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
	db, err := sql.Open("postgres", superUserConnStr)
	if err != nil {
		t.Fatalf("Failed to connect to database postgres. Err = %s", err)
	}
	_, err = db.Exec(fmt.Sprintf("DROP DATABASE %s", dbName))
	if err != nil {
		t.Fatalf("Failed to drop database %s. Err = %s", dbName, err)
	}
}
