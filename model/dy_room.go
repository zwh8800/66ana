package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

const DyRoomTableName = "dy_room"

type DyRoom struct {
	gorm.Model

	Rid          int64
	CateId       int64
	Name         string
	Status       DyRoomStatus
	Thumb        string
	Avatar       string
	FansCount    int
	OwnerName    string
	Weight       int
	LastLiveTime time.Time
}

func (*DyRoom) TableName() string {
	return DyRoomTableName
}
