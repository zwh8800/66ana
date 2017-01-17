package model

import "github.com/jinzhu/gorm"

const DySuperDanmuTableName = "dy_super_danmu"

type DySuperDanmu struct {
	gorm.Model

	Sdid       int64
	RoomId     int64
	JumpRoomId int64
	Content    string
}
