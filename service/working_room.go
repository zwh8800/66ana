package service

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/zwh8800/66ana/model"
	"gopkg.in/redis.v5"
)

const (
	workersKey             = "workers:"
	workingRoomKey         = "working-room:"
	workingRoomQueueKey    = "working-room-queue:"
	workingRoomQueueAllKey = "working-room-queue:all"
	workingRoomTTL         = 30 * time.Second
)

func AddWorker(workerInfo *model.BasicWorkerInfo) error {
	if workerInfo.WorkerId == "" {
		return errors.New("workerInfo.WorkerId == \"\"")
	}

	score := workingRoomTTL.Seconds() + float64(time.Now().Unix())
	if err := redisClient.ZAdd(workersKey, redis.Z{Score: score, Member: workerInfo.WorkerId}).Err(); err != nil {
		return err
	}
	workerInfoJson, err := json.Marshal(workerInfo)
	if err != nil {
		return err
	}
	return redisClient.Set(workersKey+workerInfo.WorkerId, workerInfoJson, workingRoomTTL).Err()
}

func ListWorkerIdList() ([]string, error) {
	return redisClient.ZRange(workersKey, 0, -1).Result()
}

func GetWorkerInfo(workerId string) (string, error) {
	return redisClient.Get(workersKey + workerId).Result()
}

func RemoveExpireWorker() error {
	return redisClient.ZRemRangeByScore(workersKey, "0", strconv.FormatInt(time.Now().Unix(), 10)).Err()
}

func AddToWorkingRoomQueue(workerId string, rid int64) error {
	key := workingRoomQueueKey + workerId
	if err := redisClient.SAdd(key, rid).Err(); err != nil {
		return err
	}
	return redisClient.Expire(key, workingRoomTTL).Err()
}

func CountWorkingRoomQueue(workerId string) (int64, error) {
	return redisClient.SCard(workingRoomQueueKey + workerId).Result()
}

func ListWorkingRoomQueue(workerId string) ([]int64, error) {
	roomStrList, err := redisClient.SMembers(workingRoomQueueKey + workerId).Result()
	if err != nil {
		return nil, err
	}
	return stringToInt64List(roomStrList)
}

func RemoveFromWorkingRoomQueue(workerId string, rid int64) error {
	key := workingRoomQueueKey + workerId

	if err := redisClient.SRem(key, rid).Err(); err != nil {
		return err
	}
	return redisClient.Expire(key, workingRoomTTL).Err()
}

func IsInWorkingRoomQueue(rid int64) (bool, error) {
	workerIdList, err := ListWorkerIdList()
	if err != nil {
		return false, err
	}
	keys := make([]string, 0, len(workersKey))
	for _, workerId := range workerIdList {
		key := workingRoomQueueKey + workerId
		if redisClient.Exists(key).Val() {
			keys = append(keys, key)
		}
	}
	if len(keys) == 0 {
		return false, nil
	}
	if err := redisClient.SUnionStore(workingRoomQueueAllKey, keys...).Err(); err != nil {
		return false, err
	}
	return redisClient.SIsMember(workingRoomQueueAllKey, rid).Result()
}

func AddWorkingRoom(rid int64) error {
	score := workingRoomTTL.Seconds() + float64(time.Now().Unix())
	return redisClient.ZAdd(workingRoomKey, redis.Z{Score: score, Member: rid}).Err()
}

func ListWorkingRoom() ([]int64, error) {
	roomStrList, err := redisClient.ZRange(workingRoomKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return stringToInt64List(roomStrList)
}

func RemoveWorkingRoom(rid int64) error {
	return redisClient.ZRem(workingRoomKey, rid).Err()
}

func RemoveExpireWorkingRoom() error {
	return redisClient.ZRemRangeByScore(workingRoomKey, "0", strconv.FormatInt(time.Now().Unix(), 10)).Err()
}

func IsWorkingRoom(rid int64) (bool, error) {
	_, err := redisClient.ZRank(workingRoomKey, strconv.FormatInt(rid, 10)).Result()
	if err == nil {
		return true, nil
	} else if err == redis.Nil {
		return false, nil
	} else {
		return false, err
	}
}

func stringToInt64List(strList []string) ([]int64, error) {
	intList := make([]int64, 0, len(strList))
	for _, str := range strList {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return nil, err
		}
		intList = append(intList, i)
	}
	return intList, nil
}
