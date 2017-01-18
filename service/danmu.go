package service

import (
	"strconv"

	"github.com/zwh8800/66ana/model"
)

func InsertDyDanmu(message map[string]string) (*model.DyDanmu, error) {
	committed := false
	tx := dbConn.Begin()
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()
	room, user, danmu, err := cookModelFromDanmu(message)
	if err != nil {
		return nil, err
	}

	if err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where(room).FirstOrCreate(room).Error; err != nil {
		return nil, err
	}

	user.FirstAppearedRoomId = int64(room.ID)
	user.LastAppearedRoomId = int64(room.ID)
	assignUser := *user
	assignUser.FirstAppearedRoomId = 0

	if err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where(model.DyUser{Uid: user.Uid}).
		Attrs(user).FirstOrCreate(user).Error; err != nil {
		return nil, err
	}
	if !user.Equals(&assignUser) {
		if err := tx.Model(user).Omit("first_appeared_room_id").
			Update(assignUser).Error; err != nil {
			return nil, err
		}
	}

	danmu.RoomId = int64(room.ID)
	danmu.UserId = int64(user.ID)
	if err := tx.Create(danmu).Error; err != nil {
		return nil, err
	}

	tx.Commit()
	committed = true
	return danmu, nil
}

func cookModelFromDanmu(message map[string]string) (*model.DyRoom, *model.DyUser, *model.DyDanmu, error) {
	rid, err := strconv.ParseInt(message["rid"], 10, 64)
	if err != nil {
		return nil, nil, nil, err
	}
	room := &model.DyRoom{
		Rid: rid,
	}

	uid, err := strconv.ParseInt(message["uid"], 10, 64)
	if err != nil {
		return nil, nil, nil, err
	}
	level, _ := strconv.ParseInt(message["level"], 10, 64)
	giftRank, _ := strconv.ParseInt(message["gt"], 10, 64)
	pg, err := strconv.ParseInt(message["pg"], 10, 64)
	if err != nil {
		pg = 1
	}
	dlv, _ := strconv.ParseInt(message["dlv"], 10, 64)
	dc, _ := strconv.ParseInt(message["dc"], 10, 64)
	bdlv, _ := strconv.ParseInt(message["bdlv"], 10, 64)
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

	color, _ := strconv.ParseInt(message["col"], 10, 64)
	client, _ := strconv.ParseInt(message["ct"], 10, 64)
	danmu := &model.DyDanmu{
		Cid:     message["cid"],
		Color:   model.DyDanmuColor(color),
		Client:  model.DyClientType(client),
		Content: message["txt"],
	}

	return room, user, danmu, nil
}
