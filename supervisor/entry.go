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

			started = true
		}
	})
}
