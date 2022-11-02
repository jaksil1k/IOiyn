package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/IOiyn/view", app.gameView)
	mux.HandleFunc("/IOiyn/create", app.gameCreate)
	mux.HandleFunc("/OIiyn/catalogView", app.catalogView)

	return mux
}
