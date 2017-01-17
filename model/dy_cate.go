package model

import "github.com/jinzhu/gorm"

const DyCateTableName = "dy_cate"

type DyCate struct {
	gorm.Model

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
