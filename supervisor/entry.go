package supervisor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/zwh8800/66ana/model"
	"github.com/zwh8800/66ana/service"
	"github.com/zwh8800/66ana/util"
)

func Run() {
	updateCateInfo()
	removeExpireWorkingRoom()

	service.SubscribeReport(func(payload *model.ReportPayload, err error) {
		log.Println("Report:", util.JsonStringify(payload, false), "err:", err)
		if err != nil {
			return
		}
		dispatchTask(payload)
	})

	service.SubscribeSpiderClosed(func(payload *model.SpiderClosedPayload, err error) {
		log.Println("SubscribeSpiderClosed:", util.JsonStringify(payload, false), "err:", err)
		if err != nil {
			return
		}

		if err := service.RemoveWorkingRoom(payload.RoomId); err != nil {
			log.Println("service.RemoveWorkingRoom(", payload.RoomId, "):", err)
		}
		dispatchTask(payload.ReportPayload)
	})

	go func() {
		for {
			removeExpireWorkingRoom()
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for {
			updateCateInfo()
			time.Sleep(1 * time.Minute)
		}
	}()
}

// FIXME: 用更细粒度的锁，或者放redis里解决
var (
	dispatchLock sync.Mutex
)

func dispatchTask(report *model.ReportPayload) {
	dispatchLock.Lock()
	defer dispatchLock.Unlock()
	if err := service.AddWorker(report.BasicWorkerInfo); err != nil {
		log.Println("service.AddWorker(", report.WorkerId, "):", err)
		return
	}

	for _, room := range report.Workers {
		if err := service.RemoveFromWorkingRoomQueue(report.WorkerId, room.RoomId); err != nil {
			log.Println("service.RemoveFromWorkingRoomQueue(", report.WorkerId, room.RoomId, "):", err)
			return
		}
		if err := service.AddWorkingRoom(room.RoomId); err != nil {
			log.Println("service.AddWorkingRoom(", room.RoomId, "):", err)
			return
		}
	}

	queueLen, err := service.CountWorkingRoomQueue(report.WorkerId)
	if err != nil {
		log.Println("service.CountWorkingRoomQueue", report.WorkerId, "):", err)
	}

	n := int64(report.Capacity) - (int64(report.Working) + queueLen)
	toStartList := make([]int64, 0)
out:
	for i := 0; n > 0; i++ {
		list, err := getLiveList(i)
		if err != nil {
			log.Println("getLiveList(", i, "):", err)
			i--
			time.Sleep(3 * time.Second)
			continue
		}
		for _, roomId := range list {
			exist, err := service.IsWorkingRoom(roomId)
			if err != nil {
				log.Println("service.IsWorkingRoom(", roomId, "):", err)
				continue
			}
			if exist {
				continue
			}
			exist, err = service.IsInWorkingRoomQueue(roomId)
			if err != nil {
				log.Println("service.IsInWorkingRoomQueue(", roomId, "):", err)
				continue
			}
			if exist {
				continue
			}

			if err := service.AddToWorkingRoomQueue(report.WorkerId, roomId); err != nil {
				log.Println("service.AddToWorkingRoomQueue(", report.WorkerId, roomId, "):", err)
				continue
			}

			toStartList = append(toStartList, roomId)

			n--
			if n <= 0 {
				break out
			}
		}
	}

	for _, roomId := range toStartList {
		if err := service.PublishStartSpider(report.WorkerId, &model.StartSpiderPayload{
			RoomId: roomId,
		}); err != nil {
			log.Println("service.InsertDyCate:", err)
		}
	}
}

func removeExpireWorkingRoom() {
	if err := service.RemoveExpireWorker(); err != nil {
		log.Println("service.RemoveExpireWorker():", err)
		return
	}
	if err := service.RemoveExpireWorkingRoom(); err != nil {
		log.Println("service.RemoveExpireWorkingRoom():", err)
		return
	}
}

func updateCateInfo() {
	cateInfos, err := getCateList()
	if err != nil {
		log.Println("getCateList():", err)
		return
	}
	for _, cateInfo := range cateInfos {
		if _, err := service.InsertDyCate(cateInfo); err != nil {
			log.Println("service.InsertDyCate:", err)
		}
	}
}

func getLiveList(page int) ([]int64, error) {
	resp, err := http.Get(fmt.Sprintf("http://api.douyutv.com/api/v1/live?offset=%d&limit=30", page*30))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var live model.LiveInfoJson
	if err := json.Unmarshal(data, &live); err != nil {
		return nil, err
	}

	if live.Error != 0 {
		return nil, err
	}

	roomIdList := make([]int64, 0, len(live.Data))
	for _, l := range live.Data {
		rid, err := strconv.ParseInt(l.RoomID, 10, 64)
		if err != nil {
			continue
		}
		roomIdList = append(roomIdList, rid)
	}
	return roomIdList, nil
}

func getCateList() ([]*model.CateInfo, error) {
	resp, err := http.Get("http://open.douyucdn.cn/api/RoomApi/game")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var cate model.CateInfoJson
	if err := json.Unmarshal(data, &cate); err != nil {
		return nil, err
	}

	if cate.Error != 0 {
		return nil, err
	}

	return cate.Data, nil
}
