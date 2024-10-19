package lazydb

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createDummyTable(db *sql.DB) {
	var sql string

	sql = "CREATE TABLE IF NOT EXISTS test_table (content text NOT NULL, val INT);"
	db.Exec(sql)

	sql = "INSERT INTO test_table (content, val) VALUES (\"abc\", 123), (\"def\", 456);"
	db.Exec(sql)
}

func TestWrapperNilDb(t *testing.T) {
	db := &LazyDB{db: nil}
	var err error

	// Test Query
	_, err = db.Query("abc")
	assert.ErrorIs(t, err, ErrNilDatabase)

	// Test QueryRow
	_, err = db.QueryRow("abc")
	assert.ErrorIs(t, err, ErrNilDatabase)

	// Test Exec
	_, err = db.Exec("abc")
	assert.ErrorIs(t, err, ErrNilDatabase)

	// Test ExecMultiple
	_, err = db.ExecMultiple([]ParamQuery{{"abc", nil}})
	assert.ErrorIs(t, err, ErrNilDatabase, "Expected ErrNilDatabase error return")
}

// This test only check situation when db is nil value & param slice is nil.
func TestExecMultipleNilSlice(t *testing.T) {
	// Test for nil slice
	db := &LazyDB{db: &sql.DB{}}

	_, err := db.ExecMultiple(nil)
	assert.ErrorIs(t, err, ErrEmptyStmt, "Expected ErrEmptyStmt error return")
}

func TestQueryRow(t *testing.T) {
	// Prepare temp directory
	temp := t.TempDir()

	// Prepare db path & ensure it is not exist
	dbPath := filepath.Join(temp, "query_row.db")
	os.Remove(dbPath)

	// Prepare db
	db := New(
		DbPath(dbPath),
	)
	db.Connect()
	defer db.Close()

	// Run Create query
	createDummyTable(db.DB())

	// Query Row
	row, err := db.QueryRow("SELECT val FROM test_table WHERE content=\"abc\"")
	assert.Nilf(t, err, "<QueryRow> No error should appear in query row.")

	var value int
	err = row.Scan(&value)
	assert.Nilf(t, err, "<QueryRow> No error should appear in scan row.")

	assert.EqualValuesf(t, 123, value, "unexpected value: %d", value)
}

func TestQuery(t *testing.T) {
	// Prepare temp directory
	temp := t.TempDir()

	// Prepare db path & ensure it is not exist
	dbPath := filepath.Join(temp, "query.db")
	os.Remove(dbPath)

	// Prepare db
	db := New(
		DbPath(dbPath),
	)
	db.Connect()
	defer db.Close()

	// Run Create query
	createDummyTable(db.DB())

	// Query Row
	rows, err := db.Query("SELECT content FROM test_table WHERE 1=1")
	assert.Nilf(t, err, "<Query> No error should appear in query row.")

	var results []string

	for rows.Next() {
		var content string
		err = rows.Scan(&content)
		assert.Nilf(t, err, "<Query> No error should appear in scan rows.")

		results = append(results, content)
	}

	assert.EqualValuesf(t, []string{"abc", "def"}, results, "<Query> unexpected value: %v", results)
}

