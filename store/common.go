package store

import (
	"fmt"
	"github.com/g10guang/graduation/constdef"
	"path"
)

type commonStorage struct {
	dirPath string
}

func (*commonStorage) genFileName(fid int64, format ...constdef.ImageFormat) string {
	if len(format) == 0 || format[0] == constdef.InvalidImageFormat {
		return fmt.Sprintf("%d", fid)
	} else {
		return fmt.Sprintf("%d_%d", fid, format[0])
	}
}

func (h *commonStorage) genFilePath(fid int64, format ...constdef.ImageFormat) string {
	return path.Join(h.dirPath, h.genFileName(fid, format...))
}
