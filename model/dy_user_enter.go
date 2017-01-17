package model

import "github.com/jinzhu/gorm"

const DyUserEnterTableName = "dy_user_enter"

type DyUserEnter struct {
	gorm.Model

	UserId int64
	RoomId int64
}

func (*DyUserEnter) TableName() string {
	return DyUserEnterTableName
}
