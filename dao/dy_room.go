package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/zwh8800/66ana/model"
)

func InsertBasicRoom(db *gorm.DB, room *model.DyRoom) error {
	if err := db.Where(room).FirstOrCreate(room).Error; err != nil {
		return err
	}
	return nil
}

func InsertOrUpdateRoom(db *gorm.DB, room *model.DyRoom) error {
	updatedRoom := *room
	if err := db.Where(model.DyRoom{Rid: room.Rid}).
		Attrs(room).FirstOrCreate(room).Error; err != nil {
		return err
	}
	if !room.Equals(updatedRoom) {
		if err := db.Model(room).Update(updatedRoom).
			Error; err != nil {
			return err
		}
	}
	return nil
}
