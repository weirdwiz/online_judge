package service

import (
	"log"
	"net/http"
)

func StarWebServer(port string) {
	r := NewRouter()
	http.Handle("/", r)
	log.Println("Starting on Port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println(err.Error())
	}
}
