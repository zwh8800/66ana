package model

const DyUserRoomTableName = "dy_user_room"

type DyUserRoom struct {
	Model

	UserId        int64
	RoomId        int64
	RoomPrivilege DyPrivilege
}

func (*DyUserRoom) TableName() string {
	return DyUserRoomTableName
}
