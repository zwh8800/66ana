package jobworker

import (
	"log"

	"github.com/zwh8800/66ana/jobworker/createfurthertable"
	"github.com/zwh8800/66ana/model"
	"github.com/zwh8800/66ana/service"
)

func Run() {
	go func() {
		for {
			payload, err := service.PullJob()
			if err != nil {
				log.Println("service.PullJob():", err)
			}
			go dispatchJob(payload)
		}
	}()
}

func dispatchJob(payload *model.JobPayload) {
	switch payload {
	case "createFurtherTable":
		createfurthertable.CreateFurtherTable(payload.Data)
	}
}
