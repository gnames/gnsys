package gnsys

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/cheggaaa/pb/v3"
)

// Ping checks if a server is reachable.
// Host should be in format "host:port" (eg "google.com:80")
func Ping(host string, seconds int) bool {
	_, err := net.DialTimeout(
		"tcp",
		host,
		time.Second*time.Duration(seconds),
	)

	return err == nil
}

// Download copies remote file to local drive. It provides the name
// of downloaded file and error as output. Supports http://, https://,
// and file:// URLs.
func Download(rawURL, destDir string, showProgress bool) (string, error) {
	// Parse the URL to determine the scheme
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", &ErrDownload{URL: rawURL, Err: err}
	}

	// Get the filename from the URL
	filename := filepath.Base(parsedURL.Path)
	destPath := filepath.Join(destDir, filename)

	// Create the destination file
	outFile, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	var reader io.Reader
	var contentLength int64

	switch parsedURL.Scheme {
	case "file":
		// Handle local file copy
		srcPath := parsedURL.Path
		srcFile, err := os.Open(srcPath)
		if err != nil {
			return "", &ErrDownload{URL: rawURL, Err: err}
		}
		defer srcFile.Close()

		// Get file size for progress bar
		fileInfo, err := srcFile.Stat()
		if err != nil {
			return "", &ErrDownload{URL: rawURL, Err: err}
		}
		contentLength = fileInfo.Size()
		reader = srcFile

	case "http", "https":
		// Issue HTTP GET request
		resp, err := http.Get(rawURL)
		if err != nil {
			return "", &ErrDownload{URL: rawURL, Err: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf(
				"download failed: server returned status %d",
				resp.StatusCode,
			)
			return "", &ErrDownload{URL: rawURL, Err: err}
		}

		contentLength = resp.ContentLength
		reader = resp.Body

	default:
		err = fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
		return "", &ErrDownload{URL: rawURL, Err: err}
	}

	if showProgress {
		// Create the progress bar
		bar := pb.Full.Start64(contentLength)
		bar.Set(pb.CleanOnFinish, true)
		reader = bar.NewProxyReader(reader)

		// Finish the progress bar
		defer bar.Finish()
	}

	_, err = io.Copy(outFile, reader)
	if err != nil {
		return "", &ErrDownload{URL: rawURL, Err: err}
	}

	return destPath, nil
}
