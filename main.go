package main

import (
	"database/sql"
	"flag"
	"fmt"
	"gfile/handler"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// var root = "C:/test"

func main() {
	// port,directory
	port := flag.String("p", "localhost:8000", "TCP address port")
	dir := flag.String("d", ".", "User's directory")
	flag.Parse()

	// database
	db, err := openDB("./file.db")
	if err != nil {
		err = fmt.Errorf("Open db error: %w", err)
		log.Panic(err)
	}
	defer db.Close()

	// static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", handler.Index(*dir))
	http.HandleFunc("/dl", handler.Download(*dir))
	http.HandleFunc("/zip", handler.Zip(*dir))
	http.HandleFunc("/search", handler.Search(db))
	http.HandleFunc("/rebuild", handler.Rebuild(*dir, db))

	fmt.Println("start...", *port)

	server := &http.Server{
		Addr: *port,
	}

	log.Fatal(server.ListenAndServe())
}

func openDB(cnn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cnn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
