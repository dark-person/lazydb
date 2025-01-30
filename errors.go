package lazydb

import "errors"

// Error when user pass empty string as database path parameter.
var ErrEmptyPath = errors.New("empty database file path")

// Error when user try to create a new database that not ".db" extension.
var ErrInvalidExt = errors.New("invalid file extension of database")

// Error when database is nil value, no operation can perform.
var ErrNilDatabase = errors.New("database is nil")

// Error when try to execute multiple statement with nil/empty slice.
var ErrEmptyStmt = errors.New("no statement to execute")

// Error when migration directory is empty string.
var ErrEmptyDir = errors.New("empty string for migration directory")

// Error when migration directory structure is not correct.
var ErrInvalidDir = errors.New("invalid migration directory structure")
