package main

import (
	"log"
	"net/http"

	"github.com/weirdwiz/online_judge/authentication/internal/app/dbclient"
	"github.com/weirdwiz/online_judge/authentication/internal/app/service"
)

func InitializeHttpServer() {
	r := service.NewRouter()
	log.Panic(http.ListenAndServe(":8080", r))
}

func InitializeDatabase() {
	service.DBClient = &dbclient.DBClient{}
	service.DBClient.Initialize("bolt.db")
}

func main() {
	InitializeDatabase()
	InitializeHttpServer()
}
