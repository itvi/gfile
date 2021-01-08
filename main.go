package main

import (
	"flag"
	"fmt"
	"gfile/handler"
	"log"
	"net/http"
)

// var root = "C:/test"

func main() {
	// port,directory
	port := flag.String("p", "localhost:8000", "TCP address port")
	dir := flag.String("d", ".", "User's directory")
	flag.Parse()

	// static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler.Index(*dir))
	http.HandleFunc("/dl", handler.Download(*dir))
	http.HandleFunc("/zip", handler.Zip(*dir))

	fmt.Println("start...", *port)

	server := &http.Server{
		Addr: *port,
	}

	log.Fatal(server.ListenAndServe())
}
