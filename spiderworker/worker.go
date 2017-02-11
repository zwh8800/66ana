package spiderworker

import (
	"log"
	"time"

	"github.com/zwh8800/66ana/model"
	"github.com/zwh8800/66ana/service"
	"github.com/zwh8800/66ana/spider"
	"github.com/zwh8800/66ana/util"
	"golang.org/x/net/proxy"
)

const minOnlineCount = 100

var proxyPool = NewProxyPoll(proxyList)

type worker struct {
	roomId    int64
	closeChan chan int64
	closed    chan bool
	spider    *spider.Spider

	speeder *util.Speedometer
}

func newWorker(roomId int64, closeChan chan int64) *worker {
	w := &worker{
		roomId:    roomId,
		closeChan: closeChan,
		closed:    make(chan bool),

		speeder: util.NewSpeedometer(),
	}
	go func() {
		var err error

		var dialer proxy.Dialer = nil
		proxyAddr := proxyPool.GetRandomProxy()
		if proxyAddr != "" {
			dialer, err = proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
			if err != nil {
				log.Println("proxy.SOCKS5:", err)
				proxyPool.DeleteProxy(proxyAddr)
				log.Println("proxyPool length:", proxyPool.Length())

				closeChan <- roomId
				return
			}
		}

		w.spider, err = spider.NewSpider(roomId, dialer)
		if err != nil {
			log.Println("spider.NewSpider:", err)
			if proxyAddr != "" {
				proxyPool.DeleteProxy(proxyAddr)

				log.Println("proxyPool length:", proxyPool.Length())
			}

			closeChan <- roomId
			return
		}

		go w.run()
	}()

	return w
}

func (w *worker) run() {
	w.pullRoomInfo()

	ticker := time.Tick(5 * time.Second)
	round := 0
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
				if round%2 == 0 {
					w.pullRoomInfo()
				} else {
					w.checkSpiderStatus()
				}
				round++
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
	needIncreaseCounter := true
	switch message["type"] {
	case "chatmsg":
		_, err := service.InsertDyDanmu(message)
		if err != nil {
			log.Println("service.InsertDyDanmu:", err)
		}

	case "dgb":
		_, err := service.InsertDyGiftHistory(message)
		if err != nil {
			log.Println("service.InsertDyGift:", err)
		}
	default:
		needIncreaseCounter = false
	}
	if needIncreaseCounter {
		w.speeder.Add()
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
	if room.OnlineCount < minOnlineCount {
		w.close()
	}
}

func (w *worker) close() {
	w.spider.Close()
	close(w.closed)

	w.closeChan <- w.roomId
}

func (w *worker) GetWorkerInfo() *model.WorkerInfo {
	return &model.WorkerInfo{
		RoomId: w.roomId,
		Speed:  w.speeder.GetSpeed(),
	}
}
