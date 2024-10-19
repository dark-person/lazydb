package lazydb

import (
	"database/sql"
	"fmt"
)

// Container for create & execute prepared statements. Read only after creation.
type ParamQuery struct {
	Query string
	Args  []any
}

// Create ParamQuery by given parameters. Value are not able to modify after creation.
func Param(query string, args ...any) ParamQuery {
	return ParamQuery{
		Query: query,
		Args:  args,
	}
}

// Return filled version of query.
func (p *ParamQuery) Filled() string {
	return asFilledQuery(p.Query, p.Args...)
}

// ----------------------------------------------------------------

// Execute given query, by prepared statement.
func (d *LazyDB) Exec(query string, args ...any) (sql.Result, error) {
	if d.db == nil {
		return nil, ErrNilDatabase
	}

	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Exec(args...)
}

// Execute given query, by prepared statement & transaction.
//
// Any failed query will cause a rollback, and return nil []sql.Result.
// Only all queries successful will return valid sql.Result slices.
func (d *LazyDB) ExecMultiple(pQueries []ParamQuery) ([]sql.Result, error) {
	if d.db == nil {
		return nil, ErrNilDatabase
	}

	if pQueries == nil {
		return nil, ErrEmptyStmt
	}

	results := make([]sql.Result, 0)

	// Begin transaction
	tx, err := d.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback() // Any failed tx will cause rollback

	// Exec query with prepare statement
	for _, query := range pQueries {
		stmt, err := tx.Prepare(query.Query)
		if err != nil {
			return nil, fmt.Errorf("Failed to prepare '%s': %v", query.Query, err)
		}

		result, err := stmt.Exec(query.Args...)
		if err != nil {
			return nil, fmt.Errorf("Failed to exec '%s' with (%v): %v", query.Query, query.Args, err)
		}

		// Record result of query
		results = append(results, result)
	}

	// Commit
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Only return results when no errors
	return results, nil
}

// Wrapper for Query() function, using prepared statement.
func (d *LazyDB) Query(query string, args ...any) (*sql.Rows, error) {
	if d.db == nil {
		return nil, ErrNilDatabase
	}

	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.Query(args...)
}

// Wrapper for QueryRow() function, using prepared statement.
func (d *LazyDB) QueryRow(query string, args ...any) (*sql.Row, error) {
	if d.db == nil {
		return nil, ErrNilDatabase
	}

	stmt, err := d.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	return stmt.QueryRow(args...), nil
}
