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

	// configuration
	config, db := handler.Config(*DIR)
	defer db.Close()

	go config.File.Watchdog(*DIR)

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
