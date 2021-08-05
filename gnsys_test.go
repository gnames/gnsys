package gnsys_test

import (
	"testing"

	"github.com/gnames/gnsys"
	"github.com/matryer/is"
)

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
		t.Run(v.name, func(t *testing.T) {
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
		t.Run(v.name, func(t *testing.T) {
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
		t.Run(v.name, func(t *testing.T) {
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
		t.Run(v.name, func(t *testing.T) {
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
		t.Run(v.name, func(t *testing.T) {
			isdir := gnsys.IsDir(v.path)
			is.Equal(v.isDir, isdir)
		})
	}
}
