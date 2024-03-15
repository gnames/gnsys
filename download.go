package gnsys

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

func Download(url, destDir string) error {
	// Get the filename from the URL
	filename := filepath.Base(url)
	destPath := filepath.Join(destDir, filename)

	// Create the destination file
	outFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Issue HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: server returned status %d", resp.StatusCode)
	}

	// Get the total file size from the content-length header
	contentLength := resp.ContentLength

	// Create the progress bar
	bar := pb.Full.Start64(contentLength)
	barReader := bar.NewProxyReader(resp.Body)

	// Copy data with progress updates
	_, err = io.Copy(outFile, barReader)
	if err != nil {
		return err
	}

	// Finish the progress bar
	bar.Finish()
	return nil
}
