package model

import "time"

const DyGiftHistoryTableName = "dy_gift_history"

type DyGiftHistory struct {
	Model

	UserId    int64
	RoomId    int64
	GiftId    int64
	Count     int
	Hits      int
	GiftStyle string
}

func (m *DyGiftHistory) TableName() string {
	return DyGiftHistoryTableName + "_" + m.CreatedAt.Format("20060102")
}

type DyGiftHistoryWithDay struct {
	DyGiftHistory

	day time.Time `gorm:"-"`
}

func NewDyGiftHistoryWithDay(day time.Time) *DyGiftHistoryWithDay {
	return &DyGiftHistoryWithDay{
		day: day,
	}
}

func (m *DyGiftHistoryWithDay) TableName() string {
	return DyGiftHistoryTableName + "_" + m.day.Format("20060102")
}
