package model

type DyGiftType int

func (gt DyGiftType) IsYuwan() bool {
	return gt == 1
}

func (gt DyGiftType) IsYuchi() bool {
	return gt == 2
}
