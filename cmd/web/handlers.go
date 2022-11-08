package main

import (
	"IOiyn.kz/internal/models"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type userCreateForm struct {
	Name        string
	Nickname    string
	Email       string
	Password    string
	FieldErrors map[string]string
}

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

	user, err := app.users.GetById(game.CreatedBy)
	user.Password = ""
	game.Author = user
	data := app.newTemplateData(r)
	data.Game = game

	app.render(w, http.StatusOK, "gameView.tmpl", data)
}
func (app *application) gameCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	id, err := app.games.Insert(1, "dota", "kind of shit", 0, time.Date(2012, time.July, 9, 0, 0, 0, 0, time.UTC))

	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/game/view/%d", id), http.StatusSeeOther)
}

func (app *application) userCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
		return
	}

	form := userCreateForm{
		Name:        r.PostForm.Get("name"),
		Nickname:    r.PostForm.Get("nickname"),
		Email:       r.PostForm.Get("email"),
		Password:    r.PostForm.Get("password"),
		FieldErrors: map[string]string{},
	}

	fieldErrors := make(map[string]string)

	if strings.TrimSpace(form.Name) == "" {
		fieldErrors["name"] = "this field can not be blank"
	} else if utf8.RuneCountInString(form.Name) > 50 {
		fieldErrors["name"] = "this field can not be more than 50 characters long"
	}

	if strings.TrimSpace(form.Nickname) == "" {
		fieldErrors["nickname"] = "this field can not be blank"
	} else if utf8.RuneCountInString(form.Nickname) > 30 {
		fieldErrors["nickname"] = "this field can not be more than 30 characters long"
	}

	if strings.TrimSpace(form.Email) == "" {
		fieldErrors["email"] = "This field can not be blank"
	}

	if strings.TrimSpace(form.Password) == "" {
		fieldErrors["email"] = "This field can not be blank"
	} else if utf8.RuneCountInString(form.Password) < 8 {
		fieldErrors["email"] = "This field can not be less than 8"
	}

	if len(fieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "userCreate.tmpl", data)
	}

	id, err := app.users.Insert(form.Name, form.Nickname, 1000, form.Email, form.Password)
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
	data := app.newTemplateData(r)
	data.Form = userCreateForm{}
	app.render(w, http.StatusOK, "userCreate.tmpl", data)
}

func (app *application) userView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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
