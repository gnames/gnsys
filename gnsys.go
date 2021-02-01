package gnsys

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
)

// MakeDir a directory out of a given unless it already exists.
func MakeDir(dir string) error {
	var err error
	dir, err = ConvertTilda(dir)
	if err != nil {
		return err
	}
	path, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	if path.Mode().IsRegular() {
		return fmt.Errorf("'%s' is a file, not a directory", dir)
	}
	return nil
}

// FileExists checks if a file exists, and that it is a regular file.
func FileExists(f string) (bool, error) {
	path, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false, nil
	}
	if !path.Mode().IsRegular() {
		return false, fmt.Errorf("'%s' is not a regular file, "+
			"delete or move it and try again.", f)
	}
	return true, nil
}

// CleanDir removes all files from a directory.
func CleanDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()

	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// ConvertTilda expands paths with `~/` to an actual home directory.
func ConvertTilda(path string) (string, error) {
	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		home, err := homedir.Dir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[2:])
	}
	return path, nil
}
