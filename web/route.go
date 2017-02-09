package web

import (
	"encoding/json"
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
		ret := make(map[string]interface{})
		ret["workingRoomIdList"] = workingRidList
		ret["workingRoomList"] = workingRoomList

		c.JSON(http.StatusOK, ret)
		return nil
	})

	e.GET("/workers", func(c echo.Context) error {
		workers, err := service.ListWorkerIdList()
		if err != nil {
			return err
		}
		workerDetail := make(map[string]map[string]interface{})

		for _, workerId := range workers {
			workerInfoStr, err := service.GetWorkerInfo(workerId)
			if err != nil {
				return err
			}
			var workerInfo map[string]interface{}
			if err := json.Unmarshal([]byte(workerInfoStr), &workerInfo); err != nil {
				return err
			}

			queueRidList, err := service.ListWorkingRoomQueue(workerId)
			if err != nil {
				return err
			}
			queueRoomList, err := service.FindRoomByRidList(queueRidList)
			if err != nil {
				return err
			}
			workerInfo["queueRoomList"] = queueRoomList
			workerDetail[workerId] = workerInfo
		}

		c.JSON(http.StatusOK, workerDetail)
		return nil
	})
}
