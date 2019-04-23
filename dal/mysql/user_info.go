package mysql

import (
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

type UserInfoMySql struct {
	Conn *gorm.DB
}

func NewUserInfoMySql() *UserInfoMySql {
	var err error
	h := &UserInfoMySql{}
	driver, url := getUserInfoMySqlConfig()
	logrus.Infof("driver: %s url: %s", driver, url)
	h.Conn, err = gorm.Open(driver, url)
	if err != nil {
		logrus.Panicf("Create UserInfo MySQL connection Error: %s", err)
		panic(err)
	}
	// In test env, print SQL gorm execute
	if !tools.IsProductEnv() {
		h.Conn = h.Conn.Debug()
	}
	return h
}

func getUserInfoMySqlConfig() (string, string) {
	return "mysql", constdef.MySqlUrl
}

func (h *UserInfoMySql) Save(conn *gorm.DB, user *model.User) (err error) {
	if conn == nil {
		conn = h.Conn
	}

	if err = conn.Save(user).Error; err != nil {
		logrus.Errorf("Save user %+v info Error: %s", user, err)
	}
	return
}

func (h *UserInfoMySql) Get(uid int64) (user *model.User, err error) {
	user = new(model.User)
	if err = h.Conn.Where("uid IN (?)", uid).Find(user).Error; err != nil {
		logrus.Errorf("Get uid: %d Error: %s", uid, err)
	}
	return
}
