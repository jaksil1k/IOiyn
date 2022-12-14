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

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/game/view/:id", dynamic.ThenFunc(app.gameView))
	router.Handler(http.MethodGet, "/game/catalogView", dynamic.ThenFunc(app.catalogView))

	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))

	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/user/view/:id", protected.ThenFunc(app.userView))
	router.Handler(http.MethodGet, "/user/another_view/:id", protected.ThenFunc(app.userView))
	router.Handler(http.MethodGet, "/game/create", protected.ThenFunc(app.gameCreate))
	router.Handler(http.MethodPost, "/game/create", protected.ThenFunc(app.gameCreatePost))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPost))
	router.Handler(http.MethodGet, "/user/change_info", protected.ThenFunc(app.changeInfo))
	router.Handler(http.MethodPost, "/user/change_info", protected.ThenFunc(app.changeInfoPut))
	router.Handler(http.MethodGet, "/user/update_balance", protected.ThenFunc(app.updateBalance))
	router.Handler(http.MethodPost, "/user/update_balance", protected.ThenFunc(app.updateBalancePut))
	router.Handler(http.MethodGet, "/user/change_password", protected.ThenFunc(app.updatePassword))
	router.Handler(http.MethodPost, "/user/change_password", protected.ThenFunc(app.updatePasswordPut))
	router.Handler(http.MethodPost, "/game/purchase/:id", protected.ThenFunc(app.purchaseGame))

	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
