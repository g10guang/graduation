package model

type User struct {
	Uid    int64  `json:"uid";gorm:"uid"`
	Status int8   `json:"status";gorm:"status"`
	Extra  string `json:"extra";gorm:"extra"`
}

func (User) TableName() string {
	return "user"
}
