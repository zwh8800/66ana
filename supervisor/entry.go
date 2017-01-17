package supervisor

import (
	"log"

	"github.com/zwh8800/66ana/service"
	"github.com/zwh8800/66ana/util"
)

func Run() {
	started := false
	service.SubscribeReport(func(payload *service.ReportPayload, err error) {
		log.Println("Report:", util.JsonStringify(payload, false), "err:", err)
		if !started {
			if err != nil {
				return
			}

			service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
				RoomId: 67373,
			})
			service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
				RoomId: 613093,
			})
			service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
				RoomId: 606171,
			})
			service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
				RoomId: 10,
			})
			service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
				RoomId: 138286,
			})
			service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
				RoomId: 96291,
			})
			service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
				RoomId: 58428,
			})
			service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
				RoomId: 522423,
			})
			service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
				RoomId: 229346,
			})
			service.PublishStartSpider(payload.WorkerId, &service.StartSpiderPayload{
				RoomId: 573449,
			})

			started = true
		}
	})
}
