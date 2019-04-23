package model

import "time"

type File struct {
	Fid        int64     `json:"fid";gorm:"fid"`
	Uid        int64     `json:"uid";gorm:"uid"`
	Name       string    `json:"name";gorm:"name"`
	Size       int64     `json:"size";gorm:"fid"`
	Md5        string    `json:"md5";gorm:"md5"`
	CreateTime time.Time `json:"create_time";gorm:"create_time"`
	UpdateTime time.Time `json:"update_time";gorm:"update_time"`
	Extra      string    `json:"extra";gorm:"extra"`
	Status     int8      `json:"status";gorm:"status"`
}

func (File) TableName() string {
	return "file"
}
