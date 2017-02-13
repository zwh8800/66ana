package main

import (
	"log"
	"time"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	reportRequest, err := zmq.NewSocket(zmq.REQ)
	if err != nil {
		panic(err)
	}
	if err := reportRequest.Bind("tcp://*:13378"); err != nil {
		panic(err)
	}

	clientMap := make(map[string]bool)

	for {
		for i := 0; len(clientMap) == 0 || i < len(clientMap); i++ {
			if _, err := reportRequest.Send("need report", 0); err != nil {
				panic(err)
			}

			log.Println("send ok")

			data, err := reportRequest.Recv(0)
			if err != nil {
				panic(err)
			}
			log.Println(data)
			clientMap[data] = true
		}

		time.Sleep(1 * time.Second)
	}
}
