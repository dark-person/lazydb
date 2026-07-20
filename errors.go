package lazydb

import "errors"

// ErrEmptyPath appear when user pass empty string as database path parameter.
var ErrEmptyPath = errors.New("empty database file path")

// ErrInvalidExt appear when user try to create a new database that not ".db" extension.
var ErrInvalidExt = errors.New("invalid file extension of database")

// ErrNilDatabase appear when database is nil value, no operation can perform.
var ErrNilDatabase = errors.New("database is nil")

// ErrEmptyStmt appear when try to execute multiple statement with nil/empty slice.
var ErrEmptyStmt = errors.New("no statement to execute")

// ErrEmptyDir appear when migration directory is empty string.
var ErrEmptyDir = errors.New("empty string for migration directory")

// ErrInvalidDir appear when migration directory structure is not correct.
var ErrInvalidDir = errors.New("invalid migration directory structure")
