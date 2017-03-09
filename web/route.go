package web

import (
	"encoding/json"
	"log"
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

	e.GET("/working-room/count", func(c echo.Context) error {
		count, err := service.CountWorkingRoom()
		if err != nil {
			log.Println("service.CountWorkingRoom():", err)
			c.JSON(http.StatusOK, model.NewApiResponse(
				model.CodeInternalServerError, model.MessageInternalServerError))
			return nil
		}

		ret := model.NewApiOKResponse()
		ret.Data = &model.ApiWorkingRoomCount{
			Count: count,
		}

		c.JSON(http.StatusOK, ret)
		return nil
	})

	e.GET("/working-room/list", func(c echo.Context) error {
		var input model.ReqPager
		if err := c.Bind(&input); err != nil {
			return err
		}
		if input.Offset < 0 || input.Limit < 0 {
			c.JSON(http.StatusOK, model.NewApiResponse(
				model.CodeOutOfRange, model.MessageOutOfRange))
			return nil
		}

		count, err := service.CountWorkingRoom()
		if err != nil {
			log.Println("service.CountWorkingRoom():", err)
			c.JSON(http.StatusOK, model.NewApiResponse(
				model.CodeInternalServerError, model.MessageInternalServerError))
			return nil
		}

		workingRidList, err := service.ListWorkingRoom()
		if err != nil {
			log.Println("service.ListWorkingRoom():", err)
			c.JSON(http.StatusOK, model.NewApiResponse(
				model.CodeInternalServerError, model.MessageInternalServerError))
			return nil
		}
		workingRoomList, err := service.FindRoomByRidListPage(workingRidList, input.Offset, input.Limit)
		if err != nil {
			log.Println("service.FindRoomByRidList():", err)
			c.JSON(http.StatusOK, model.NewApiResponse(
				model.CodeInternalServerError, model.MessageInternalServerError))
			return nil
		}

		ret := model.NewApiOKResponse()
		ret.Data = &model.ApiWorkingRoomList{
			Total:           count,
			WorkingRoomList: workingRoomList,
		}

		c.JSON(http.StatusOK, ret)
		return nil
	})

}
