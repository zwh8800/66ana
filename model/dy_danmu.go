package model

const DyDanmuTableName = "dy_danmu"

type DyDanmu struct {
	Model

	Cid     string
	UserId  int64
	RoomId  int64
	Content string
	Color   DyDanmuColor
	Client  DyClientType
}

func (*DyDanmu) TableName() string {
	return DyDanmuTableName
}
