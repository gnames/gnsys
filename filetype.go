package gnsys

import (
	"strings"
)

type FileType int

const (
	UnknownFT FileType = iota
	ZipFT              // .zip
	TarFT              // .tar
	TarGzFT            // .tar.gz
	TarXzFt            // .tar.xz
	TarBzFT            //.tar.bz2
)

var ftMap = map[FileType]string{
	UnknownFT: "unknown",
	ZipFT:     "zip",
	TarFT:     "tar",
	TarGzFT:   "tar-gzip",
	TarXzFt:   "tar-xz",
	TarBzFT:   "tar-bz2",
}

func (ft FileType) String() string {
	return ftMap[ft]
}

func GetFileType(file string) FileType {
	switch {
	case strings.HasSuffix(file, ".zip"):
		return ZipFT
	case strings.HasSuffix(file, ".tar"):
		return TarFT
	case strings.HasSuffix(file, ".tar.gz"):
		return TarGzFT
	case strings.HasSuffix(file, ".tar.xz"):
		return TarXzFt
	case strings.HasSuffix(file, ".tar.bz2"):
		return TarBzFT
	default:
		return UnknownFT
	}
}
