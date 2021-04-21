package service

import (
	"log"
	"net/http"
	"github.com/gorilla/handlers"
)

func StartWebServer(port string) {
	r := NewRouter()
	http.Handle("/", r)



	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"*"})




	log.Println("Starting on port " + port)
	err := http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(r))
	// err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Println("there's an error: " + err.Error())
	}
}
