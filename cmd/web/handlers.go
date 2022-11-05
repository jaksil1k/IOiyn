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
	w.Write([]byte("Hello from IOiyn"))
}
func (app *application) gameView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	fmt.Println(id)
	game, err := app.games.GetByGameId(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
	}

	fmt.Fprintf(w, "%+v", game)
}
func (app *application) gameCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	userName := "Zaur"
	userNickname := "Lagmazavr"
	userBalance := 100
	userEmail := "zaur@gmail.com"
	userPassword := "password"

	userId, err := app.games.InsertIntoUsers(userName, userNickname, userBalance, userEmail, userPassword)
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := app.games.InsertIntoGames(userId, "dota", "kind of shit", 0, time.Date(2012, time.July, 9, 0, 0, 0, 0, time.UTC))
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
	id, err := app.games.InsertIntoUsers(name, nickname, balance, email, password)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/user/view?id=%d", id), http.StatusSeeOther)
}

func (app *application) catalogView(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Display a specific catalog of games ")
}
