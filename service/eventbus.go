package service

import (
	"encoding/json"
	"errors"

	"github.com/zwh8800/66ana/constants"
	"github.com/zwh8800/66ana/model"
)

func SubscribeReport(h func(*model.ReportPayload, error)) error {
	return subscribe(constants.EventReport, func() interface{} { return &model.ReportPayload{} },
		func(payload interface{}, err error) {
			report, _ := payload.(*model.ReportPayload)
			h(report, err)
		})
}

func PublishReport(payload *model.ReportPayload) error {
	return publish(constants.EventReport, payload)
}

func SubscribeStartSpider(workerId string, h func(*model.StartSpiderPayload, error)) error {
	return subscribe(constants.EventStartSpider+":"+workerId, func() interface{} { return &model.StartSpiderPayload{} },
		func(payload interface{}, err error) {
			report, _ := payload.(*model.StartSpiderPayload)
			h(report, err)
		})
}

func PublishStartSpider(id string, payload *model.StartSpiderPayload) error {
	return publish(constants.EventStartSpider+":"+id, payload)
}

func SubscribeSpiderClosed(h func(*model.SpiderClosedPayload, error)) error {
	return subscribe(constants.EventSpiderClosed, func() interface{} { return &model.SpiderClosedPayload{} },
		func(payload interface{}, err error) {
			report, _ := payload.(*model.SpiderClosedPayload)
			h(report, err)
		})
}

func PublishSpiderClosed(payload *model.SpiderClosedPayload) error {
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
				return
			}
			payload := genPayload()
			if err := json.Unmarshal([]byte(message.Payload), payload); err != nil {
				h(nil, err)
				return
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
