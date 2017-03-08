package model

const DyCateTableName = "dy_cate"

type DyCate struct {
	Model

	Cid       int64
	GameName  string
	ShortName string
	GameUrl   string
	GameSrc   string
	GameIcon  string
}

func (*DyCate) TableName() string {
	return DyCateTableName
}

func (a DyCate) Equals(b DyCate) bool {
	// omit some field
	a.Model = b.Model

	return a == b
}

func (a *DyCate) Assign(b *DyCate) {
	// omit some field
	m := a.Model

	*a = *b
	a.Model = m
}
