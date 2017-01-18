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
	w.logMessage(message)
	switch message["type"] {
	case "chatmsg":
		danmu, err := service.InsertDyDanmu(message)
		if err != nil {
			log.Println("service.InsertDyDanmu:", err)
		}

		log.Println(danmu, "inserted")
	case "dgb":
	default:
	}
}

func (w *worker) logMessage(message map[string]string) {
	switch message["type"] {
	case "chatmsg":
		colorCode := ""
		switch message["col"] {
		case "1": // 红
			colorCode = "\033[1;91m"
		case "2": // 蓝
			colorCode = "\033[1;94m"
		case "3": // 绿
			colorCode = "\033[1;92m"
		case "4": // 黄
			colorCode = "\033[1;93m"
		case "5": // 紫
			colorCode = "\033[1;38;5;129m"
		case "6": // 粉
			colorCode = "\033[1;38;5;213m"
		default:
			colorCode = "\033[1m"
		}
		log.Printf("%s(%s): %s%s\033[0m", message["nn"], message["uid"], colorCode, message["txt"])
	case "dgb":
		hits := message["hits"]
		if hits == "" {
			hits = "1"
		}
		log.Printf("%s(%s) \033[90m送出 %s (%s 连击)\033[0m", message["nn"], message["uid"], w.spider.GetGiftMap()[message["gfid"]], hits)
	default:
		// log.Printf("%#v", message)
	}
}

func (w *worker) pullRoomInfo() {
	// TODO: 拉取房间状态
}

func (w *worker) GetRoomId() int64 {
	return w.roomId
}
