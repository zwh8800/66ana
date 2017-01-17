package service

import (
	"encoding/json"
	"errors"

	"github.com/zwh8800/66ana/constants"
)

type ReportPayload struct {
	WorkerId string
	Capacity int
	Working  int

	RoomIdList []int64

	CpuCount int
	MemUsage int64
}

type reportHandler func(*ReportPayload, error)

func SubscribeReport(h reportHandler) error {
	return subscribe(constants.EventReport, func() interface{} { return &ReportPayload{} },
		func(payload interface{}, err error) {
			report, _ := payload.(*ReportPayload)
			h(report, err)
		})
}

func PublishReport(payload *ReportPayload) error {
	return publish(constants.EventReport, payload)
}

type StartSpiderPayload struct {
	RoomId int64
}

type startSpiderHandler func(*StartSpiderPayload, error)

func SubscribeStartSpider(workerId string, h startSpiderHandler) error {
	return subscribe(constants.EventStartSpider+":"+workerId, func() interface{} { return &StartSpiderPayload{} },
		func(payload interface{}, err error) {
			report, _ := payload.(*StartSpiderPayload)
			h(report, err)
		})
}

func PublishStartSpider(id string, payload *StartSpiderPayload) error {
	return publish(constants.EventStartSpider+":"+id, payload)
}

type SpiderClosedPayload struct {
	WorkerId string
	RoomId   int64

	ReportPayload
}

type spiderClosedHandler func(*SpiderClosedPayload, error)

func SubscribeSpiderClosed(h spiderClosedHandler) error {
	return subscribe(constants.EventSpiderClosed, func() interface{} { return &SpiderClosedPayload{} },
		func(payload interface{}, err error) {
			report, _ := payload.(*SpiderClosedPayload)
			h(report, err)
		})
}

func PublishSpiderClosed(payload *SpiderClosedPayload) error {
	return publish(constants.EventSpiderClosed, payload)
}

func subscribe(channel string, genPayload func() interface{}, h func(interface{}, error)) error {
	if h == nil {
		return errors.New("h should not be nil")
	}
	sub, err := redisClient.Subscribe(channel)
	if err != nil {
		return err
	}
	go func() {
		for {
			message, err := sub.ReceiveMessage()
			if err != nil {
				h(nil, err)
			}
			payload := genPayload()
			if err := json.Unmarshal([]byte(message.Payload), payload); err != nil {
				h(nil, err)
			}
			h(payload, nil)
		}
	}()

	return nil
}

func publish(channel string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	redisClient.Publish(channel, string(data))
	return nil
}
