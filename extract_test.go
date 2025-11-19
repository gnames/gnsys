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

func TestExtractBZ2(t *testing.T) {
	assert := assert.New(t)
	bzFile := filepath.Join("testdata", "text.txt.bz2")
	tempDir, err := os.MkdirTemp("", "gnsys-test")
	assert.Nil(err)

	err = gnsys.ExtractBz2(bzFile, tempDir)
	assert.Nil(err)

	exists, err := gnsys.FileExists(filepath.Join(tempDir, "text.txt"))
	assert.Nil(err)
	assert.True(exists)

	err = os.RemoveAll(tempDir)
	assert.Nil(err)
}

func TestExtractXZ(t *testing.T) {
	assert := assert.New(t)
	xzFile := filepath.Join("testdata", "text.txt.xz")
	tempDir, err := os.MkdirTemp("", "gnsys-test")
	assert.Nil(err)

	err = gnsys.ExtractXz(xzFile, tempDir)
	assert.Nil(err)

	exists, err := gnsys.FileExists(filepath.Join(tempDir, "text.txt"))
	assert.Nil(err)
	assert.True(exists)

	err = os.RemoveAll(tempDir)
	assert.Nil(err)
}
