package mysql

import (
	"errors"

	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

type FileInfoMySql struct {
	conn *gorm.DB
}

func NewFileInfoMySql() *FileInfoMySql {
	var err error
	h := &FileInfoMySql{}
	driver, url := getFileInfoMySql()
	logrus.Infof("driver: %s url: %s", driver, url)
	h.conn, err = gorm.Open(driver, url)
	if err != nil {
		logrus.Panicf("Create FileInfo MySQL connection Error: %s", err)
		panic(err)
	}
	h.conn = h.conn.Table("file")
	// In test env, print SQL gorm execute
	if !tools.IsProductEnv() {
		h.conn = h.conn.Debug()
	}
	return h
}

func getFileInfoMySql() (string, string) {
	return "mysql", constdef.MySqlUrl
}

// 写需要涉及到事务，所以由外部传递 connection
func (h *FileInfoMySql) Save(conn *gorm.DB, file *model.File) error {
	if conn == nil {
		conn = h.conn
	}
	err := h.conn.Create(file).Error
	if err != nil {
		logrus.Errorf("Save Error: %s", err)
	}
	return err
}

// 删需要涉及到事务，所以由外部传递 connection
func (h *FileInfoMySql) Delete(conn *gorm.DB, fids []int64, uid int64) (err error) {
	if conn == nil {
		conn = h.conn
	}
	if err = conn.Where("fid IN (?) AND uid = ?", fids, uid).Delete(nil).Error; err != nil {
		logrus.Errorf("Delete Error: %s", err)
	}
	return
}

func (h *FileInfoMySql) Get(fid int64) (meta model.File, err error) {
	if err = h.conn.Where("fid = ?", fid).Find(&meta).Error; err != nil {
		logrus.Errorf("Get Error: %s", err)
	}
	return
}

func (h *FileInfoMySql) MGet(fids []int64) (metas []*model.File, err error) {
	if len(fids) == 0 {
		return nil, errors.New("len(fids) == 0")
	}
	if err = h.conn.Where("fid IN (?)", fids).Find(&metas).Error; err != nil {
		logrus.Errorf("MGet Error: %vs", err)
	}
	return
}

func (h *FileInfoMySql) Begin() *gorm.DB {
	return h.conn.Begin()
}

func (h *FileInfoMySql) UpdateMd5(fid int64, md5 string) error {
	if err := h.conn.Where("fid = ?", fid).Update(map[string]string{
		"md5": md5,
	}).Error; err != nil {
		logrus.Errorf("UpdateMd5 fid: %d Error: %s", fid, err)
		return err
	}
	return nil
}

func (h *FileInfoMySql) GetFileByUid(uid, offset, limit int64) (files []*model.File, err error) {
	if err = h.conn.Where("uid = ?", uid).Offset(offset).Limit(limit).Find(&files).Error; err != nil {
		logrus.Errorf("GetFileByUid Error: %s", err)
	}
	return
}
