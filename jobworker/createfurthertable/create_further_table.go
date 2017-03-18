package createfurthertable

import (
	"log"

	"github.com/zwh8800/66ana/service"
)

func CreateFurtherTable(interface{}) {
	log.Println("Job CreateFurtherTable start")
	if err := service.CreateFurtherTable(5); err != nil {
		log.Println("service.CreateFurtherTable(5):", err)
	}
	log.Println("Job CreateFurtherTable end")
}
