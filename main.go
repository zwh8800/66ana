package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zwh8800/66ana/conf"
	"github.com/zwh8800/66ana/spiderworker"
	"github.com/zwh8800/66ana/supervisor"
	"github.com/zwh8800/66ana/web"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	rand.Seed(time.Now().UnixNano())

	if conf.Conf.Supervisor.IsSupervisor {
		go supervisor.Run()
	}

	if conf.Conf.SpiderWorker.IsSpiderWorker {
		go spiderworker.Run()
	}

	if conf.Conf.Web.IsWeb {
		go web.Run()
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch

	stack := make([]byte, 4*1024*1024)
	runtime.Stack(stack, true)
	log.Println(string(stack))
}
