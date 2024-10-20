package lazydb

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
)

// Create a backup of current database.
//
// If backup directory is empty string, then no backup will be created.
func (l *LazyDB) createBackup() (dest string, err error) {
	// Prevent backup directory is empty string
	if l.backupDir == "" {
		return "", nil // Consider as graceful return
	}

	// Prepare timestamp for backup
	str := time.Now().Format("20060102150405")

	// Get ext
	ext := filepath.Ext(l.dbPath)

	// Get original database name
	base := filepath.Base(l.dbPath)
	base = strings.Replace(base, ext, "", 1)

	// Perform backup
	dest = filepath.Join(l.backupDir, base+"_bk_"+str+ext)

	err = createDbFile(dest)
	if err != nil {
		return dest, err
	}

	err = copyFile(l.dbPath, dest)
	return dest, err
}

// Start backup process. If version is latest (i.e. no need to update), then no backup will be performed.
func (l *LazyDB) backup(m *migrate.Migrate) (dest string, err error) {
	if m == nil {
		return "", fmt.Errorf("migrate instance is nil")
	}

	// Get current database version
	current, _, err := m.Version()

	// Check err, Ignore nil version as it may be new db
	if err != nil && err != migrate.ErrNilVersion {
		return "", err
	}

	// Determine database is newly created or not
	isNew := (err == migrate.ErrNilVersion)

	// Get latest database schema version
	latest, err := LargestSchemaVer(l.migrateFs, l.migrateDir)
	if err != nil {
		return "", err
	}

	if latest > current && !isNew {
		return l.createBackup()
	}

	return "", nil // Consider as graceful return
}
