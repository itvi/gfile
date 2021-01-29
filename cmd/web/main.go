package main

import (
	"gfile/internal/handler"
	"log"
	"net/http"
)

func main() {

	config, db := handler.Config()
	defer db.Close()

	server := &http.Server{
		Addr:    ":9000",
		Handler: config.Route(),
	}

	log.Println("Starting...", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
