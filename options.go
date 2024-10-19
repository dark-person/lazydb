package lazydb

import (
	"embed"
)

// Final options for create database. Internal usage only.
type databaseOpts struct {
	DbPath        string   // Absolute path of .db file
	MigrateFS     embed.FS // FS to be used for migration
	MigrateDir    string   // directory that contains migration sql files
	SchemaVersion uint     // Schema version that using
	BackupDir     string   // directory that used to backup database file
}

// Option of database.
type DatabaseOption interface {
	apply(*databaseOpts)
}

// ---------------------------------------------------
type dbPath string

func (p dbPath) apply(opts *databaseOpts) {
	opts.DbPath = string(p)
}

// Use given database path.
func DbPath(path string) DatabaseOption {
	return dbPath(path)
}

// ---------------------------------------------------
type migrateParam struct {
	MigrateFS  embed.FS
	MigrateDir string
}

func (m migrateParam) apply(opts *databaseOpts) {
	opts.MigrateFS = m.MigrateFS
	opts.MigrateDir = m.MigrateDir
}

// Use given embed file system to perform migration.
func Migrate(f embed.FS, dir string) DatabaseOption {
	return migrateParam{f, dir}
}

// ---------------------------------------------------

type schemaVer uint

func (s schemaVer) apply(opts *databaseOpts) {
	opts.SchemaVersion = uint(s)
}

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