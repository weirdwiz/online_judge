package main

import (
	"fmt"
	"log"
	"net/http"
)

func InitializeHttpServer() {
	r := NewRouter()
	log.Panic(http.ListenAndServe(":8080", r))
}

func InitializeDatabase() {
	DB = &DBClient{}
	DB.Initialize("bolt.db")
}

func main() {
	fmt.Println("Hello, World!")
	InitializeDatabase()
	InitializeHttpServer()
}
