package main

import (
	"IOiyn.kz/internal/models"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	games, err := app.games.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Games = games
	app.render(w, http.StatusOK, "home.tmpl", data)
}
func (app *application) gameView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	game, err := app.games.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Game = game
	app.render(w, http.StatusOK, "gameView.tmpl", data)
}
func (app *application) gameCreatePost(w http.ResponseWriter, r *http.Request) {
	user, err := app.users.GetById(1)
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := app.games.Insert(user.ID, "dota", "kind of shit", 0, time.Date(2012, time.July, 9, 0, 0, 0, 0, time.UTC))
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/game/view/%d", id), http.StatusSeeOther)
}

func (app *application) userCreatePost(w http.ResponseWriter, r *http.Request) {
	name := "Zaur"
	nickname := "Lagmazavr"
	balance := 100
	email := "zaur@gmail.com"
	password := "password"
	id, err := app.users.Insert(name, nickname, balance, email, password)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/user/view/%d", id), http.StatusSeeOther)
}

func (app *application) gameCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "gameCreate.tmpl", data)
}

func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display the form for creating a new snippet..."))
}

func (app *application) userView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	user, err := app.users.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	data := app.newTemplateData(r)
	data.User = user
	app.render(w, http.StatusOK, "userView.tmpl", data)
}

func (app *application) catalogView(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Display a specific catalog of games ")
}
