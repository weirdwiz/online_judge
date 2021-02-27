package main

import (
	"fmt"

	"github.com/weirdwiz/online_judge/compile_microservice/internal/app/service"
)

func main() {
	fmt.Printf("Starting\n")
	service.StartWebServer("5050")
	fmt.Println("Started at 5050")
}
