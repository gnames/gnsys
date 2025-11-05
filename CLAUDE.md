# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`gnsys` is a Go helper package for filesystem operations, focusing on file/directory management, archive extraction, and file downloads. It's designed as a utility library to be imported by other Go projects.

## Development Commands

### Testing
```bash
# Run all tests
go test

# Run tests with verbose output
go test -v

# Run a specific test
go test -run TestFunctionName

# Run tests with coverage
go test -cover
```

### Building
```bash
# Build (if needed for validation)
go build

# Install as a module (for use in other projects)
go mod tidy
```

### Linting
```bash
# Format code
go fmt ./...

# Run go vet
go vet ./...
```

## Architecture

### Core Components

**gnsys.go**: Core filesystem operations
- `DirState` enum: Represents directory states (Unknown, NotDir, DirAbsent, DirEmpty, DirNotEmpty)
- Directory operations: `GetDirState`, `MakeDir`, `CleanDir`, `DirExists`
- File operations: `FileExists`, `IsFile`, `IsDir`, `IsTextFile`, `CopyFile`
- Path utilities: `ConvertTilda` (expands `~/`), `SplitPath` (breaks path into dir/base/ext)

**extract.go**: Archive extraction functionality
- Defines `Extractor` function type: `func(src, dst string) error`
- Extraction functions for multiple formats:
  - `ExtractZip`: .zip archives
  - `ExtractTar`: .tar archives
  - `ExtractGz`: .gz compressed files
  - `ExtractTarGz`: .tar.gz archives
  - `ExtractTarBz2`: .tar.bz2 archives
  - `ExtractTarXz`: .tar.xz archives
- Internal `untar` helper for tar-based formats

**download.go**: Network operations
- `Ping(host, seconds)`: TCP connectivity check (format: "host:port")
- `Download(url, destDir, showProgress)`: HTTP file download with optional progress bar using `github.com/cheggaaa/pb/v3`

**filetype.go**: File type detection
- `FileType` enum with constants: ZipFT, GzFT, TarFT, TarGzFT, TarXzFt, TarBzFT, SqlFT, SqliteFT
- `GetFileType(file)`: Determines file type based on extension (suffix matching)

**errors.go**: Custom error types
- `ErrFileMissing`: File not found at path
- `ErrNotFile`: Path is not a regular file
- `ErrNotDir`: Path is not a directory
- `ErrExtract`: Archive extraction failure
- `ErrDownload`: Download operation failure

### Design Patterns

1. **Error Wrapping**: Custom error types wrap underlying errors with context (path/URL)
2. **Functional Types**: `Extractor` type allows flexibility in extraction implementations
3. **Progress Reporting**: Download function supports optional progress bar display
4. **Path Safety**: `ConvertTilda` ensures paths with `~/` work correctly

### Testing

- Tests use `testify/assert` and `matryer/is` for assertions
- Test data stored in `testdata/` directory
- Network test downloads from `opendata.globalnames.org` (may require internet connection)
- Tests create temporary directories for validation

## Important Notes

- **Progress Bar**: The `Download` function with `showProgress=true` uses `pb.CleanOnFinish` to clear the progress bar after completion (see download.go:66)
- **Tar Extraction**: The `untar` function validates extracted directory is not empty to detect bad tar files (see extract.go:215-221)
- **Text File Detection**: `IsTextFile` checks only first 20 lines and considers null bytes, non-printable character ratio (>30%), and line length (see gnsys.go:129-178)
- **File Type Order**: When detecting file types, check for multi-extension formats (`.tar.gz`, `.tar.xz`, `.tar.bz2`) before single extensions (`.gz`) to avoid false positives
