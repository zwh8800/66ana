package spiderworker

import (
	"log"
	"time"

	"github.com/zwh8800/66ana/service"
	"github.com/zwh8800/66ana/spider"
)

type worker struct {
	roomId    int64
	closeChan chan int64
	spider    *spider.Spider
}

func newWorker(roomId int64, closeChan chan int64) *worker {
	w := &worker{
		roomId:    roomId,
		closeChan: closeChan,
	}
	var err error
	w.spider, err = spider.NewSpider(roomId, nil)
	if err != nil {
		log.Println("spider.NewSpider:", err)

		closeChan <- roomId
		return nil
	}

	go w.run()

	return w
}

func (w *worker) run() {
	w.pullRoomInfo()

	ticker := time.Tick(10 * time.Second)
	for {
		select {
		case message := <-w.spider.GetMessageChan():
			w.handleMessage(message)

		case <-ticker:
			w.pullRoomInfo()
		}
	}
}

func (w *worker) handleMessage(message map[string]string) {
	switch message["type"] {
	case "chatmsg":
		_, err := service.InsertDyDanmu(message)
		if err != nil {
			log.Println("service.InsertDyDanmu:", err)
		}

	case "dgb":
	default:
	}
}

func (w *worker) pullRoomInfo() {
	// TODO: 拉取房间状态
	roomInfo, err := w.spider.GetRoomInfo()
	if err != nil {
		log.Println("spider.GetRoomInfo:", err)
		return
	}
	if _, err := service.InsertDyRoom(roomInfo); err != nil {
		log.Println("service.InsertDyRoom:", err)
		return
	}
}

func (w *worker) GetRoomId() int64 {
	return w.roomId
}
