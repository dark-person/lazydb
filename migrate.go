package lazydb

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	sqlite "github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Create a new `migrate.Migrate` instance, which can be used to migrate up/down.
//
// Purpose of this function is to reuse code to create a new migrate object.
func (l *LazyDB) migrateInstance() (*migrate.Migrate, error) {
	// Prevent db is nil
	if l.db == nil {
		return nil, ErrNilDatabase
	}

	// Prevent empty migration directory
	if l.migrateDir == "" {
		return nil, ErrEmptyDir
	}

	// Get sqlite3 instance
	instance, err := sqlite.WithInstance(l.db, &sqlite.Config{})
	if err != nil {
		return nil, err
	}

	// Create iofs with embedded filesystem, with directory specified
	f, err := iofs.New(l.migrateFs, l.migrateDir)
	if err != nil {
		return nil, err
	}

	return migrate.NewWithInstance("iofs", f, "sqlite3", instance)
}

// Migrate database to latest supported version, which is defined when create new LazyDB.
//
// If backup directory is set, and migration is actually performed,
// then this function will also return backup database path.
// Otherwise empty string will be returned.
func (l *LazyDB) Migrate() (backupPath string, err error) {
	return l.MigrateTo(l.schemaVersion)
}

// Migrate database to specified version.
//
// When version is 0, then it will migrate to latest version as possible,
// otherwise it will migrate to specified version.
//
// If backup directory is set, and migration is actually performed,
// then this function will also return backup database path.
// Otherwise empty string will be returned.
func (l *LazyDB) MigrateTo(version uint) (backupPath string, err error) {
	// Prepare migration instance
	l.mig, err = l.migrateInstance()
	if err != nil {
		return "", err
	}

	// Run backup
	backupPath, err = l.backup(l.mig)
	if err != nil {
		return backupPath, err
	}

	// Perform migration depend on version is equals to 0
	if version == 0 {
		err = l.mig.Up()
	} else {
		err = l.mig.Migrate(version)
	}

	// Early return if no error
	if err == nil {
		return backupPath, nil
	}

	// No changes applied, which is acceptable
	if errors.Is(err, migrate.ErrNoChange) {
		return backupPath, nil
	}

	return backupPath, err
}
