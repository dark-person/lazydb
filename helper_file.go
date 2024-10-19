package lazydb

import (
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Create a new database file, WITHOUT apply any schema changes.
// This function will create all necessary directory for given path.
//
// Any incorrect path or creation failed will return an error,
// except database file is already exist.
func createDbFile(path string) error {
	// Prevent invalid database path
	if path == "" {
		return ErrEmptyPath
	}

	// Prevent Not database file
	if filepath.Ext(path) != ".db" {
		return ErrInvalidExt
	}

	// Prevent already existing database file
	if IsFileExist(path) {
		return nil
	}

	// Prevent no directory is created
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	// Start Create Database
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return nil
}

// Check the file is exist or not.
//
// Please note that this function will not check filepath valid or not.
func IsFileExist(path string) bool {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return err == nil
}

// Get largest schema version in given embed.FS, by checking prefix of embed sql filename.
// This function NOT ensure migration can be performed successfully.
//
// Any failure when attempt will return error and version = 0.
func LargestSchemaVer(fs embed.FS, folder string) (uint, error) {
	var maxVersion int

	// Get Directory List 1st
	entries, err := fs.ReadDir(folder)
	if err != nil {
		return 0, err
	}

	// Loop entries to find largest schema version
	for _, entry := range entries {
		// Ensure No nested directories inside
		if entry.IsDir() {
			return 0, ErrInvalidDir
		}

		var v int

		// Scan digit in prefix of sql schema filename
		_, err := fmt.Sscanf(entry.Name(), "%d_", &v)

		// Skip not related file
		if err != nil {
			continue
		}

		// Check version
		if v > maxVersion {
			maxVersion = v
		}
	}

	// Return
	return uint(maxVersion), nil
}

// Copy file from source to destination.
func copyFile(src, dst string) error {
	stat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !stat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}
