package main

import (
	"os"
	"os/signal"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zwh8800/66ana/conf"
	"github.com/zwh8800/66ana/spiderworker"
	"github.com/zwh8800/66ana/supervisor"
)

func main() {
	if conf.Conf.Supervisor.IsSupervisor {
		supervisor.Run()
	}

	if conf.Conf.SpiderWorker.IsSpiderWorker {
		spiderworker.Run()
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, os.Kill)
	<-ch
}
