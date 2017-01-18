package model

import "github.com/jinzhu/gorm"

const DyDanmuTableName = "dy_danmu"

type DyDanmu struct {
	gorm.Model

	Cid     string
	UserId  int64
	RoomId  int64
	Content string
	Color   DyDanmuColor
	Client  DyClientType
}

func (*DyDanmu) TableName() string {
	return DyDanmuTableName
}
