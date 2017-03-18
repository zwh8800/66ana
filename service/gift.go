package service

import (
	"strconv"
	"time"

	"github.com/zwh8800/66ana/model"
)

func InsertDyGiftHistory(message map[string]string) (*model.DyGiftHistory, error) {
	committed := false
	tx := dbConn.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	room, user, gift, giftHistory, err := cookModelFromGift(message)
	if err != nil {
		return nil, err
	}

	if err := tx.Where(&model.DyRoom{Rid: room.Rid}).
		Attrs(room).FirstOrCreate(room).Error; err != nil {
		return nil, err
	}

	gift.RoomId = int64(room.ID)
	if err := tx.Where(gift).FirstOrCreate(gift).Error; err != nil {
		return nil, err
	}

	user.FirstAppearedRoomId = int64(room.ID)
	user.LastAppearedRoomId = int64(room.ID)
	updatedUser := *user
	updatedUser.FirstAppearedRoomId = 0
	if err := tx.Where(model.DyUser{Uid: user.Uid}).
		Attrs(user).FirstOrCreate(user).Error; err != nil {
		return nil, err
	}
	if !user.Equals(updatedUser) {
		user.Assign(&updatedUser)
		if err := tx.Save(user).
			Error; err != nil {
			return nil, err
		}
	}

	giftHistory.UserId = int64(user.ID)
	giftHistory.RoomId = int64(room.ID)
	giftHistory.GiftId = int64(gift.ID)
	createdAt := giftHistory.CreatedAt
	if err := tx.Create(giftHistory).Error; err != nil {
		return nil, err
	}
	tx.Model(giftHistory).Update("CreatedAt", createdAt)

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	committed = true
	return giftHistory, nil
}

func cookModelFromGift(message map[string]string) (*model.DyRoom, *model.DyUser, *model.DyGift, *model.DyGiftHistory, error) {
	rid, err := strconv.ParseInt(message["rid"], 10, 64)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	weight, _ := strconv.ParseInt(message["dw"], 10, 64)
	room := &model.DyRoom{
		Rid:    rid,
		Weight: int(weight),
	}

	uid, err := strconv.ParseInt(message["uid"], 10, 64)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	level, _ := strconv.ParseInt(message["level"], 10, 64)
	giftRank, _ := strconv.ParseInt(message["gt"], 10, 64)
	pg, err := strconv.ParseInt(message["pg"], 10, 64)
	if err != nil {
		pg = 1
	}
	dlv, _ := strconv.ParseInt(message["dlv"], 10, 64)
	dc, _ := strconv.ParseInt(message["dc"], 10, 64)
	bdlv, _ := strconv.ParseInt(message["bdl"], 10, 64)
	user := &model.DyUser{
		Uid:               uid,
		Nickname:          message["nn"],
		Level:             int(level),
		GiftRank:          int(giftRank),
		PlatformPrivilege: model.DyPrivilege(pg),
		DeserveLevel:      int(dlv),
		DeserveCount:      int(dc),
		BdeserveLevel:     int(bdlv),
	}

	gfid, err := strconv.ParseInt(message["gfid"], 10, 64)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	gift := &model.DyGift{
		Gid: gfid,
	}

	count, _ := strconv.ParseInt(message["gfcnt"], 10, 64)
	hits, _ := strconv.ParseInt(message["hits"], 10, 64)

	giftHistory := &model.DyGiftHistory{
		Count:     int(count),
		Hits:      int(hits),
		GiftStyle: message["gs"],
	}
	if timestamp, err := strconv.ParseInt(message["timestamp"], 10, 64); err == nil {
		giftHistory.CreatedAt = time.Unix(0, timestamp)
	}

	return room, user, gift, giftHistory, nil
}
