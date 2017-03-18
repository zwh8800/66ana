package supervisor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/robfig/cron"
	"github.com/zwh8800/66ana/conf"
	"github.com/zwh8800/66ana/model"
	"github.com/zwh8800/66ana/service"
	"github.com/zwh8800/66ana/util"
)

func Run() {
	if err := service.CreateFurtherTable(5); err != nil {
		panic(err)
	}

	setupCron()

	updateCateInfo()
	removeExpire()

	removeClosedSuccessChan := make(chan bool, 50)

	go func() {
		for {
			select {
			case <-removeClosedSuccessChan:
				log.Println("dispatchLoop by removeClosedSuccessChan")
			case <-time.After(10 * time.Second):
				log.Println("dispatchLoop by ticker")
			}
			dispatchLoop()
		}
	}()

	go func() {
		for {
			removeClosed()
			removeClosedSuccessChan <- true
		}
	}()

	go func() {
		for {
			removeExpire()
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

func setupCron() {
	crontab := cron.New()
	crontab.AddFunc("0 50 23 * * *", func() {
		log.Println("push createFurtherTable job")
		service.SupervisorPushJob(&model.JobPayload{
			JobName: "createFurtherTable",
		})
	})

	crontab.Start()
}

func dispatchLoop() {
	workerCount, err := service.CountWorkers()
	if err != nil {
		log.Println("service.CountWorkers()", err)
		return
	}

	countMap := make(map[string]int)
	for i := int64(0); workerCount == 0 || i < workerCount; i++ {
		report, err := service.GetWorkerReport()
		if err != nil {
			log.Println("service.GetWorkerReport():", err)
			continue
		}
		if err := service.AddWorker(report.BasicWorkerInfo); err != nil {
			log.Println("service.AddWorker(", util.JsonStringify(report.BasicWorkerInfo, false), "):", err)
			continue
		}
		for _, room := range report.Workers {
			if err := service.AddWorkingRoom(room.RoomId); err != nil {
				log.Println("service.AddWorkingRoom(", room.RoomId, "):", err)
				return
			}
		}
		n := report.Capacity - report.Working
		countMap[report.WorkerId] = n

		thisWorkerCount, err := service.CountWorkers()
		if err != nil {
			log.Println("service.CountWorkers()", err)
			continue
		}
		workerCount = thisWorkerCount
	}

	totalDispatchCount := 0
	for _, count := range countMap {
		totalDispatchCount += count
	}

	roomIdToStart := make(map[int64]bool, totalDispatchCount)
out:
	for i := 0; len(roomIdToStart) < totalDispatchCount; i++ {
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

			roomIdToStart[roomId] = true
			if len(roomIdToStart) >= totalDispatchCount {
				break out
			}
		}
	}

	for roomId := range roomIdToStart {
		if service.DispatchWork(&model.StartSpiderPayload{
			RoomId: roomId,
		}); err != nil {
			log.Println("service.DispatchWork:", err)
		}
	}
}

func removeClosed() {
	payload, err := service.PullSpiderClosed()
	if err != nil {
		log.Println("service.RemoveExpireWorker():", err)
		return
	}
	if err := service.AddWorker(payload.BasicWorkerInfo); err != nil {
		log.Println("service.AddWorker(", util.JsonStringify(payload.BasicWorkerInfo, false), "):", err)
		return
	}
	for _, room := range payload.Workers {
		if err := service.AddWorkingRoom(room.RoomId); err != nil {
			log.Println("service.AddWorkingRoom(", room.RoomId, "):", err)
			return
		}
	}
	if err := service.RemoveWorkingRoom(payload.RoomId); err != nil {
		log.Println("service.RemoveWorkingRoom(", payload.RoomId, "):", err)
		return
	}
}

func removeExpire() {
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
		if l.Online < conf.Conf.Spider.MinOnlineCount {
			continue
		}
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
