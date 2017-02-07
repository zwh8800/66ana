package service

import (
	"strconv"
	"time"

	"gopkg.in/redis.v5"
)

const workingRoomKey = "working-room"
const workingRoomTTL = 60 * time.Second

func AddWorkingRoom(rid int64) (bool, error) {
	score := workingRoomTTL.Seconds() + float64(time.Now().Unix())
	exist, err := redisClient.ZAdd(workingRoomKey, redis.Z{Score: score, Member: rid}).Result()
	if err != nil {
		return false, err
	}
	return exist == 0, nil
}

func ListWorkingRoom() ([]int64, error) {
	roomStrList, err := redisClient.ZRange(workingRoomKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	roomIdList := make([]int64, 0, len(roomStrList))
	for _, roomStr := range roomStrList {
		roomId, err := strconv.ParseInt(roomStr, 10, 64)
		if err != nil {
			return nil, err
		}
		roomIdList = append(roomIdList, roomId)
	}
	return roomIdList, nil
}

func RemoveWorkingRoom(rid int64) error {
	return redisClient.ZRem(workingRoomKey, rid).Err()
}

func RemoveExpireWorkingRoom() error {
	return redisClient.ZRemRangeByScore(workingRoomKey, "0", strconv.FormatInt(time.Now().Unix(), 10)).Err()
}

//func IsWorkingRoom(rid int64) (bool, error) {
//	return redisClient.ZScore(workingRoomKey, rid).Result()
//}
