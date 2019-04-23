package store

import (
	"bytes"
	"errors"
	"github.com/g10guang/graduation/constdef"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// concurrency safe
type LocalStorage struct {
	*commonStorage
	dirPath string
}

func NewLocalStorage() *LocalStorage {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		panic(errors.New("GOPATH not exists in env"))
	}
	dir := path.Join(goPath, "src/github.com/g10guang/graduation/oss")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.Mkdir(dir, 0777); err != nil {
			panic(err)
		}
	}
	h := &LocalStorage{
		commonStorage: &commonStorage{dirPath: dir},
		dirPath: dir,
	}
	return h
}

func (h *LocalStorage) Write(fid int64, reader io.Reader, format ...constdef.ImageFormat) error {
	filePath := h.genFilePath(fid, format...)
	return h.write(filePath, reader)
}

func (h *LocalStorage) Read(fid int64, format ...constdef.ImageFormat) (reader io.Reader, err error) {
	filePath := h.genFilePath(fid, format...)
	return h.read(filePath)
}

// 删除需要将其他格式一并删除
func (h *LocalStorage) Delete(fid int64) error {
	go os.Remove(h.genFilePath(fid))
	for _, f := range constdef.ImageFormatList {
		go os.Remove(h.genFilePath(fid, f))
	}
	return nil
}

func (h *LocalStorage) write(path string, reader io.Reader) error {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		logrus.Errorf("local write read from io.Reader Error: %s", err)
		return err
	}
	err = ioutil.WriteFile(path, b, 0666)
	if err != nil {
		logrus.Errorf("write %s Error: %s", path, err)
	}
	return err
}

func (h *LocalStorage) read(path string) (io.Reader, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Errorf("read %s Error: %s", path, err)
	}
	return bytes.NewReader(b), err
}
