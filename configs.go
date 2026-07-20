package lazydb

import "embed"

// DatabaseType imply supported database driver type. Only support sqlite3.
const DatabaseType = "sqlite3"

// Latest imply the latest schema version that supported.
// Use zero to get latest schema version as possible.
var Latest uint = 0

// Dir represent the default folder path that store schema migration script in embed.FS.
var Dir = "schema"

// Path represent default File path of database.
var Path = "data.db"

// Get default options when creating database.
func defaultOpts() databaseOpts {
	return databaseOpts{
		DBPath:        Path,
		MigrateFS:     embed.FS{},
		MigrateDir:    Dir,
		SchemaVersion: Latest,
	}
}
