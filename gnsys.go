package gnsys

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DirState represents the state of a directory.
type DirState int

const (
	Unknown DirState = iota
	NotDir
	DirAbsent
	DirEmpty
	DirNotEmpty
)

// String returns a string representation of the DirState.
func (d DirState) String() string {
	switch d {
	case NotDir:
		return "NotDir"
	case DirAbsent:
		return "DirAbsent"
	case DirEmpty:
		return "DirEmpty"
	case DirNotEmpty:
		return "DirNotEmpty"
	}
	return "Unknown"
}

// GetDirState returns the state of a directory.
func GetDirState(dir string) DirState {
	st, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return DirAbsent
	}
	if st == nil {
		return NotDir
	}
	if !st.Mode().IsDir() {
		return NotDir
	}

	d, err := os.Open(dir)
	if err != nil {
		return Unknown
	}
	defer d.Close()

	_, err = d.Readdirnames(1)
	if err == io.EOF {
		return DirEmpty
	} else if err != nil {
		return Unknown
	}
	return DirNotEmpty
}

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
	if path == nil {
		return nil
	}
	if path.Mode().IsRegular() {
		return &ErrNotDir{Path: path.Name()}
	}
	return nil
}

// FileExists checks if a file exists, and that it is a regular file.
func FileExists(f string) (bool, error) {
	path, err := os.Stat(f)
	if os.IsNotExist(err) {
		return false, nil
	}
	if path == nil {
		return false, nil
	}
	if !path.Mode().IsRegular() {
		return false, &ErrNotFile{Path: f}
	}
	return true, nil
}

// DirExists checks if directory exists and if it is empty
func DirExists(path string) (exists bool, empty bool, err error) {
	st, err := os.Stat(path)
	if os.IsNotExist(err) || st.Mode().IsRegular() {
		return false, false, err
	}

	d, err := os.Open(path)
	if err != nil {
		return false, false, err
	}
	defer d.Close()

	_, err = d.Readdirnames(1)
	if err == io.EOF {
		return true, true, nil
	} else if err != nil {
		return false, false, err
	}
	return true, false, nil
}

func IsFile(path string) bool {
	res, _ := FileExists(path)
	return res
}

func IsDir(path string) bool {
	res, _, _ := DirExists(path)
	return res
}

// CleanDir removes all files from a directory or creates the directory if
// it is absent.
func CleanDir(dir string) error {
	exists, _, err := DirExists(dir)
	if err != nil {
		return err
	}
	if !exists {
		return MakeDir(dir)
	}

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
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(home, path[2:])
	}
	return path, nil
}

// SplitPath breaks path into directory, file name and extension.
func SplitPath(path string) (dir, base, ext string) {
	if len(path) < 2 {
		return path, "", ""
	}

	e := path[len(path)-1]
	if e == '/' {
		return path[:len(path)-1], "", ""
	}

	dir = filepath.Dir(path)
	base = filepath.Base(path)
	ext = filepath.Ext(path)
	base = base[:len(base)-len(ext)] // Remove extension from base
	return
}
