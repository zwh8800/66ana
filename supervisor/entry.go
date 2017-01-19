package supervisor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/zwh8800/66ana/model"
	"github.com/zwh8800/66ana/service"
	"github.com/zwh8800/66ana/util"
)

func Run() {
	started := false
	service.SubscribeReport(func(payload *service.ReportPayload, err error) {
		log.Println("Report:", util.JsonStringify(payload, false), "err:", err)
		if err != nil {
			return
		}
		if !started {
			cateInfos := getCateList()

			for _, cateInfo := range cateInfos {
				if _, err := service.InsertDyCate(cateInfo); err != nil {
					log.Println("service.InsertDyCate:", err)
				}
			}

			list := getLiveList()

			for _, roomId := range list {
				if err := service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
					RoomId: roomId,
				}); err != nil {
					log.Println("service.InsertDyCate:", err)
				}
			}
			started = true
		}
	})
}

// FIXME: panic
func getLiveList() []int64 {
	resp, err := http.Get("http://api.douyutv.com/api/v1/live")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var live model.LiveInfoJson
	if err := json.Unmarshal(data, &live); err != nil {
		return nil
	}

	if live.Error != 0 {
		return nil
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
		return nil
	}
	var cate model.CateInfoJson
	if err := json.Unmarshal(data, &cate); err != nil {
		return nil
	}

	if cate.Error != 0 {
		return nil
	}

	return cate.Data
}
