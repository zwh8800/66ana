package createfurthertable

import (
	"log"

	"github.com/zwh8800/66ana/service"
)

func CreateFurtherTable(interface{}) {
	if err := service.CreateFurtherTable(5); err != nil {
		log.Println("service.CreateFurtherTable(5):", err)
	}
}
