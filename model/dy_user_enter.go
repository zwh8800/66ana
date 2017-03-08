package model

const DyUserEnterTableName = "dy_user_enter"

type DyUserEnter struct {
	Model

	UserId int64
	RoomId int64
}

func (*DyUserEnter) TableName() string {
	return DyUserEnterTableName
}
