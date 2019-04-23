package store

import (
	"github.com/g10guang/graduation/constdef"
	"io"
)

// define the interface for storage layer
type Storage interface {
	Write(fid int64, reader io.Reader, format ...constdef.ImageFormat) error
	Read(fid int64, format ...constdef.ImageFormat) (reader io.Reader, err error)
	Delete(fid int64) error
}
