package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/zwh8800/66ana/model"
)

func InsertDanmu(db *gorm.DB, danmu *model.DyDanmu) error {
	if err := db.Create(danmu).Error; err != nil {
		return err
	}
	return nil
}
