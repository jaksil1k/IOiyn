package main

import (
	"IOiyn.kz/internal/models"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	games, err := app.games.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, http.StatusOK, "home.tmpl", &templateData{
		Games: games,
	})
}
func (app *application) gameView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

	app.render(w, http.StatusOK, "gameView.tmpl", &templateData{
		Game: game,
	})
}
func (app *application) gameCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

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
	http.Redirect(w, r, fmt.Sprintf("/game/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
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

	http.Redirect(w, r, fmt.Sprintf("/user/view?id=%d", id), http.StatusSeeOther)
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
	app.render(w, http.StatusOK, "userView.tmpl", &templateData{
		User: user,
	})
}

func (app *application) catalogView(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Display a specific catalog of games ")
}
