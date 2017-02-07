package service

import (
	"fmt"
	"strconv"
	"time"

	"github.com/zwh8800/66ana/model"
)

func InsertDyRoom(roomInfo *model.RoomInfo) (*model.DyRoom, error) {
	committed := false
	tx := dbConn.Begin()
	defer func() {
		if !committed {
			tx.Rollback()
		}
	}()

	cate, room, giftList, err := cookModelFromRoomInfo(roomInfo)
	if err != nil {
		return nil, err
	}
	if err := tx.Where(model.DyCate{Cid: cate.Cid}).
		Attrs(cate).FirstOrCreate(cate).Error; err != nil {
		return nil, err
	}

	room.CateId = int64(cate.ID)
	updatedRoom := *room
	if err := tx.Where(model.DyRoom{Rid: room.Rid}).
		Attrs(room).FirstOrCreate(room).Error; err != nil {
		return nil, err
	}
	if !room.Equals(updatedRoom) {
		if err := tx.Model(room).Update(updatedRoom).
			Error; err != nil {
			return nil, err
		}
	}

	for _, gift := range giftList {
		gift.RoomId = int64(room.ID)
		updatedGift := *gift
		if err := tx.Where(model.DyGift{RoomId: gift.RoomId, Gid: gift.Gid}).
			Attrs(gift).FirstOrCreate(gift).Error; err != nil {
			return nil, err
		}
		if !gift.Equals(updatedGift) {
			if err := tx.Model(gift).Update(updatedGift).
				Error; err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	committed = true
	return room, nil
}

func cookModelFromRoomInfo(roomInfo *model.RoomInfo) (*model.DyCate, *model.DyRoom, []*model.DyGift, error) {
	cid, err := strconv.ParseInt(roomInfo.CateID, 10, 64)
	if err != nil {
		return nil, nil, nil, err
	}
	cate := &model.DyCate{
		Cid:      cid,
		GameName: roomInfo.CateName,
	}

	rid, err := strconv.ParseInt(roomInfo.RoomID, 10, 64)
	if err != nil {
		return nil, nil, nil, err
	}
	status, err := strconv.ParseInt(roomInfo.RoomStatus, 10, 64)
	if err != nil {
		return nil, nil, nil, err
	}
	fansCount, err := strconv.ParseInt(roomInfo.FansNum, 10, 64)
	if err != nil {
		return nil, nil, nil, err
	}
	weight, err := parseWeight(roomInfo.OwnerWeight)
	if err != nil {
		return nil, nil, nil, err
	}
	lastLiveTime, err := time.ParseInLocation("2006-01-02 15:04", roomInfo.StartTime, time.Local)
	if err != nil {
		return nil, nil, nil, err
	}
	room := &model.DyRoom{
		Rid:          rid,
		Name:         roomInfo.RoomName,
		Status:       model.DyRoomStatus(status),
		Thumb:        roomInfo.RoomThumb,
		Avatar:       roomInfo.Avatar,
		FansCount:    int(fansCount),
		OnlineCount:  roomInfo.Online,
		OwnerName:    roomInfo.OwnerName,
		Weight:       weight,
		LastLiveTime: lastLiveTime,
	}

	giftList := make([]*model.DyGift, 0, len(roomInfo.Gift))
	for _, giftInfo := range roomInfo.Gift {
		gid, err := strconv.ParseInt(giftInfo.ID, 10, 64)
		if err != nil {
			return nil, nil, nil, err
		}
		giftType, err := strconv.ParseInt(giftInfo.Type, 10, 64)
		if err != nil {
			return nil, nil, nil, err
		}
		gift := &model.DyGift{
			Gid:          gid,
			Name:         giftInfo.Name,
			GiftType:     model.DyGiftType(giftType),
			Price:        giftInfo.Pc,
			Contribution: giftInfo.Gx,
			Intro:        giftInfo.Intro,
			Mimg:         giftInfo.Mimg,
			Himg:         giftInfo.Himg,
		}
		giftList = append(giftList, gift)
	}
	return cate, room, giftList, nil
}

func parseWeight(weightStr string) (int, error) {
	unit := 1
	l := len(weightStr)
	if weightStr[l-2:] == "kg" {
		unit = 1000
		weightStr = weightStr[:l-2]
	} else if weightStr[l-1:] == "g" {
		unit = 1
		weightStr = weightStr[:l-1]
	} else if weightStr[l-1:] == "t" {
		unit = 1000 * 1000
		weightStr = weightStr[:l-1]
	} else {
		return 0, fmt.Errorf("unexpected weight string: %s", weightStr)
	}

	num, err := strconv.ParseFloat(weightStr, 64)
	if err != nil {
		return 0, err
	}
	return int(num * float64(unit)), nil
}

func FindRoomByRidList(ridList []int64) ([]*model.DyRoom, error) {
	var roomList []*model.DyRoom
	if err := dbConn.Where("rid in (?)", ridList).Order("online_count desc").Find(&roomList).Error; err != nil {
		return nil, err
	}
	return roomList, nil
}
