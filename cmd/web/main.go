package main

import (
	"flag"
	"gfile/internal/handler"
	"log"
	"net/http"
)

func main() {
	var PORT = flag.String("p", ":8000", "TCP address port")
	var DIR = flag.String("d", ".", "User's directory")
	flag.Parse()

	config, db := handler.Config(*DIR)
	defer db.Close()

	server := &http.Server{
		Addr:    *PORT,
		Handler: config.Route(),
	}

	log.Println("Starting...", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
