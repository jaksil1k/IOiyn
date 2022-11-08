package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/game/view", app.gameView)
	mux.HandleFunc("/game/create", app.gameCreate)
	mux.HandleFunc("/game/catalogView", app.catalogView)
	mux.HandleFunc("/user/view", app.userView)
	mux.HandleFunc("/user/create", app.userCreate)

	standart := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standart.Then(mux)
}
