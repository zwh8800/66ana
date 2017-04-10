package model

import "time"

const DyDeserveTableName = "dy_deserve"

type DyDeserve struct {
	Model

	UserId int64
	RoomId int64
	level  int
	Count  int
	Hits   int
}

func (m *DyDeserve) TableName() string {
	return DyDeserveTableName + "_" + m.CreatedAt.Format("20060102")
}

type DyDeserveWithDay struct {
	DyDeserve

	day time.Time `gorm:"-"`
}

func NewDyDeserveWithDay(day time.Time) *DyDeserveWithDay {
	return &DyDeserveWithDay{
		day: day,
	}
}

func (m *DyDeserveWithDay) TableName() string {
	return DyDeserveTableName + "_" + m.day.Format("20060102")
}
