package model

import "github.com/jinzhu/gorm"

const DyGiftHistoryTableName = "dy_gift_history"

type DyGiftHistory struct {
	gorm.Model

	UserId    int64
	RoomId    int64
	GiftId    int64
	Count     int
	Hits      int
	GiftStyle string
}

func (*DyGiftHistory) TableName() string {
	return DyGiftHistoryTableName
}
