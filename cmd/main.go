package main

import (
	"fmt"

	"github.com/weirdwiz/compile_microservice/cmd/service"
)

func main() {
	fmt.Printf("Starting\n")
	service.StartWebServer("5050")
	fmt.Println("Started at 5050")
}
