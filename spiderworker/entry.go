package spiderworker

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/satori/go.uuid"
	"github.com/zwh8800/66ana/conf"
	"github.com/zwh8800/66ana/service"
	"github.com/zwh8800/66ana/util"
)

var (
	workerId   string
	workers    map[int64]*worker
	workerLock sync.RWMutex
	closeChan  chan int64
)

func init() {
	workerId = uuid.NewV4().String()
	workers = make(map[int64]*worker, conf.Conf.SpiderWorker.Capacity)
	closeChan = make(chan int64, conf.Conf.SpiderWorker.Capacity/10)
}

func Run() {
	service.SubscribeStartSpider(workerId, func(payload *service.StartSpiderPayload, err error) {
		if err != nil {
			log.Println("SubscribeStartSpider:", err)
			return
		}
		newJob(payload.RoomId)
	})

	go report()
	go checkClosed()
}

func newJob(roomId int64) {
	workerLock.Lock()
	defer workerLock.Unlock()
	worker := newWorker(roomId, closeChan)
	workers[roomId] = worker
}

func report() {
	for {
		service.PublishReport(generateReport())
		time.Sleep(time.Duration(10+rand.Intn(4)-2) * time.Second) // +/- 2s
	}
}

func checkClosed() {
	for {
		roomId := <-closeChan
		service.PublishSpiderClosed(&service.SpiderClosedPayload{
			WorkerId:      workerId,
			RoomId:        roomId,
			ReportPayload: *generateReport(),
		})
	}
}

func getWorkingRoomIdList() []int64 {
	list := make([]int64, 0, len(workers))
	workerLock.RLock()
	defer workerLock.RUnlock()
	for _, worker := range workers {
		list = append(list, worker.GetRoomId())
	}
	return list
}

func generateReport() *service.ReportPayload {
	return &service.ReportPayload{
		WorkerId: workerId,
		Capacity: conf.Conf.SpiderWorker.Capacity,
		Working:  len(workers),

		RoomIdList: getWorkingRoomIdList(),

		CpuCount: util.CpuCount(),
		MemUsage: util.MemUsage(),
	}
}
