package service

import (
	"log"
	"net/http"
)

func StartWebServer(port string) {
	r := NewRouter()
	http.Handle("/", r)
	log.Println("Starting on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("there's an error: " + err.Error())
	}
}
