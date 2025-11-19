package gnsys

import (
	"strings"
)

type FileType int

const (
	UnknownFT FileType = iota
	ZipFT              // .zip
	GzFT               // .gz
	TarFT              // .tar
	TarGzFT            // .tar.gz
	TarXzFt            // .tar.xz
	TarBzFT            // .tar.bz2
	Bz2FT              // .bz2
	XzFT               // .xz
	SqlFT              // .sql
	SqliteFT           // .sqlite
)

var ftMap = map[FileType]string{
	UnknownFT: "unknown",
	ZipFT:     "zip",
	GzFT:      "gz",
	TarFT:     "tar",
	TarGzFT:   "tar-gzip",
	TarXzFt:   "tar-xz",
	TarBzFT:   "tar-bz2",
	Bz2FT:     "bz2",
	XzFT:      "xz",
	SqlFT:     "sql",
	SqliteFT:  "sqlite",
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
	case strings.HasSuffix(file, ".gz"):
		return GzFT
	case strings.HasSuffix(file, ".tar.xz"):
		return TarXzFt
	case strings.HasSuffix(file, ".xz"):
		return XzFT
	case strings.HasSuffix(file, ".tar.bz2"):
		return TarBzFT
	case strings.HasSuffix(file, ".bz2"):
		return Bz2FT
	case strings.HasSuffix(file, ".sql"):
		return SqlFT
	case strings.HasSuffix(file, ".sqlite"):
		return SqliteFT
	default:
		return UnknownFT
	}
}
