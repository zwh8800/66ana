package model

import "github.com/jinzhu/gorm"

const DyGiftTypeTableName = "dy_gift_type"

type DyGift struct {
	gorm.Model

	RoomId       int64
	Gid          int64
	name         string
	GiftType     DyGiftType
	Price        float64
	contribution int
	intro        string
	description  string
	himg         string
	mimg         string
}

func (*DyGift) TableName() string {
	return DyGiftTypeTableName
}
