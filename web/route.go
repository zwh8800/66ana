package web

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/zwh8800/66ana/service"
)

func route(e *echo.Echo) {
	e.File("/", "static/index.html")
	e.Static("/static", "static")

	e.GET("/working-room", func(c echo.Context) error {
		workingRidList, err := service.ListWorkingRoom()
		if err != nil {
			return err
		}
		workingRoomList, err := service.FindRoomByRidList(workingRidList)
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, workingRoomList)
		return nil
	})
}
