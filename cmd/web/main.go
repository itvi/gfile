package main

import (
	"gfile/internal/handler"
)

// MOVE main.go TO ROOT PATH !

func main() {
	// var PORT = flag.String("p", ":8000", "TCP address port")
	// var DIR = flag.String("d", ".", "Monitor directory")
	// var SERVICE = flag.String("s", "", "Control the system service.")
	// flag.Parse()

	// // configuration
	// config, db := handler.Config(*DIR)
	// defer db.Close()

	// go config.File.Watchdog(*DIR)

	// server := &http.Server{
	// 	Addr:    *PORT,
	// 	Handler: config.Route(),
	// }

	// log.Println("Starting...", server.Addr)
	// err := server.ListenAndServe()
	// if err != nil {
	// 	log.Println(err)
	// }

	// handler.Service(*PORT, *DIR, *SERVICE)
	handler.Service()
}
