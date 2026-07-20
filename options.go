package lazydb

import "io/fs"

// Final options for create database. Internal usage only.
type databaseOpts struct {
	DBPath         string // Absolute path of .db file
	ConnectOptions string // DSN Options in connection string, e.g. _journal_mode, _mutex
	MigrateFS      fs.FS  // FS to be used for migration
	MigrateDir     string // directory that contains migration sql files
	SchemaVersion  uint   // Schema version that using
	BackupDir      string // directory that used to backup database file
}

// Option of database.
type DatabaseOption interface {
	apply(*databaseOpts)
}

// ---------------------------------------------------
type dbPath string

func (p dbPath) apply(opts *databaseOpts) {
	opts.DBPath = string(p)
}

// DbPath will provide option for lazyDB constructor to use given database path.
// Please note if any database options is involved in path string will lead to panic during Connect().
//
// e.g. DSN of "file:test.db?cache=shared&mode=memory", should use below format when passing option to LazyDB.New():
//     DbPath("test.db"), DSNOption("cache=shared&mode=memory")
func DbPath(path string) DatabaseOption {
	return dbPath(path)
}

// ---------------------------------------------------
type connectOptions string

func (p connectOptions) apply(opts *databaseOpts) {
	opts.ConnectOptions = string(p)
}

// ConnectOptions imply Data Source Name string options,
// which are append after the filename of the SQLite database.
// The database filename and options are separated by an ? (Question Mark).
//
// For example, to use a database of with WAL journal mode,
// 		lazydb.ConnectOptions("_journal_mode=WAL")
//
// Details should refer to https://github.com/mattn/go-sqlite3#connection-string
func ConnectOptions(option string) DatabaseOption {
	return connectOptions(option)
}

// ---------------------------------------------------
type migrateParam struct {
	MigrateFS  fs.FS
	MigrateDir string
}

func (m migrateParam) apply(opts *databaseOpts) {
	opts.MigrateFS = m.MigrateFS
	opts.MigrateDir = m.MigrateDir
}

// Use given file system (e.g. embed.FS) to perform migration.
func Migrate(f fs.FS, dir string) DatabaseOption {
	return migrateParam{f, dir}
}

// ---------------------------------------------------

type schemaVer uint

func (s schemaVer) apply(opts *databaseOpts) {
	opts.SchemaVersion = uint(s)
}

// Specific version of schema to be used.
//
// If value is <= 0, then it will use as latest version as possible
// that defined in migration schema fs.
func Version(ver int) DatabaseOption {
	return schemaVer(ver)
}

// ---------------------------------------------------
type backupDir string

func (b backupDir) apply(opts *databaseOpts) {
	opts.BackupDir = string(b)
}

// Backup to given directory before migration.
// The backup filename is fixed to {original_name}_bk_{time}.{ext}.
func BackupDir(path string) DatabaseOption {
	return backupDir(path)
}
