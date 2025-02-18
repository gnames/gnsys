package gnsys

import "fmt"

// ErrFileMissing indicates that a file at the specified path could not be
// found.
type ErrFileMissing struct {
	Path string
}

func (e *ErrFileMissing) Error() string {
	return fmt.Sprintf("file not found '%s'", e.Path)
}

// ErrNotFile indicates that the path provided does not refer to a regular
// file. It might be a directory, symbolic link, or another type of file system
// entry.
type ErrNotFile struct {
	Path string
}

func (e *ErrNotFile) Error() string {
	return fmt.Sprintf("not a file '%s'", e.Path)
}

// ErrNotDir indicates that the path provided does not refer to a directory.
// It might be a file, symbolic link, or another type of file system entry.
type ErrNotDir struct {
	Path string
}

func (e *ErrNotDir) Error() string {
	return fmt.Sprintf("not a directory '%s'", e.Path)
}

// ErrExtract is returned when the extraction of a file (e.g., Zip, Tar file)
// fails. The Path field specifies the file that was being extracted, and the
// Err field contains the underlying error that caused the extraction to fail.
type ErrExtract struct {
	Path string
	Err  error
}

func (e *ErrExtract) Error() string {
	return fmt.Sprintf("extracting '%s' failed: %v", e.Path, e.Err)
}

// ErrDownload is returned when a file download operation fails. The URL field
// specifies the URL that was being downloaded, and the Err field contains the
// underlying error that caused the download to fail.
type ErrDownload struct {
	URL string
	Err error
}

func (e *ErrDownload) Error() string {
	return fmt.Sprintf("cannot download file: %s", e.Err)
}
