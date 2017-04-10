package main

import (
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	zmq "github.com/pebbe/zmq4"
	"github.com/zwh8800/66ana/conf"
	"github.com/zwh8800/66ana/jobworker"
	"github.com/zwh8800/66ana/spiderworker"
	"github.com/zwh8800/66ana/supervisor"
	"github.com/zwh8800/66ana/web"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	rand.Seed(time.Now().UnixNano())

	printVersions()

	if conf.Conf.Supervisor.IsSupervisor {
		go supervisor.Run()
	}

	if conf.Conf.SpiderWorker.IsSpiderWorker {
		go spiderworker.Run()
	}

	if conf.Conf.Web.IsWeb {
		go web.Run()
	}

	if conf.Conf.JobWorker.IsJobWorker {
		go jobworker.Run()
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch

	stack := make([]byte, 4*1024*1024)
	runtime.Stack(stack, true)
	log.Println(string(stack))
}

func printVersions() {
	major, minor, patch := zmq.Version()
	log.Printf("zeromq version: %d.%d.%d", major, minor, patch)
}
