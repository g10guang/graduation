package store

import (
	"github.com/g10guang/graduation/constdef"
	"github.com/sirupsen/logrus"
	"github.com/vladimirvivien/gowfs"
	"io"
)

type HdfsStorage struct {
	*commonStorage
	client *gowfs.FileSystem
}

func NewHdfsStorage(webhdfsAddr, user, dir string) *HdfsStorage {
	webgowfsClient, err := gowfs.NewFileSystem(gowfs.Configuration{Addr: webhdfsAddr, User: user})
	if err != nil {
		panic(err)
	}
	s := &HdfsStorage{
		commonStorage: &commonStorage{dirPath: dir},
		client: webgowfsClient,
	}
	return s
}

func (s *HdfsStorage) Write(fid int64, reader io.Reader, format ...constdef.ImageFormat) error {
	filePath := s.genFilePath(fid, format...)
	return s.write(filePath, reader)
}

func (s *HdfsStorage) Read(fid int64, format ...constdef.ImageFormat) (io.Reader, error) {
	filePath := s.genFilePath(fid, format...)
	return s.read(filePath)
}

func (s *HdfsStorage) Delete(fid int64) error {
	go s.delete_(s.genFilePath(fid))
	for _, f := range constdef.ImageFormatList {
		go func(id int64, f constdef.ImageFormat) {
			s.delete_(s.genFilePath(id, f))
		}(fid, f)
	}
	return nil
}

// 这里需要先进行一次 redirect，然后再将数据写入 data node
func (s *HdfsStorage) write(filepath string, reader io.Reader) error {
	isCreate, err := s.client.Create(reader, gowfs.Path{Name: filepath}, true, 0, 2, 0755, 0, "")
	if err != nil {
		logrus.Errorf("webhdfs create file Error: %s", err)
		return err
	}
	logrus.Debugf("webhdfs path: %s create file success: %v", filepath, isCreate)
	return nil
}

func (s *HdfsStorage) read(filePath string) (io.Reader, error) {
	reader, err := s.client.Open(gowfs.Path{Name: filePath}, 0, 0, 0)
	if err != nil {
		logrus.Errorf("webhdfs open path: %s Error: %s", filePath, err)
		return reader, err
	}
	return reader, nil
}

func (s *HdfsStorage) delete_(filePath string) error {
	isDelete, err := s.client.Delete(gowfs.Path{Name: filePath}, false)
	if err != nil {
		logrus.Errorf("webhdfs delete path: %s Error: %s", filePath, err)
		return err
	}
	logrus.Debugf("webhdfs delete path: %s success: %v", filePath, isDelete)
	return nil
}

