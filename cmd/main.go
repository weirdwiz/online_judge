package main

import (
	"fmt"

	"github.com/weirdwiz/compile_microservice/cmd/service"
)

func main() {
	fmt.Printf("Starting\n")
	initializeBoltClient()
	service.StartWebServer("5050")
	fmt.Println("Started at 5050")
}

func initializeBoltClient() {
		service.DBClient = &dbclient.BoltClient{}
		service.DBClient.OpenBoltDb()
		service.DBClient.Seed()
}
