package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/zwh8800/66ana/model"
)

func InsertOrUpdateDyUser(db *gorm.DB, user *model.DyUser) error {
	updatedUser := *user
	updatedUser.FirstAppearedRoomId = 0
	if err := db.Set("gorm:query_option", "FOR UPDATE").
		Where(model.DyUser{Uid: user.Uid}).
		Attrs(user).FirstOrCreate(user).Error; err != nil {
		return err
	}
	if !user.Equals(updatedUser) {
		user.Assign(&updatedUser)
		if err := db.Save(user).
			Error; err != nil {
			return err
		}
	}
	return nil
}
