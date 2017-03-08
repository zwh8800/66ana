package model

const DyGiftTableName = "dy_gift"

type DyGift struct {
	Model

	RoomId       int64
	Gid          int64
	Name         string
	GiftType     DyGiftType
	Price        float64
	Contribution int
	Intro        string
	Description  string
	Himg         string
	Mimg         string
}

func (*DyGift) TableName() string {
	return DyGiftTableName
}

func (a DyGift) Equals(b DyGift) bool {
	// omit some field
	a.Model = b.Model

	return a == b
}

func (a *DyGift) Assign(b *DyGift) {
	// omit some field
	m := a.Model

	*a = *b
	a.Model = m
}
