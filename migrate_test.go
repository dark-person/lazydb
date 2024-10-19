package lazydb

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const dirNormalTestV2 = "test_schema/normal_v2"
const dirNormalTestV3 = "test_schema/normal_v3"

//go:embed test_schema/normal_v2/*
var fsNormalTestV2 embed.FS

//go:embed test_schema/normal_v3/*
var fsNormalTestV3 embed.FS

// Function to prepare usable AppDB object.
// Purpose of function is to reuse code.
//
// Developer MUST close AppDB after testing is completed,
// to prevent filesystem blocking database removal.
func prepareTestAppDB(t *testing.T, dbName string) (*LazyDB, error) {
	// Remove any existing database
	os.Remove(dbName)

	// Create AppDB object
	l := New(
		DbPath(dbName),
		Version(3),
		Migrate(fsNormalTestV3, "test_schema/normal_v3"),
	)

	// Connect to database
	err := l.Connect()
	if err != nil {
		t.Fatal("Failed to connect, return error: ", err)
	}
	return l, err
}

// Function to prepare usable AppDB object.
// Purpose of function is to reuse code.
//
// Developer MUST close AppDB after testing is completed,
// to prevent filesystem blocking database removal.
func prepareTestLatestAppDB(t *testing.T, dbName string) (*LazyDB, error) {
	// Use existing func
	a, err := prepareTestAppDB(t, dbName)
	if err != nil {
		panic(err)
	}

	// Step to latest
	_, err = a.Migrate()
	if err != nil {
		t.Fatal("Failed to prepare db as latest: ", err)
	}

	// Return as success
	return a, err
}

// Get user_version value from database.
//
// This utility function is ONLY for internal use.
func getUserVersion(db *sql.DB) (int, error) {
	if db == nil {
		return -1, ErrNilDatabase
	}

	// Get User Version
	row := db.QueryRow("PRAGMA user_version")
	if row == nil {
		return -1, fmt.Errorf("nil row query")
	}

	// Scan value into var
	var userVersion int
	err := row.Scan(&userVersion)
	if err != nil {
		return -1, err
	}

	// Return
	return userVersion, nil
}

// Test migration to latest. Act as smoke test only.
func TestToLatest(t *testing.T) {
	// Get temp directory
	tempDir := t.TempDir()

	// Create AppDB object
	a, err := prepareTestAppDB(t, filepath.Join(tempDir, "test.db"))
	if err != nil {
		t.Fatal("Failed to prepare db: ", err)
	}

	// Migration
	_, err = a.Migrate()
	if err != nil {
		t.Error("MigrateUp return error: ", err)
	}

	ver, err := getUserVersion(a.db)
	assert.Nil(t, err, "Migration should not get error when get user version.")
	assert.EqualValuesf(t, 3, ver, "Version incorrect.")

	// Migrate again
	_, err = a.Migrate()
	if err != nil {
		t.Error("Migrate Again return error: ", err)
	}

	ver, err = getUserVersion(a.db)
	assert.Nil(t, err, "Migration should not get error when get user version.")
	assert.EqualValuesf(t, 3, ver, "Version incorrect.")

	a.Close()
}

// Test migration to specified version.
func TestToSpecific(t *testing.T) {
	// Get temp directory
	tempDir := t.TempDir()

	// Test case struct, use fsTesting for all tests
	type testCase struct {
		dbName      string // Name of database file
		schemaVer   uint   // Schema Version number
		wantErr     bool   // Should error appear
		wantVersion int    // user_version that should be appear in migrated database
	}

	tests := []testCase{
		// Normal Input
		{"test0.db", 1, false, 1},
		{"test1.db", 2, false, 2},
		{"test2.db", 3, false, 3},

		// To Latest
		{"test3.db", 0, false, 3},

		// Invalid input
		{"test4.db", 4, true, -1},
	}

	// Run tests
	for idx, tt := range tests {
		// Create AppDB object
		a, err := prepareTestLatestAppDB(t, filepath.Join(tempDir, tt.dbName))

		if err != nil {
			t.Fatal("Failed to prepare db: ", err)
		}
		defer a.Close() // Close to prevent leak

		// Main test scope
		_, err = a.MigrateTo(tt.schemaVer)
		if tt.wantErr {
			assert.Errorf(t, err, "Case %d should return error.", idx)
		} else {
			assert.Nilf(t, err, "Case %d should not return error, but got %v.", idx, err)
		}

		// Get user version
		if !tt.wantErr {
			ver, _ := getUserVersion(a.db)
			assert.EqualValuesf(t, tt.wantVersion, ver, "user_version not matched in case %d: %d", idx, ver)
		}
	}
}

func TestInvalidMigrate(t *testing.T) {
	// Get temp directory
	tempDir := t.TempDir()

	// Prepare db path
	noEmbedFs := filepath.Join(tempDir, "no_embed.db")
	noSchema := filepath.Join(tempDir, "empty_schema.db")
	noDirFound := filepath.Join(tempDir, "dir_not_found.db")
	notConnect := filepath.Join(tempDir, "no_connect.db")

	// Clear all exist database file
	os.Remove(noEmbedFs)
	os.Remove(noSchema)
	os.Remove(noDirFound)
	os.Remove(notConnect)

	// Case: No embed fs
	l1 := New(DbPath(noEmbedFs))
	l1.Connect()
	defer l1.Close()

	_, err := l1.Migrate()
	assert.NotNilf(t, err, "No embed fs should return error.")

	// Case: empty schema
	l2 := New(
		DbPath(noSchema),
		Migrate(fsNormalTestV3, ""),
	)
	l2.Connect()
	defer l2.Close()

	_, err = l2.Migrate()
	assert.ErrorIs(t, err, ErrEmptyDir, "Empty directory should return ErrEmptyDir")

	// Case: directory not found in embed fs
	l3 := New(
		DbPath(noDirFound),
		Version(1),
		Migrate(fsNormalTestV3, "no_no_no"),
	)
	l3.Connect()
	defer l3.Close()

	_, err = l3.Migrate()
	assert.NotNilf(t, err, "Directory not found in embed.fs should return error.")

	// Case: database is not connected
	l4 := New(
		DbPath(notConnect),
		Version(1),
		Migrate(fsNormalTestV3, "test_schema"),
	)
	defer l4.Close()

	_, err = l4.Migrate()
	assert.ErrorIs(t, err, ErrNilDatabase, "Not connected db should return ErrNilDatabase")
}
