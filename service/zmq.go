package service

import (
	"encoding/json"
	"fmt"

	zmq "github.com/pebbe/zmq4"
	"github.com/zwh8800/66ana/conf"
	"github.com/zwh8800/66ana/model"
)

const zmqAddressFormat = "tcp://%s:%d"

const (
	workerReportPort = 23330 + iota
	pushPullJobPort
)

var (
	workerReportReqSocket *zmq.Socket = nil
	workerReportRepSocket *zmq.Socket = nil
)

func GetWorkerReport() (*model.ReportPayload, error) {
	if workerReportReqSocket == nil {
		var err error
		addr := fmt.Sprintf(zmqAddressFormat, conf.Conf.Zeromq.Addr, workerReportPort)
		workerReportReqSocket, err = zmq.NewSocket(zmq.REQ)
		if err != nil {
			return nil, err
		}
		if err := workerReportReqSocket.Bind(addr); err != nil {
			return nil, err
		}
	}
	if _, err := workerReportReqSocket.Send("", 0); err != nil {
		return nil, err
	}

	data, err := workerReportReqSocket.Recv(0)
	if err != nil {
		return nil, err
	}
	report := &model.ReportPayload{}
	if err := json.Unmarshal([]byte(data), report); err != nil {
		return nil, err
	}
	return report, nil
}

func ServeWorkerReport(reportCallback func(err error) *model.ReportPayload) error {
	if workerReportRepSocket == nil {
		var err error
		addr := fmt.Sprintf(zmqAddressFormat, conf.Conf.Zeromq.Addr, workerReportPort)
		workerReportRepSocket, err = zmq.NewSocket(zmq.REP)
		if err != nil {
			return err
		}
		if err := workerReportRepSocket.Connect(addr); err != nil {
			return err
		}
	}

	go func() {
		for {
			if _, err := workerReportRepSocket.Recv(0); err != nil {
				reportCallback(err)
				continue
			}

			report := reportCallback(nil)
			data, err := json.Marshal(report)
			if err != nil {
				reportCallback(err)
				continue
			}

			if _, err := workerReportRepSocket.Send(string(data), 0); err != nil {
				reportCallback(err)
				continue
			}
		}
	}()
	return nil
}
