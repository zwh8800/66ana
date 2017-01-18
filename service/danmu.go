package service

import (
	"log"
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

	tx.Where(room).FirstOrCreate(room)

	user.FirstAppearedRoomId = int64(room.ID)
	user.LastAppearedRoomId = int64(room.ID)
	assignUser := *user
	assignUser.FirstAppearedRoomId = 0

	tx.Set("gorm:query_option", "FOR UPDATE").
		Where(model.DyUser{Uid: user.Uid}).
		Attrs(user).
		Assign(assignUser).
		FirstOrCreate(user)

	danmu.RoomId = int64(room.ID)
	danmu.UserId = int64(user.ID)

	tx.Create(danmu)

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
	level, err := strconv.ParseInt(message["level"], 10, 64)
	if err != nil {
		log.Println("giftRank:", level)
	}
	giftRank, err := strconv.ParseInt(message["gt"], 10, 64)
	if err != nil {
		log.Println("giftRank:", err)
	}
	pg, err := strconv.ParseInt(message["pg"], 10, 64)
	if err != nil {
		pg = 1
		log.Println("PlatformPrivilege:", err)
	}
	dlv, err := strconv.ParseInt(message["dlv"], 10, 64)
	if err != nil {
		log.Println("DeserveLevel:", err)
	}
	dc, err := strconv.ParseInt(message["dc"], 10, 64)
	if err != nil {
		log.Println("DeserveCount:", err)
	}
	bdlv, err := strconv.ParseInt(message["bdlv"], 10, 64)
	if err != nil {
		log.Println("BdeserveLevel:", err)
	}
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

	color, err := strconv.ParseInt(message["color"], 10, 64)
	if err != nil {
		log.Println("color:", err)
	}
	client, err := strconv.ParseInt(message["client"], 10, 64)
	if err != nil {
		log.Println("client:", err)
	}
	danmu := &model.DyDanmu{
		Cid:     message["cid"],
		Color:   model.DyDanmuColor(color),
		Client:  model.DyClientType(client),
		Content: message["txt"],
	}

	return room, user, danmu, nil
}
