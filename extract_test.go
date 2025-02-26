package gnsys_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnsys"
	"github.com/stretchr/testify/assert"
)

func TestExtractGZ(t *testing.T) {
	assert := assert.New(t)
	gzFile := filepath.Join("testdata", "text.txt.gz")
	tempDir, err := os.MkdirTemp("", "gnsys-test")
	assert.Nil(err)

	err = gnsys.ExtractGz(gzFile, tempDir)
	assert.Nil(err)

	exists, err := gnsys.FileExists(filepath.Join(tempDir, "text.txt"))
	assert.Nil(err)
	assert.True(exists)

	err = os.RemoveAll(tempDir)
	assert.Nil(err)
}
