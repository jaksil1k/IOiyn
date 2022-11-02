package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/IOiyn/view", gameView)
	mux.HandleFunc("/IOiyn/create", gameCreate)
	mux.HandleFunc("/OIiyn/catalogView", catalogView)
	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
