package model

import "github.com/jinzhu/gorm"

const DyUserRoomTableName = "dy_user_room"

type DyUserRoom struct {
	gorm.Model

	UserId        int64
	RoomId        int64
	RoomPrivilege DyPrivilege
}

func (*DyUserRoom) TableName() string {
	return DyUserRoomTableName
}
