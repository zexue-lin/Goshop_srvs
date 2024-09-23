package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseModel struct {
	ID        int32     `gorm:"primary_key;AUTO_INCREMENT"`
	CreatedAt time.Time `gorm:"column:add_time"`
	UpdatedAt time.Time `gorm:"column:update_time"`
	DeletedAt gorm.DeletedAt
	IsDeleted bool
}

type User struct {
	BaseModel
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null;comment:手机号码"`
	Password string     `gorm:"type:varchar(100);not null;comment:密码"`
	NickName string     `gorm:"type:varchar(20);comment:昵称"`
	Birthday *time.Time `gorm:"type:datetime;comment:生日"`
	Gender   string     `gorm:"column:gender;default:male;type:varchar(6);comment:female表示女,male表示男"`
	Role     int        `gorm:"column:role;default:1;type:int;comment:1表示普通用户，2表示管理员"`
}
