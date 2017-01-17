package model

import "github.com/jinzhu/gorm"

const DyDeserveTableName = "dy_deserve"

type DyDeserve struct {
	gorm.Model

	UserId int64
	RoomId int64
	level  int
	Count  int
	Hits   int
}

func (*DyDeserve) TableName() string {
	return DyDeserveTableName
}
