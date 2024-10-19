package lazydb

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Specific test for int values boundary.
func TestAsFilledQueryInt(t *testing.T) {
	type testCase struct {
		sql  string
		args []any
		want string
	}

	const minInt8 int8 = math.MinInt8
	const maxInt8 int8 = math.MaxInt8
	const minInt16 int16 = math.MinInt16
	const maxInt16 int16 = math.MaxInt16
	const minInt32 int32 = math.MinInt32
	const maxInt32 int32 = math.MaxInt32
	const minInt64 int64 = math.MinInt64
	const maxInt64 int64 = math.MaxInt64

	tests := []testCase{
		{`SELECT * FROM test WHERE col1 = ?`, []any{minInt8}, `SELECT * FROM test WHERE col1 = -128`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{maxInt8}, `SELECT * FROM test WHERE col1 = 127`},

		{`SELECT * FROM test WHERE col1 = ?`, []any{minInt16}, `SELECT * FROM test WHERE col1 = -32768`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{maxInt16}, `SELECT * FROM test WHERE col1 = 32767`},

		{`SELECT * FROM test WHERE col1 = ?`, []any{minInt32}, `SELECT * FROM test WHERE col1 = -2147483648`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{maxInt32}, `SELECT * FROM test WHERE col1 = 2147483647`},

		{`SELECT * FROM test WHERE col1 = ?`, []any{minInt64}, `SELECT * FROM test WHERE col1 = -9223372036854775808`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{maxInt64}, `SELECT * FROM test WHERE col1 = 9223372036854775807`},
	}

	for idx, tt := range tests {
		result := asFilledQuery(tt.sql, tt.args...)
		assert.EqualValuesf(t, tt.want, result, "Case %d: incorrect result on int value", idx)
	}
}

// Specific test for uint values boundary.
func TestAsFilledQueryUInt(t *testing.T) {
	type testCase struct {
		sql  string
		args []any
		want string
	}

	const minUint8 uint8 = 0
	const maxUint8 uint8 = 255
	const minUint16 uint16 = 0
	const maxUint16 uint16 = 65535
	const minUint32 uint32 = 0
	const maxUint32 uint32 = 4294967295
	const minUint64 uint64 = 0
	const maxUint64 uint64 = 18446744073709551615

	tests := []testCase{
		{`SELECT * FROM test WHERE col1 = ?`, []any{minUint8}, `SELECT * FROM test WHERE col1 = 0`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{maxUint8}, `SELECT * FROM test WHERE col1 = 255`},

		{`SELECT * FROM test WHERE col1 = ?`, []any{minUint16}, `SELECT * FROM test WHERE col1 = 0`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{maxUint16}, `SELECT * FROM test WHERE col1 = 65535`},

		{`SELECT * FROM test WHERE col1 = ?`, []any{minUint32}, `SELECT * FROM test WHERE col1 = 0`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{maxUint32}, `SELECT * FROM test WHERE col1 = 4294967295`},

		{`SELECT * FROM test WHERE col1 = ?`, []any{minUint64}, `SELECT * FROM test WHERE col1 = 0`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{maxUint64}, `SELECT * FROM test WHERE col1 = 18446744073709551615`},
	}

	for idx, tt := range tests {
		result := asFilledQuery(tt.sql, tt.args...)
		assert.EqualValuesf(t, tt.want, result, "Case %d: incorrect result on uint value", idx)
	}
}

func TestAsFilledQuery(t *testing.T) {
	type testCase struct {
		sql  string
		args []any
		want string
	}

	tests := []testCase{
		// Single value for type check
		{`SELECT * FROM test WHERE col1 = ?`, []any{"abc"}, `SELECT * FROM test WHERE col1 = "abc"`},

		{`SELECT * FROM test WHERE col1 = ?`, []any{123}, `SELECT * FROM test WHERE col1 = 123`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{int8(123)}, `SELECT * FROM test WHERE col1 = 123`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{int16(123)}, `SELECT * FROM test WHERE col1 = 123`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{int32(123)}, `SELECT * FROM test WHERE col1 = 123`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{int64(123)}, `SELECT * FROM test WHERE col1 = 123`},

		{`SELECT * FROM test WHERE col1 = ?`, []any{uint(123)}, `SELECT * FROM test WHERE col1 = 123`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{uint8(123)}, `SELECT * FROM test WHERE col1 = 123`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{uint16(123)}, `SELECT * FROM test WHERE col1 = 123`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{uint32(123)}, `SELECT * FROM test WHERE col1 = 123`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{uint64(123)}, `SELECT * FROM test WHERE col1 = 123`},

		{`SELECT * FROM test WHERE col1 = ?`, []any{float32(123.45)}, `SELECT * FROM test WHERE col1 = 123.45`},
		{`SELECT * FROM test WHERE col1 = ?`, []any{float64(123.45)}, `SELECT * FROM test WHERE col1 = 123.45`},

		// Multiple values
		{`SELECT * FROM test WHERE col1 = ? AND col2 = ?`, []any{123, "abc"}, `SELECT * FROM test WHERE col1 = 123 AND col2 = "abc"`},
		{`SELECT * FROM test WHERE col1 = ? AND col2 = ?`, []any{123}, `SELECT * FROM test WHERE col1 = 123 AND col2 = ?`},

		// No values
		{`SELECT * FROM test WHERE col1 = ? AND col2 = ?`, []any{}, `SELECT * FROM test WHERE col1 = ? AND col2 = ?`},
	}

	for idx, tt := range tests {
		result := asFilledQuery(tt.sql, tt.args...)
		assert.EqualValuesf(t, tt.want, result, "Case %d: incorrect result", idx)
	}

	// Panic case
	assert.Panics(t, func() { asFilledQuery(`SELECT * FROM test WHERE col1 = ? AND col2 = ?`, nil) })
}
