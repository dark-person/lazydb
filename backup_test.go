package lazydb

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Create an outdated LazyDB file.
func createOutdatedDb(path string) (*LazyDB, error) {
	// Prepare database of version 2
	original := New(
		DbPath(path),
		Migrate(fsNormalTestV2, dirNormalTestV2))

	err := original.Connect()
	if err != nil {
		return nil, err
	}
	defer original.Close()

	bk, err := original.Migrate()
	if err != nil {
		return nil, err
	}

	if bk != "" {
		return nil, fmt.Errorf("Unexpected backup when create.")
	}

	return original, nil
}

// Ensure backup perform when backup directory is set.
func TestCreateAutoBackup(t *testing.T) {
	tmpDir := t.TempDir()

	path := filepath.Join(tmpDir, "backupTest.db")

	// ================== Old database =================
	original, err := createOutdatedDb(path)
	if err != nil {
		t.Errorf("Failed when connect original db")
		return
	}
	original.Close()

	// ================== Update Database =====================
	updated := New(
		DbPath(path),
		Migrate(fsNormalTestV3, dirNormalTestV3),
		BackupDir(filepath.Join(tmpDir, "bk")),
		Version(0),
	)

	err = updated.Connect()
	if err != nil {
		t.Errorf("Failed when connect updated db")
		return
	}
	defer updated.Close()

	bk, err := updated.Migrate()
	assert.Nilf(t, err, "Error when migrate")
	assert.NotEqualValuesf(t, "", bk, "Unexpected backup location")

	// Check filepath location
	flag := IsFileExist(bk)
	assert.EqualValuesf(t, true, flag, "File location is not contains path: %s", bk)
}

// Ensure no backup when backup directory not set.
func TestCreateAutoBackupNoDir(t *testing.T) {
	tmpDir := t.TempDir()

	path := filepath.Join(tmpDir, "data.db")

	// ================== Old database =================
	original, err := createOutdatedDb(path)
	if err != nil {
		t.Errorf("Failed when connect original db")
		return
	}
	original.Close()

	// ================== Update Database =====================
	updated := New(
		DbPath(path),
		Migrate(fsNormalTestV3, dirNormalTestV3),
		Version(0),
	)

	err = updated.Connect()
	if err != nil {
		t.Errorf("Failed when connect updated db")
		return
	}
	defer updated.Close()

	bk, err := updated.Migrate()
	assert.Nilf(t, err, "Error when migrate")
	assert.EqualValuesf(t, "", bk, "Unexpected backup location")
}

// Test no backup will created when new database is created.
func TestBackupWithNewDb(t *testing.T) {
	var err error
	tmpDir := t.TempDir()

	path := filepath.Join(tmpDir, "data.db")

	// ================== Create Database =====================
	newOne := New(
		DbPath(path),
		Migrate(fsNormalTestV3, dirNormalTestV3),
		BackupDir(filepath.Join(tmpDir, "bk")),
		Version(0),
	)

	err = newOne.Connect()
	if err != nil {
		t.Errorf("Failed when connect new created db")
		return
	}
	defer newOne.Close()

	// Perform migration (MUST be 1st time)
	bk, err := newOne.Migrate()
	assert.Nilf(t, err, "Error when migrate new db")
	assert.EqualValuesf(t, "", bk, "No backup should appear in new created db.")
}

// Ensure no backup when backup directory not set.
func TestCreateAutoBackupVersionSpecified(t *testing.T) {
	tmpDir := t.TempDir()

	path := filepath.Join(tmpDir, "data.db")

	// ================== Old database =================
	original, err := createOutdatedDb(path)
	if err != nil {
		t.Errorf("Failed when connect original db")
		return
	}
	original.Close()

	// ================== Update Database =====================
	updated := New(
		DbPath(path),
		Migrate(fsNormalTestV3, dirNormalTestV3),
		BackupDir(filepath.Join(tmpDir, "bk")),
		Version(2),
	)

	err = updated.Connect()
	if err != nil {
		t.Errorf("Failed when connect updated db")
		return
	}
	defer updated.Close()

	bk, err := updated.Migrate()
	assert.Nilf(t, err, "Error when migrate")
	assert.EqualValuesf(t, "", bk, "Unexpected backup location")
}
