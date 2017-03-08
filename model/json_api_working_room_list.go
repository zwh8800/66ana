package model

type ApiWorkingRoomList struct {
	Total           int64     `json:"total"`
	WorkingRoomList []*DyRoom `json:"workingRoomList"`
}
