package loader

import (
	"bytes"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/dal/redis"
	"github.com/g10guang/graduation/store"
	"github.com/sirupsen/logrus"
	"io/ioutil"
)

const LoaderName_FileContent = "file_content_loader"

type FileContentLoader struct {
	fid     int64
	storage store.Storage
	format  constdef.ImageFormat
}

func NewFileContentLoader(fid int64, storage store.Storage, format constdef.ImageFormat) *FileContentLoader {
	return &FileContentLoader{
		fid:     fid,
		storage: storage,
		format:  format,
	}
}

func (l *FileContentLoader) GetName() string {
	return LoaderName_FileContent
}

// 先从 redis 缓存中获取，没有再到 HDFS
func (l *FileContentLoader) Run() (interface{}, error) {
	// redis 非核心逻辑允许失败
	b, err := redis.ContentRedis.Get(l.fid, l.format)
	if err == nil {
		logrus.Debugf("fid: %d content cache hit", l.fid)
		return bytes.NewReader(b), nil
	}
	r, err := l.storage.Read(l.fid, l.format)
	if err != nil {
		return nil, err
	}
	b, err = ioutil.ReadAll(r)
	go l.saveRedis(b)
	return bytes.NewReader(b), nil
}

func (l *FileContentLoader) saveRedis(b []byte) {
	redis.ContentRedis.Set(l.fid, b, l.format)
}
