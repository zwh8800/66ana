package model

const DySuperDanmuTableName = "dy_super_danmu"

type DySuperDanmu struct {
	Model

	Sdid       int64
	RoomId     int64
	JumpRoomId int64
	Content    string
}
