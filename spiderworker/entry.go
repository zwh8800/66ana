package spiderworker

import (
	"log"
	"sync"

	"github.com/satori/go.uuid"
	"github.com/zwh8800/66ana/conf"
	"github.com/zwh8800/66ana/model"
	"github.com/zwh8800/66ana/service"
	"github.com/zwh8800/66ana/util"
)

var (
	workerId       string
	workers        map[int64]*worker
	workersLock    sync.RWMutex
	closeChan      chan int64
	pullNewJobSema util.Semaphore
)

func init() {
	workerId = uuid.NewV4().String()
	workers = make(map[int64]*worker, conf.Conf.SpiderWorker.Capacity)
	closeChan = make(chan int64, conf.Conf.SpiderWorker.Capacity/10)
	pullNewJobSema = util.NewSemaphore(conf.Conf.SpiderWorker.Capacity)
}

func Run() {
	service.ServeWorkerReport(func(err error) *model.ReportPayload {
		if err != nil {
			log.Println("ServeWorkerReport:", err)
			return nil
		}
		return generateReport()
	})

	go checkClosed()
	go pullNewJob()
}

func checkClosed() {
	for {
		roomId := <-closeChan
		pullNewJobSema.V(1)
		func() {
			workersLock.Lock()
			defer workersLock.Unlock()
			delete(workers, roomId) // gc will handle all
		}()

		if err := service.PushWorkerClosed(&model.SpiderClosedPayload{
			RoomId:        roomId,
			ReportPayload: generateReport(),
		}); err != nil {
			log.Println("service.PullWork():", err)
			continue
		}
	}
}

// 这个 goroutine 应该比上面那个慢
func pullNewJob() {
	for {
		pullNewJobSema.P(1)
		payload, err := service.PullWork()
		log.Println("service.PullWork()", util.JsonStringify(payload, false))
		if err != nil {
			log.Println("service.PullWork():", err)
			continue
		}
		newJob(payload.RoomId)
	}
}

func newJob(roomId int64) {
	worker := newWorker(roomId, closeChan)

	workersLock.Lock()
	defer workersLock.Unlock()
	workers[roomId] = worker
}

func generateReport() *model.ReportPayload {
	wil, speed := getWorkerInfoListAndSpeed()
	basicInfo := &model.BasicWorkerInfo{
		WorkerId: workerId,
		Capacity: conf.Conf.SpiderWorker.Capacity,
		Working:  len(workers),
		CpuCount: util.CpuCount(),
		MemUsage: util.MemUsage(),
		Speed:    speed,
	}
	return &model.ReportPayload{
		BasicWorkerInfo: basicInfo,
		Workers:         wil,
	}
}

func getWorkerInfoListAndSpeed() ([]*model.WorkerInfo, float64) {
	speed := 0.0
	list := make([]*model.WorkerInfo, 0, len(workers))
	workersLock.RLock()
	defer workersLock.RUnlock()
	for _, worker := range workers {
		list = append(list, worker.GetWorkerInfo())
		speed += worker.speeder.GetSpeed()
	}
	return list, speed
}
