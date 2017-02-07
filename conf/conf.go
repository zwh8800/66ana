package conf

import (
	"flag"

	"github.com/BurntSushi/toml"
)

type config struct {
	DB struct {
		Driver        string
		Dsn           string
		MaxConnection int
	}

	Redis struct {
		Addr     string
		Password string
		DB       int
	}

	Supervisor struct {
		IsSupervisor bool
	}

	SpiderWorker struct {
		IsSpiderWorker bool
		Capacity       int
	}

	Web struct {
		IsWeb   bool
		Address string
	}
}

var Conf config

func init() {
	configPath := flag.String("config", "66ana.toml", "specify a config file")
	flag.Parse()
	if _, err := toml.DecodeFile(*configPath, &Conf); err != nil {
		panic(err)
	}
}
