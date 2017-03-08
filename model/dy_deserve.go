package model

const DyDeserveTableName = "dy_deserve"

type DyDeserve struct {
	Model

	UserId int64
	RoomId int64
	level  int
	Count  int
	Hits   int
}

func (*DyDeserve) TableName() string {
	return DyDeserveTableName
}
