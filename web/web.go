package web

import (
	"log"

	"github.com/labstack/echo"
	"github.com/zwh8800/66ana/conf"
)

func Run() {
	e := echo.New()
	route(e)
	if err := e.Start(conf.Conf.Web.Address); err != nil {
		log.Println("web error:", err)
	}
}
