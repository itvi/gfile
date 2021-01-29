package main

import (
	"flag"
	"gfile/internal/handler"
	"log"
	"net/http"
)

func main() {
	//port := flag.String("p", ":8000", "TCP address port")
	//dir := flag.String("d", ".", "User's directory")
	flag.Parse()

	config, db := handler.Config()
	defer db.Close()

	server := &http.Server{
		Addr:    "localhost:9000",
		Handler: config.Route(),
	}

	log.Println("Starting...", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
