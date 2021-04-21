package service

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
)

func StarWebServer(port string) {
	r := NewRouter()
	http.Handle("/", r)

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})

	log.Println("Starting on Port " + port)
	err := http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(r))
	if err != nil {
		log.Println(err.Error())
	}
}
