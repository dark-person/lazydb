package lazydb

import "embed"

// Database driver type. Only support sqlite3.
const DatabaseType = "sqlite3"

// Latest schema version that supported.
// Use zero to get latest schema version as possible.
var Latest uint = 0

// Default folder path that store schema migration script in embed.FS.
var Dir = "schema"

// Default File path of database.
var Path = "data.db"

// Get default options when creating database.
func defaultOpts() databaseOpts {
	return databaseOpts{
		DbPath:        Path,
		MigrateFS:     embed.FS{},
		MigrateDir:    Dir,
		SchemaVersion: Latest,
	}
}
