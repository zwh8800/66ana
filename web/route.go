package web

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/zwh8800/66ana/model"
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

	e.GET("/working-room-queue", func(c echo.Context) error {
		workers, err := service.ListWorkers()
		if err != nil {
			return err
		}
		workerDetail := make(map[string][]*model.DyRoom)

		for _, workerId := range workers {
			workingRidList, err := service.ListWorkingRoomQueue(workerId)
			if err != nil {
				return err
			}
			workingRoomList, err := service.FindRoomByRidList(workingRidList)
			if err != nil {
				return err
			}
			workerDetail[workerId] = workingRoomList
		}

		c.JSON(http.StatusOK, workerDetail)
		return nil
	})
}
