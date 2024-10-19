package lazydb

import (
	"database/sql"
	"embed"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

// Test function of Create LazyDB, check its final value.
// This test will not check embed.fs due to extra impl is needed.
func TestNew(t *testing.T) {
	// Case 1: Default
	db1 := New()

	assert.EqualValuesf(t, &LazyDB{
		dbPath:        "data.db",
		migrateDir:    "schema",
		migrateFs:     embed.FS{},
		schemaVersion: 0,
		backupDir:     "",
	}, db1, "Incorrect value in default.")

	// Case 2~3: Partly Modify
	db2 := New(DbPath("data.sqlite"), Version(2))
	assert.EqualValuesf(t, &LazyDB{
		dbPath:        "data.sqlite",
		migrateDir:    "schema",
		migrateFs:     embed.FS{},
		schemaVersion: 2,
		backupDir:     "",
	}, db2, "Incorrect value in modified 1.")

	db3 := New(Migrate(embed.FS{}, "kk"))
	assert.EqualValuesf(t, &LazyDB{
		dbPath:        "data.db",
		migrateDir:    "kk",
		migrateFs:     embed.FS{},
		schemaVersion: 0,
		backupDir:     "",
	}, db3, "Incorrect value in modified 2.")

	db4 := New(BackupDir("./abc"))
	assert.EqualValuesf(t, &LazyDB{
		dbPath:        "data.db",
		migrateDir:    "schema",
		migrateFs:     embed.FS{},
		schemaVersion: 0,
		backupDir:     "./abc",
	}, db4, "Incorrect value in modified 2.")

	// Case 4: All modify
	db5 := New(DbPath("data2.db"), Version(14), Migrate(embed.FS{}, "ee"), BackupDir("./abc"))
	assert.EqualValuesf(t, &LazyDB{
		dbPath:        "data2.db",
		migrateDir:    "ee",
		migrateFs:     embed.FS{},
		schemaVersion: 14,
		backupDir:     "./abc",
	}, db5, "Incorrect value in all modified.")
}

// Test Method of Connect().
// This test will NOT consider $HOME directory as a case,
// all tests are using custom path only.
func TestConnect(t *testing.T) {
	// Test Case type
	type testCase struct {
		a    *LazyDB
		want error
	}

	// Existing database creation
	existPath := filepath.Join(t.TempDir(), "exist_test.db")
	f, _ := os.Create(existPath)
	f.Close()

	// Prepare Tests
	tests := []testCase{
		// Database that not exist
		{&LazyDB{dbPath: filepath.Join(t.TempDir(), "connect_test.db")}, nil},
		// Existed database
		{&LazyDB{dbPath: existPath}, nil},
		// Empty Path
		{&LazyDB{dbPath: ""}, ErrEmptyPath},
	}

	// Run Tests
	for idx, tt := range tests {
		err := tt.a.Connect()

		assert.ErrorIsf(t, err, tt.want, "Case %d: Unexpected error result: %v", idx, err)

		// Close connection (by sql library but not LazyDB)
		if tt.a.db != nil {
			tt.a.db.Close()
		}
	}
}

// Test Close Database connection, which will:
//   - Enable to run even database is nil
//   - DB must be nil after close
//
// This test is rely on createFile() function to create a new database.
func TestClose(t *testing.T) {
	// Test Case Type
	type testCase struct {
		path    *string // Path of database, nil if wanted db is nil
		wantErr bool    // Determine non-nil error will be returned
	}

	// Prepare Test Case
	gracePath := filepath.Join(t.TempDir(), "db_close.db")
	tests := []testCase{
		// Graceful Case
		{&gracePath, false},
		// Empty database
		{nil, false},
	}

	// Run Tests
	for idx, tt := range tests {
		var l *LazyDB

		if tt.path == nil {
			// Create LazyDB with nil db value
			l = &LazyDB{db: nil}
		} else {
			// Create database file with createFile()
			l = &LazyDB{dbPath: *tt.path}
			createDbFile(*tt.path)

			// Connect database with connect(), with db is non-nil value ONLY
			l.db, _ = sql.Open(DatabaseType, l.dbPath)
		}

		// Close database connection with LazyDB.Close()
		err := l.Close()

		// Check want errors
		assert.EqualValues(t, tt.wantErr, err != nil,
			"Case %d: Unexpected error as %v", idx, err)

		// Ensure Nil Value of sql.DB
		assert.Nilf(t, l.db, "Case %d: Unexpected non-nil database connection.", idx)
	}
}
