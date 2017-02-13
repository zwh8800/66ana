package main

import (
	"fmt"
	"log"

	zmq "github.com/pebbe/zmq4"
	uuid "github.com/satori/go.uuid"
)

func main() {
	reportRequest, err := zmq.NewSocket(zmq.REP)
	if err != nil {
		panic(err)
	}
	if err := reportRequest.Connect("tcp://localhost:13378"); err != nil {
		panic(err)
	}

	workerId := uuid.NewV4().String()

	for {
		data, err := reportRequest.Recv(0)
		if err != nil {
			panic(err)
		}
		log.Println(data)

		if _, err := reportRequest.Send(fmt.Sprintln("workerId:", workerId), 0); err != nil {
			panic(err)
		}
	}
}
