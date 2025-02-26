package gnsys

import (
	"archive/tar"
	"archive/zip"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/ulikunitz/xz"
)

type Extractor func(src, dst string) error

// ExtractZip extracts a zip archive located at srcPath to the destination
// directory dstDir.
func ExtractZip(srcPath, dstDir string) error {
	exists, _ := FileExists(srcPath)
	if !exists {
		return &ErrFileMissing{Path: srcPath}
	}

	// Open the zip file for reading.
	r, err := zip.OpenReader(srcPath)
	if err != nil {
		return &ErrExtract{Path: srcPath, Err: err}
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dstDir, f.Name)
		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return &ErrExtract{Path: fpath, Err: err}
		}

		// If it's a directory, move on to the next entry.
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		// Open the file within the zip.
		rc, err := f.Open()
		if err != nil {
			return &ErrExtract{Path: fpath, Err: err}
		}
		defer rc.Close()

		// Create a file in the filesystem.
		outFile, err := os.OpenFile(
			fpath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			f.Mode(),
		)
		if err != nil {
			return &ErrExtract{Path: fpath, Err: err}
		}
		defer outFile.Close()

		// Copy the contents of the file from the zip to the new file.
		_, err = io.Copy(outFile, rc)
		if err != nil {
			return &ErrExtract{Path: fpath, Err: err}
		}
	}

	return nil
}

// ExtractTar extracts a tar archive located at srcPath to the destination
// directory dstDir.
func ExtractTar(srcPath, dstDir string) error {
	// Open the tar archive for reading.
	file, err := os.Open(srcPath)
	if err != nil {
		return &ErrExtract{Path: srcPath, Err: err}
	}
	defer file.Close()

	tr := tar.NewReader(file)
	return untar(tr, srcPath, dstDir)
}

func ExtractGz(srcPath, dstDir string) error {
	// Open the .gz archive for reading.
	file, err := os.Open(srcPath)
	if err != nil {
		return &ErrExtract{Path: srcPath, Err: err}
	}
	defer file.Close()

	// Create a new gzip reader.
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return &ErrExtract{Path: srcPath, Err: err}
	}
	defer gzReader.Close()

	// Determine the destination file name.
	dstFileName := filepath.Base(srcPath)
	dstFileName = dstFileName[:len(dstFileName)-3] // Remove ".gz" extension
	dstPath := filepath.Join(dstDir, dstFileName)

	// Create the destination file.
	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return &ErrExtract{Path: dstPath, Err: err}
	}
	defer dstFile.Close()

	// Copy the file contents from the gzip reader to the destination file.
	if _, err := io.Copy(dstFile, gzReader); err != nil {
		return &ErrExtract{Path: dstPath, Err: err}
	}

	return nil
}

// ExtractTarGz extracts a tar.gz archive located at srcPath to the destination
// directory dstDir.
func ExtractTarGz(srcPath, dstDir string) error {
	// Open the .tar.gz archive for reading.
	file, err := os.Open(srcPath)
	if err != nil {
		return &ErrExtract{Path: srcPath, Err: err}
	}
	defer file.Close()

	// Create a new gzip reader.
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return &ErrExtract{Path: srcPath, Err: err}
	}
	defer gzReader.Close()

	// Create a new tar reader from the gzip reader.
	tr := tar.NewReader(gzReader)
	return untar(tr, srcPath, dstDir)
}

// ExtractTarBz2 extracts a tar.bz archive located at srcPath to the destination
// directory dstDir.
func ExtractTarBz2(srcPath, dstDir string) error {
	// Open the .tar.gz archive for reading.
	file, err := os.Open(srcPath)
	if err != nil {
		return &ErrExtract{Path: srcPath, Err: err}
	}
	defer file.Close()

	// Create a new bz2 reader.
	bzReader := bzip2.NewReader(file)

	// Create a new tar reader from the gzip reader.
	tr := tar.NewReader(bzReader)
	return untar(tr, srcPath, dstDir)
}

// ExtractTarXz extracts a tar.xz archive located at srcPath to the destination
// directory dstDir.
func ExtractTarXz(srcPath, dstDir string) error {
	// Open the .tar.gz archive for reading.
	file, err := os.Open(srcPath)
	if err != nil {
		return &ErrExtract{Path: srcPath, Err: err}
	}
	defer file.Close()

	xzReader, err := xz.NewReader(file)
	if err != nil {
		return &ErrExtract{Path: srcPath, Err: err}
	}

	// Create a new tar reader from the gzip reader.
	tr := tar.NewReader(xzReader)
	return untar(tr, srcPath, dstDir)
}

func untar(tarReader *tar.Reader, srcPath, dstDir string) error {
	var writer *os.File
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return &ErrExtract{Path: srcPath, Err: err}
		}

		// Get the individual filepath from the header.
		filepath := filepath.Join(dstDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Handle directory.
			err = os.MkdirAll(filepath, os.FileMode(header.Mode))
			if err != nil {
				return &ErrExtract{Path: srcPath, Err: err}
			}
		case tar.TypeReg:
			// Handle regular file.
			writer, err = os.Create(filepath)
			if err != nil {
				return &ErrExtract{Path: srcPath, Err: err}
			}
			io.Copy(writer, tarReader)
			writer.Close()
		default:
			return &ErrExtract{Path: srcPath, Err: err}
		}
	}
	state := GetDirState(dstDir)
	if state == DirEmpty {
		return &ErrExtract{
			Path: dstDir,
			Err:  errors.New("bad tar file"),
		}
	}
	return nil
}
