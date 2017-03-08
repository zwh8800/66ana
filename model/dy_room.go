package model

import "time"

const DyRoomTableName = "dy_room"

type DyRoom struct {
	Model

	Rid          int64        `json:"rid"`
	CateId       int64        `json:"cateId"`
	Name         string       `json:"name"`
	Status       DyRoomStatus `json:"status"`
	Thumb        string       `json:"thumb"`
	Avatar       string       `json:"avatar"`
	FansCount    int          `json:"fansCount"`
	OnlineCount  int          `json:"onlineCount"`
	OwnerName    string       `json:"ownerName"`
	Weight       int          `json:"weight"`
	LastLiveTime time.Time    `json:"lastLiveTime"`
}

func (*DyRoom) TableName() string {
	return DyRoomTableName
}

func (a DyRoom) Equals(b DyRoom) bool {
	// omit some field
	a.Model = b.Model

	return a == b
}

func (a *DyRoom) Assign(b *DyRoom) {
	// omit some field
	m := a.Model

	*a = *b
	a.Model = m
}
