package supervisor

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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
			list := getLiveList()

			for _, roomId := range list {
				service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
					RoomId: roomId,
				})
			}
			started = true
		}
	})
}

// FIXME: panic
func getLiveList() []int64 {
	resp, err := http.Get("http://api.douyutv.com/api/v1/live/")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	var live liveData
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

type liveData struct {
	Error int `json:"error"`
	Data  []struct {
		RoomID          string `json:"room_id"`
		RoomSrc         string `json:"room_src"`
		VerticalSrc     string `json:"vertical_src"`
		IsVertical      int    `json:"isVertical"`
		CateID          string `json:"cate_id"`
		RoomName        string `json:"room_name"`
		ShowStatus      string `json:"show_status"`
		Subject         string `json:"subject"`
		ShowTime        string `json:"show_time"`
		OwnerUID        string `json:"owner_uid"`
		SpecificCatalog string `json:"specific_catalog"`
		SpecificStatus  string `json:"specific_status"`
		VodQuality      string `json:"vod_quality"`
		Nickname        string `json:"nickname"`
		Online          int    `json:"online"`
		URL             string `json:"url"`
		GameURL         string `json:"game_url"`
		GameName        string `json:"game_name"`
		ChildID         string `json:"child_id"`
		Avatar          string `json:"avatar"`
		AvatarMid       string `json:"avatar_mid"`
		AvatarSmall     string `json:"avatar_small"`
		JumpURL         string `json:"jumpUrl"`
		IsHide          int    `json:"isHide,omitempty"`
		Fans            string `json:"fans"`
		Ranktype        int    `json:"ranktype"`
		AnchorCity      string `json:"anchor_city"`
	} `json:"data"`
}