func TestExec(t *testing.T) {
	// Prepare temp directory
	temp := t.TempDir()

	// Prepare db path & ensure it is not exist
	dbPath := filepath.Join(temp, "exec.db")
	os.Remove(dbPath)

	// Prepare db
	db := New(
		DbPath(dbPath),
	)
	db.Connect()
	defer db.Close()

	// Run Create query
	createDummyTable(db.DB())

	// Declare common value
	var err error
	var query string
	var row *sql.Row
	var ct int // row counter

	// ------------ Exec Row (INSERT) ---------------
	_, err = db.Exec("INSERT INTO test_table(content, val) VALUES (\"ghi\", 789)")
	assert.Nilf(t, err, "<Exec Insert> No error should appear in query row.")

	// Query if any record inserted
	query = "SELECT COUNT(*) FROM test_table WHERE content = \"ghi\" AND val = 789"
	row = db.db.QueryRow(query)

	err = row.Scan(&ct)
	assert.Nilf(t, err, "<Exec Insert> No error should appear in scan row.")

	assert.EqualValuesf(t, 1, ct, "<Exec Insert> unexpected value: %v", ct)

	//  ------------ Exec Row (UPDATE) ---------------
	_, err = db.Exec("UPDATE test_table SET val = 456 WHERE content = \"ghi\"")
	assert.Nilf(t, err, "<Exec Update> No error should appear in query row.")

	// Query if any record inserted
	query = "SELECT COUNT(*) FROM test_table WHERE val = 456"
	row = db.db.QueryRow(query)

	err = row.Scan(&ct)
	assert.Nilf(t, err, "<Exec Update> No error should appear in scan row.")

	assert.EqualValuesf(t, 2, ct, "<Exec Update> unexpected value: %v", ct)

	//  ------------ Exec Row (DELETE) ---------------
	_, err = db.Exec("DELETE FROM test_table WHERE content = \"ghi\"")
	assert.Nilf(t, err, "<Exec Delete> No error should appear in query row.")

	// Query if any record inserted
	query = "SELECT COUNT(*) FROM test_table WHERE val = 456"
	row = db.db.QueryRow(query)

	err = row.Scan(&ct)
	assert.Nilf(t, err, "<Exec Delete> No error should appear in scan row.")

	assert.EqualValuesf(t, 1, ct, "<Exec Delete> unexpected value: %v", ct)
}

func TestPureExecMultiple(t *testing.T) {
	// Prepare temp directory
	temp := t.TempDir()

	// Prepare db path & ensure it is not exist
	dbPath := filepath.Join(temp, "exec_multiple_pure.db")
	os.Remove(dbPath)

	// Prepare db
	db := New(
		DbPath(dbPath),
	)
	db.Connect()
	defer db.Close()

	// Run Create query
	createDummyTable(db.DB())

	// Declare common value
	var err error
	var query string
	var row *sql.Row
	var ct int // row counter

	// Prepare multiple INSERT statements
	var queries []ParamQuery
	queries = append(queries, ParamQuery{"INSERT INTO test_table(content, val) VALUES (\"ghi\", 789)", nil})
	queries = append(queries, ParamQuery{"INSERT INTO test_table(content, val) VALUES (?, ?)", []any{"jkl", 444}})
	queries = append(queries, ParamQuery{"INSERT INTO test_table(content, val) VALUES (\"opq\", 777)", nil})

	_, err = db.ExecMultiple(queries)
	assert.Nilf(t, err, "<ExecMultiple Insert> No error should appear in insert row.")

	// Query if any record inserted
	query = "SELECT COUNT(*) FROM test_table"
	row = db.db.QueryRow(query)

	err = row.Scan(&ct)
	assert.Nilf(t, err, "<ExecMultiple Insert> No error should appear in scan row.")

	assert.EqualValuesf(t, 5, ct, "<ExecMultiple Insert> unexpected value: %v", ct)
}

func TestMixedExecMultiple(t *testing.T) {
	// Prepare temp directory
	temp := t.TempDir()

	// Prepare db path & ensure it is not exist
	dbPath := filepath.Join(temp, "exec_multiple_pure.db")
	os.Remove(dbPath)

	// Prepare db
	db := New(
		DbPath(dbPath),
	)
	db.Connect()
	defer db.Close()

	// Run Create query
	createDummyTable(db.DB())

	// Declare common value
	var err error
	var query string
	var row *sql.Row
	var ct int // row counter

	// Prepare multiple INSERT statements
	var queries []ParamQuery
	queries = append(queries, Param("INSERT INTO test_table(content, val) VALUES (\"k89\", 789)"))
	queries = append(queries, Param("DELETE FROM test_table WHERE content=?", "k89"))

	_, err = db.ExecMultiple(queries)
	assert.Nilf(t, err, "<ExecMultiple Insert> No error should appear in insert row.")

	// Query if any record inserted
	query = "SELECT COUNT(*) FROM test_table"
	row = db.db.QueryRow(query)

	err = row.Scan(&ct)
	assert.Nilf(t, err, "<ExecMultiple Insert> No error should appear in scan row.")

	assert.EqualValuesf(t, 2, ct, "<ExecMultiple Insert> unexpected value: %v", ct)
}
