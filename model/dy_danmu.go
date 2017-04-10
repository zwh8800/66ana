package model

import "time"

const DyDanmuTableName = "dy_danmu"

// DyDanmu: only manipulate current day
type DyDanmu struct {
	Model

	Cid     string
	UserId  int64
	RoomId  int64
	Content string
	Color   DyDanmuColor
	Client  DyClientType
}

func (m *DyDanmu) TableName() string {
	return DyDanmuTableName + "_" + m.CreatedAt.Format("20060102")
}

type DyDanmuWithDay struct {
	DyDanmu

	day time.Time `gorm:"-"`
}

func NewDyDanmuWithDay(day time.Time) *DyDanmuWithDay {
	return &DyDanmuWithDay{
		day: day,
	}
}

func (m *DyDanmuWithDay) TableName() string {
	return DyDanmuTableName + "_" + m.day.Format("20060102")
}
