package spiderworker

import (
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
	workers    map[string]*worker
	workerLock sync.RWMutex
)

func init() {
	workerId = uuid.NewV4().String()
	workers = make(map[string]*worker, conf.Conf.SpiderWorker.Capacity)
}

func Run() {
	service.SubscribeStartSpider(workerId, func(payload *service.StartSpiderPayload) {

	})

	go report()
}

func report() {
	for {
		service.PublishReport(generateReport())
		time.Sleep(10*time.Second + rand.Intn(4) - 2) // +/- 2s
	}
}

func getWorkingRoomIdList() []string {
	list := make([]string, 0, len(workers))
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
