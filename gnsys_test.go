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
