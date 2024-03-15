package gnsys

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/cheggaaa/pb/v3"
)

func Download(url, destDir string, showProgress bool) (string, error) {
	// Get the filename from the URL
	filename := filepath.Base(url)
	destPath := filepath.Join(destDir, filename)

	// Create the destination file
	outFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// Issue HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(
			"download failed: server returned status %d",
			resp.StatusCode,
		)
		return "", err
	}

	// Get the total file size from the content-length header
	contentLength := resp.ContentLength

	var reader io.Reader
	if showProgress {
		// Create the progress bar
		bar := pb.Full.Start64(contentLength)
		reader = bar.NewProxyReader(resp.Body)

		// Finish the progress bar
		defer bar.Finish()
	} else {
		// Copy data without progress bar
		reader = resp.Body
	}
	_, err = io.Copy(outFile, reader)
	// Copy data with progress updates
	if err != nil {
		return "", err
	}

	return destPath, nil
}
