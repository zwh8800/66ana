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
	cateInfos := getCateList()
	for _, cateInfo := range cateInfos {
		if _, err := service.InsertDyCate(cateInfo); err != nil {
			log.Println("service.InsertDyCate:", err)
		}
	}

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

	go removeExpireWorkingRoom()
}

var (
	dispatchLock sync.Mutex
)

func dispatchTask(report *model.ReportPayload) {
	dispatchLock.Lock()
	defer dispatchLock.Unlock()

	for _, room := range report.Workers {
		if _, err := service.AddWorkingRoom(room.RoomId); err != nil {
			log.Println("service.AddWorkingRoom(", room.RoomId, "):", err)
		}
	}

	if report.Working < report.Capacity {
		n := report.Capacity - report.Working
		toStartList := make([]int64, 0, n)

	out:
		for i := 0; n > 0; i++ {
			list := getLiveList(i)
			for _, roomId := range list {
				exist, err := service.AddWorkingRoom(roomId)
				if err != nil {
					log.Println("service.AddWorkingRoom(", roomId, "):", err)
					continue
				}
				if exist {
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
}

func removeExpireWorkingRoom() {
	for {
		if err := service.RemoveExpireWorkingRoom(); err != nil {
			log.Println("service.RemoveExpireWorkingRoom():", err)
		}

		time.Sleep(1 * time.Second)
	}
}

// FIXME: panic
func getLiveList(page int) []int64 {
	resp, err := http.Get(fmt.Sprintf("http://api.douyutv.com/api/v1/live?offset=%d&limit=30", page*30))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var live model.LiveInfoJson
	if err := json.Unmarshal(data, &live); err != nil {
		panic(err)
	}

	if live.Error != 0 {
		panic(err)
	}

	roomIdList := make([]int64, 0, len(live.Data))
	for _, l := range live.Data {
		rid, err := strconv.ParseInt(l.RoomID, 10, 64)
		if err != nil {
			continue
		}
		roomIdList = append(roomIdList, rid)
	}
	return roomIdList
}

// FIXME: panic
func getCateList() []*model.CateInfo {
	resp, err := http.Get("http://open.douyucdn.cn/api/RoomApi/game")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var cate model.CateInfoJson
	if err := json.Unmarshal(data, &cate); err != nil {
		panic(err)
	}

	if cate.Error != 0 {
		panic(err)
	}

	return cate.Data
}
