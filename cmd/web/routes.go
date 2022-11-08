package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.home)
	router.HandlerFunc(http.MethodGet, "/game/view/:id", app.gameView)
	router.HandlerFunc(http.MethodGet, "/game/create", app.gameCreate)
	router.HandlerFunc(http.MethodPost, "/game/create", app.gameCreatePost)
	router.HandlerFunc(http.MethodGet, "/game/catalogView", app.catalogView)
	router.HandlerFunc(http.MethodGet, "/user/view/:id", app.userView)
	router.HandlerFunc(http.MethodGet, "/user/create", app.userCreate)
	router.HandlerFunc(http.MethodPost, "/user/create", app.userCreatePost)

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
