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
	closed    chan bool
	spider    *spider.Spider
}

func newWorker(roomId int64, closeChan chan int64) *worker {
	w := &worker{
		roomId:    roomId,
		closeChan: closeChan,
		closed:    make(chan bool),
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
		case <-w.closed:
			return
		default:
			select {
			case message := <-w.spider.GetMessageChan():
				if message == nil {
					w.checkSpiderStatus()
				} else {
					w.handleMessage(message)
				}
			case <-ticker:
				w.pullRoomInfo()
				w.checkSpiderStatus()
			}
		}
	}
}

func (w *worker) checkSpiderStatus() {
	switch w.spider.GetStatus() {
	case spider.StatusError:
		log.Println("spider.StatusError:", w.spider.GetLastError())
		w.close()
	case spider.StatusClosed:
		w.close()
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
	roomInfo, err := w.spider.GetRoomInfo()
	if err != nil {
		log.Println("spider.GetRoomInfo:", err)
		return
	}
	room, err := service.InsertDyRoom(roomInfo)
	if err != nil {
		log.Println("service.InsertDyRoom:", err)
		return
	}
	if room.OnlineCount <= 0 {
		w.close()
	}
}

func (w *worker) close() {
	w.spider.Close()
	close(w.closed)

	w.closeChan <- w.roomId
}

func (w *worker) GetRoomId() int64 {
	return w.roomId
}
