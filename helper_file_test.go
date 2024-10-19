package lazydb

import (
	"embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test_schema/with_others/*
var fsOtherTest embed.FS

//go:embed test_schema/skipped/*
var fsSkippedTest embed.FS

func TestCreateFile(t *testing.T) {
	// Test Case Type
	type testCase struct {
		dbPath  string // Database path, as this function only accept this
		wantErr bool
	}

	// Prepare a existed database
	existPath := filepath.Join(t.TempDir(), "existed.db")
	f, _ := os.Create(existPath)
	f.Close()

	// Prepare Test Case, with direct declare with `&LazyDB{}`
	tests := []testCase{
		// Valid Path, not exist database
		{filepath.Join(t.TempDir(), "new.db"), false},
		// Valid Path, existed database
		{existPath, false},
		// Valid path, with directory not exist
		{filepath.Join(t.TempDir(), "not_exist", "new.db"), false},
		// Invalid path, with database path incorrect
		{filepath.Join(t.TempDir(), "invalid", "invalid.txt"), true},
		// Empty filepath
		{"", true},
	}

	for idx, tt := range tests {
		// Run Create
		err := createDbFile(tt.dbPath)
		if tt.wantErr {
			assert.Errorf(t, err, "Case %d - Expected return error.", idx+1)
		} else {
			assert.Nilf(t, err, "Case %d - Expect no error return, but got %v", idx+1, err)
		}
	}
}

func TestIsFileExist(t *testing.T) {
	// Prepare temp directory
	dir := t.TempDir()

	// Case 1: Exist file
	path1 := filepath.Join(dir, "test1")
	tmp, _ := os.Create(path1)
	defer tmp.Close()

	// Case 2: Not exist file
	path2 := filepath.Join(dir, "test2")

	// Prepare test cases
	tests := []struct {
		path string
		want bool
	}{
		{path1, true},
		{path2, false},
		{"???", false}, // Invalid filepath
	}

	// Start Test
	for idx, tt := range tests {
		got := IsFileExist(tt.path)
		assert.EqualValuesf(t, tt.want, got, "Case %d: Unexpected result", idx)
	}
}

func TestLargestSchemaVer(t *testing.T) {
	type testCase struct {
		fs      embed.FS
		folder  string
		want    uint
		wantErr bool
	}

	tests := []testCase{
		{fsNormalTestV3, "test_schema/normal_v3", 3, false},
		{fsOtherTest, "test_schema/with_others", 2, false},
		{fsSkippedTest, "test_schema/skipped", 3, false},

		{fsNormalTestV3, "abc", 0, true},
		{embed.FS{}, "test_schema/normal_v3", 0, true},
	}

	for idx, tt := range tests {
		got, err := LargestSchemaVer(tt.fs, tt.folder)

		if tt.wantErr {
			assert.NotNilf(t, err, "Case %d: err should be non-nil", idx)
			continue
		}

		assert.Nilf(t, err, "Case %d: err should be nil, but %v", idx, err)
		assert.EqualValuesf(t, tt.want, got, "Case %d: unexpected version number.", idx)
	}
}
