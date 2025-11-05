# gnsys

A Go helper package providing utilities for filesystem operations, archive extraction, and file downloads.

## Installation

```bash
go get github.com/gnames/gnsys
```

## Features

- **File & Directory Operations**: Check existence, create directories, copy files, detect file types
- **Archive Extraction**: Extract zip, tar, gzip, xz, bzip2 archives and combinations
- **File Downloads**: HTTP downloads with optional progress bars
- **Path Utilities**: Tilde expansion, path splitting
- **Text File Detection**: Heuristic-based text file identification

## Usage

### File and Directory Operations

```go
import "github.com/gnames/gnsys"

// Check if file exists
exists, err := gnsys.FileExists("path/to/file.txt")

// Check if directory exists and if it's empty
exists, empty, err := gnsys.DirExists("path/to/dir")

// Create directory (with parent directories)
err := gnsys.MakeDir("path/to/new/dir")

// Clean directory (remove all contents or create if absent)
err := gnsys.CleanDir("path/to/dir")

// Get directory state
state := gnsys.GetDirState("path/to/dir")
// Returns: DirAbsent, DirEmpty, DirNotEmpty, NotDir, or Unknown

// Copy file
bytesCopied, err := gnsys.CopyFile("source.txt", "dest.txt")

// Detect if file is text
isText, err := gnsys.IsTextFile("path/to/file")
```

### Path Utilities

```go
// Expand tilde in path
fullPath, err := gnsys.ConvertTilda("~/documents/file.txt")

// Split path into directory, base name, and extension
dir, base, ext := gnsys.SplitPath("/path/to/file.tar.gz")
// dir: "/path/to", base: "file.tar", ext: ".gz"
```

### Archive Extraction

```go
// Extract various archive formats
err := gnsys.ExtractZip("archive.zip", "dest/dir")
err := gnsys.ExtractTar("archive.tar", "dest/dir")
err := gnsys.ExtractGz("file.gz", "dest/dir")
err := gnsys.ExtractXz("file.xz", "dest/dir")
err := gnsys.ExtractTarGz("archive.tar.gz", "dest/dir")
err := gnsys.ExtractTarXz("archive.tar.xz", "dest/dir")
err := gnsys.ExtractTarBz2("archive.tar.bz2", "dest/dir")
```

### File Type Detection

```go
// Get file type based on extension
ft := gnsys.GetFileType("archive.tar.gz")
// Returns: TarGzFT

fmt.Println(ft.String()) // Prints: "tar-gzip"

// Available file types:
// ZipFT, GzFT, XzFT, TarFT, TarGzFT, TarXzFt, TarBzFT, SqlFT, SqliteFT
```

### File Downloads

```go
// Download file without progress bar
filePath, err := gnsys.Download("https://example.com/file.zip", "/dest/dir", false)

// Download file with progress bar
filePath, err := gnsys.Download("https://example.com/file.zip", "/dest/dir", true)

// Check if server is reachable
isReachable := gnsys.Ping("example.com:80", 3) // 3 second timeout
```

## Error Types

The package provides custom error types for better error handling:

- `ErrFileMissing`: File not found at specified path
- `ErrNotFile`: Path is not a regular file
- `ErrNotDir`: Path is not a directory
- `ErrExtract`: Archive extraction failed
- `ErrDownload`: File download failed

## Testing

Run tests using:

```bash
just test
```

Or directly with Go:

```bash
go test -v
```

## License

Licensed under MIT (see LICENSE file)
