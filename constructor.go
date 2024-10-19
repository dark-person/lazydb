package lazydb

import (
	"database/sql"
	"embed"

	_ "github.com/mattn/go-sqlite3"
)

// The database implementation for lazy people. Using sqlite3.
//
// Support multiple lazy feature:
//   - Lazy creation for *sql.DB with sqlite3
//   - Wrapper function to simplify sql.stmt
//   - Migration support if migration fs & directory is properly set
type LazyDB struct {
	db *sql.DB // Database connection

	dbPath        string   // Database absolute path, for easy reuse
	migrateFs     embed.FS // FS for schema migrations sql scripts
	migrateDir    string   // Directory for storing migration script, default is "schema"
	schemaVersion uint     // version of migration script to use
	backupDir     string   // directory to backup, or empty string for no backup. Default is empty string.
}

// Create a new LazyDB.
func New(opts ...DatabaseOption) *LazyDB {
	// Get default option
	opt := defaultOpts()

	// Apply options
	for _, item := range opts {
		item.apply(&opt)
	}

	// Return object
	return &LazyDB{
		dbPath:        opt.DbPath,
		migrateDir:    opt.MigrateDir,
		migrateFs:     opt.MigrateFS,
		schemaVersion: opt.SchemaVersion,
		backupDir:     opt.BackupDir,
	}
}

// Connect to database, when path already stored in AppDB.
func (l *LazyDB) Connect() error {
	// Prevent Empty Path
	if l.dbPath == "" {
		return ErrEmptyPath
	}

	var err error

	// Create Database if need
	err = createDbFile(l.dbPath)
	if err != nil {
		return err
	}

	// Open database connection, which create file if not exist
	l.db, err = sql.Open(DatabaseType, l.dbPath)
	if err != nil {
		return err
	}

	// Test DB connection by ping
	return l.db.Ping()
}

// Close all existing database connection.
//
// If AppDB has no database connected, then this function has no effect,
// with no error returned.
func (l *LazyDB) Close() error {
	// Prevent no connection for nil pointer
	if l.db == nil {
		return nil
	}

	// Close connection
	err := l.db.Close()
	l.db = nil

	// Return error
	return err
}

// Get *sql.DB created.
func (l *LazyDB) DB() *sql.DB {
	return l.db
}
