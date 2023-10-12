package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	// Establish a sql.DB connection pool for our test database
	// Because our setup and teardown scripts contain multiple SQL statements,
	// we need to use "multiStatements=true" parameter in our DSN. This instructs
	// MySQL database driver to support executing multiple SQL statements in one 
	// db.Exec() call.
	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script form file and execute the statements
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// User the t.Cleanup() to register a function which will automatically be called by
	// Go when the current test (or sub-test) which calls the newTestDB() has finished
	// In this section we read and execute teardown script
	t.Cleanup(func() {
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}

		db.Close()
	})

	// Return the database connection pool
	return db
}