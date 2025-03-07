package gnsys_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/gnames/gnsys"
	"github.com/matryer/is"
	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	assert := assert.New(t)
	ping := gnsys.Ping("google.com:80", 3)
	assert.True(ping)

	ping = gnsys.Ping("notAserver:80", 3)
	assert.False(ping)
}

func TestConvertTilda(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		name, input            string
		outputSameSize, errNil bool
	}{
		{"no tilda", "/somedir", true, true},
		{"tilda", "~/somedir", false, true},
	}

	for _, v := range tests {
		t.Run(v.name, func(_ *testing.T) {
			path, err := gnsys.ConvertTilda(v.input)
			is.Equal(len(v.input) == len(path), v.outputSameSize)
			is.Equal(v.errNil, err == nil)
		})
	}
}

func TestFileExists(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		name, path         string
		fileExists, errNil bool
	}{
		{"file exists", "testdata/text.txt", true, true},
		{"file does not exist", "testdata/text2.txt", false, true},
		{"is dir", "testdata", false, false},
	}
	for _, v := range tests {
		t.Run(v.name, func(_ *testing.T) {
			exists, err := gnsys.FileExists(v.path)
			is.Equal(v.fileExists, exists)
			is.Equal(v.errNil, err == nil)
		})
	}
}

func TestDirExists(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		name, path                  string
		dirExists, dirEmpty, errNil bool
	}{
		{"dir exists notempty", "testdata", true, false, true},
		{"dir exists not", "testdata/nodir", false, false, false},
	}
	for _, v := range tests {
		t.Run(v.name, func(_ *testing.T) {
			exists, empty, err := gnsys.DirExists(v.path)
			is.Equal(v.dirExists, exists)
			is.Equal(v.dirEmpty, empty)
			is.Equal(v.errNil, err == nil)
		})
	}
}

func TestIsFile(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		name, path string
		isFile     bool
	}{
		{"is dir", "testdata", false},
		{"is file", "testdata/text.txt", true},
		{"is not file", "testdata/nofile", false},
	}
	for _, v := range tests {
		t.Run(v.name, func(_ *testing.T) {
			isfile := gnsys.IsFile(v.path)
			is.Equal(v.isFile, isfile)
		})
	}
}

func TestIsDir(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		name, path string
		isDir      bool
	}{
		{"is dir", "testdata", true},
		{"is not dir", "testdata/text.txt", false},
		{"is not dir", "testdata/nodir", false},
	}
	for _, v := range tests {
		t.Run(v.name, func(_ *testing.T) {
			isdir := gnsys.IsDir(v.path)
			is.Equal(v.isDir, isdir)
		})
	}
}

func TestGetDirState(t *testing.T) {
	is := is.New(t)
	makeEmptyDir(t)
	tests := []struct {
		name, path string
		state      gnsys.DirState
	}{
		{"is dir", "testdata", gnsys.DirNotEmpty},
		{"is empty dir", "testdata/empty_dir", gnsys.DirEmpty},
		{"is not dir", "testdata/text.txt", gnsys.NotDir},
		{"absent", "testdata/absent_from_tests", gnsys.DirAbsent},
	}
	for _, v := range tests {
		t.Run(v.name, func(_ *testing.T) {
			state := gnsys.GetDirState(v.path)
			is.Equal(v.state, state)
		})
	}
}

func makeEmptyDir(t *testing.T) {
	assert := assert.New(t)
	dir := filepath.Join("testdata/empty_dir")
	err := os.RemoveAll(dir)
	// err would be nil if dir does not exist
	assert.Nil(err)
	err = os.Mkdir(dir, 0775)
	assert.Nil(err)
}

func TestDownload(t *testing.T) {
	is := is.New(t)
	url := "http://opendata.globalnames.org/dwca/183-sherborn.tar.gz"
	path := os.TempDir()
	filePath, err := gnsys.Download(url, path, false)
	is.NoErr(err)
	exists, _ := gnsys.FileExists(filePath)
	is.True(exists)
	err = os.Remove(filePath)
	is.NoErr(err)
}

func TestSplitPath(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, path, dir, file, ext string
	}{
		{"empty", "", "", "", ""},
		{"root1", "/", "/", "", ""},
		{"any", "c", "c", "", ""},
		{"dir1", "/one/two/three/", "/one/two/three", "", ""},
		{"dir3", "~/one/two/three/", "~/one/two/three", "", ""},
		{"file1", "~/one/two/three", "~/one/two", "three", ""},
		{"file2", "./one/two/three.txt", "one/two", "three", ".txt"},
		{"file3", "../one/two/three.txt.gz", "../one/two", "three.txt", ".gz"},
	}

	for _, v := range tests {
		d, f, e := gnsys.SplitPath(v.path)
		fmt.Printf("%s, %s, %s", d, f, e)
		assert.Equal(v.dir, d, v.msg+":dir")
		assert.Equal(v.file, f, v.msg+":file")
		assert.Equal(v.ext, e, v.msg+":ext")
	}
}

func TestIsTextFile(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		msg, file string
		isErr     bool
		isText    bool
	}{
		{"text", "text.txt", false, true},
		{"gz", "text.txt.gz", false, false},
		{"none", "ddd-no-such-file", true, false},
	}

	for _, v := range tests {
		path := filepath.Join("testdata", v.file)
		isText, err := gnsys.IsTextFile(path)
		assert.Equal(v.isErr, err != nil)
		assert.Equal(v.isText, isText)
	}
}
