package models

import (
	"database/sql"
	"os"
	"testing"
)

// Create new database instance, put dummy data in it
// automatically delete all tables and close connection when calling test func is finished
func newTestDB(t *testing.T) *sql.DB {
	// multiStatements allows our teardown script to perform two statements in one db.Exec() call
	db, err := sql.Open("mysql", "test_web:pass@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(script))
	if err != nil {
		t.Fatal(err)
	}

	// this executes given func when test calling newTestDB is finished
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
	return db
}
