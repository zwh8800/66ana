package model

type DyRoomStatus int

func (s DyRoomStatus) IsOnline() bool {
	return s == 1
}

func (s DyRoomStatus) IsOffline() bool {
	return s == 2
}
