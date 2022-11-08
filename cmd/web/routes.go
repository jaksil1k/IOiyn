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

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/game/view/:id", dynamic.ThenFunc(app.gameView))
	router.Handler(http.MethodGet, "/game/create", dynamic.ThenFunc(app.gameCreate))
	router.Handler(http.MethodPost, "/game/create", dynamic.ThenFunc(app.gameCreatePost))
	router.Handler(http.MethodGet, "/game/catalogView", dynamic.ThenFunc(app.catalogView))
	router.Handler(http.MethodGet, "/user/view/:id", dynamic.ThenFunc(app.userView))
	router.Handler(http.MethodGet, "/user/create", dynamic.ThenFunc(app.userCreate))
	router.Handler(http.MethodPost, "/user/create", dynamic.ThenFunc(app.userCreatePost))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
