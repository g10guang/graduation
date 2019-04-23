package mysql

import (
	"github.com/g10guang/graduation/model"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestSaveFile(t *testing.T) {
	err := FileMySQL.Save(nil, &model.File{
		Fid: 1,
		Uid: 1,
		Name: "静态资源缓存层级.jpg",
		Size: 10,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Md5: "5d41402abc4b2a76b9719d911017c592",
		Extra: "hello",
		Status: 0,
	})
	if err != nil {
		logrus.Errorf("MySQL Error: %s", err)
	} else {
		logrus.Infof("Success")
	}
}

func TestGetFile(t *testing.T) {
	f, err := FileMySQL.Get(1)
	if err != nil {
		logrus.Errorf("MySQL Error: %s", err)
	} else {
		logrus.Infof("file: %+v", f)
	}
}

func TestSaveUser(t *testing.T) {
	err := UserMySQL.Save(nil, &model.User{
		Uid: 1,
	})
	if err != nil {
		logrus.Errorf("MySQL Error: %s", err)
	} else {
		logrus.Infof("Success")
	}
}

func TestGetUser(t *testing.T) {
	u, err := UserMySQL.Get(1)
	if err != nil {
		logrus.Errorf("MySQL Error: %s", err)
	} else {
		logrus.Infof("user: %+v", u)
	}
}

func TestDel(t *testing.T) {
	if err := FileMySQL.Delete(nil, []int64{0}, 0); err != nil {
		panic(err)
	}
}
