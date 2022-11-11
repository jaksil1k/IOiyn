package main

import (
	"IOiyn.kz/internal/models"
	"IOiyn.kz/internal/validator"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"time"
)

type userSignupForm struct {
	Name                string `form:"name"`
	Nickname            string `form:"nickname"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type gameCreateForm struct {
	CreatedBy           int    `form:"-"`
	Name                string `form:"name"`
	Cost                int    `form:"cost"`
	Description         string `form:"description"`
	validator.Validator `form:"-"`
}

type userChangeInfoForm struct {
	Name                string `form:"name"`
	Nickname            string `form:"nickname"`
	validator.Validator `form:"-"`
}

type userUpdateBalanceForm struct {
	Balance             int `form:"balance"`
	CurrentBalance      int `form:"currentBalance"`
	validator.Validator `form:"-"`
}

type gamePurchaseForm struct {
	Balance             int `form:"balance"`
	validator.Validator `form:"-"`
}

type userChangePasswordForm struct {
	Password            string `form:"password"`
	validator.Validator `form:"-"`
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
	user.Password = []byte("")
	game.AuthorName = user.Name
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

	var form gameCreateForm

	err = app.formDecoder.Decode(&form, r.PostForm)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "game cannot be blank")
	form.CheckField(validator.NotBlank(form.Description), "description", "description cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 100), "name", "name cannot be more than 100 characters long")
	form.CheckField(validator.MaxChars(form.Description, 10000), "description", "description cannot be more than 10000 characters long")
	form.CheckField(validator.MaxInt(form.Cost, 100000), "cost", "cost cannot be more than 100000")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusOK, "gameCreate.tmpl", data)
		return
	}

	id, err := app.games.Insert(app.sessionManager.GetInt(r.Context(), "authenticatedUserID"), form.Name, form.Description, form.Cost, time.Now())
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Game successfully created!")

	http.Redirect(w, r, fmt.Sprintf("/game/view/%d", id), http.StatusSeeOther)
}

func (app *application) gameCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = gameCreateForm{}
	app.render(w, http.StatusOK, "gameCreate.tmpl", data)
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, http.StatusOK, "userCreate.tmpl", data)
}

func (app *application) userView(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	currentId := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	if id != currentId {
		http.Redirect(w, r, fmt.Sprintf("/user/another_view/%d", id), http.StatusSeeOther)
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
	user.Password = []byte("")

	data := app.newTemplateData(r)
	data.PurchasedGames, err = app.users.UserPurchasedGamesView(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.CreatedGames, err = app.users.UserCreatedGamesView(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.User = user

	app.render(w, http.StatusOK, "userView.tmpl", data)
}

func (app *application) anotherUserView(w http.ResponseWriter, r *http.Request) {
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
	user.Password = []byte("")

	data := app.newTemplateData(r)
	data.CreatedGames, err = app.users.UserCreatedGamesView(id)
	if err != nil {
		app.serverError(w, err)
		return
	}
	data.User = user

	app.render(w, http.StatusOK, "anotherUserView.tmpl", data)
}

func (app *application) catalogView(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Display a specific catalog of games ")
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 255), "name", "Name field cannot be more than 255 character long")
	form.CheckField(validator.NotBlank(form.Nickname), "nickname", "Nickname cannot be blank")
	form.CheckField(validator.MaxChars(form.Nickname, 255), "nickname", "Nickname field cannot be more than 255 character long")
	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password cannot be less than 8")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "userCreate.tmpl", data)
		return
	}

	err = app.users.Insert(form.Name, form.Nickname, 1000, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "signup.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, http.StatusOK, "login.tmpl", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form.CheckField(validator.NotBlank(form.Email), "email", "Email field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "Email field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "Password field cannot be blank")
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")
			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, err)
		}
		return
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/game/create", http.StatusSeeOther)
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) purchaseGame(w http.ResponseWriter, r *http.Request) {
	var form gamePurchaseForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

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

	user, err := app.users.GetById(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	//form.CheckField(validator.MaxInt(game.Cost, game.))

	app.users.UpdateBalance(user.ID, -game.Cost, user.Balance)
}

func (app *application) changeInfo(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

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
	user.Password = []byte("")
	data.User = user
	app.render(w, http.StatusOK, "changeInfo.tmpl", data)
}

func (app *application) changeInfoPut(w http.ResponseWriter, r *http.Request) {
	var form userChangeInfoForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 255), "name", "Name field cannot be more than 255 character long")
	form.CheckField(validator.NotBlank(form.Nickname), "nickname", "Nickname cannot be blank")
	form.CheckField(validator.MaxChars(form.Nickname, 255), "nickname", "Nickname field cannot be more than 255 character long")

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = app.users.UpdateUserInfo(id, form.Name, form.Nickname)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your update was successful.")

	url := "/user/view/" + string(rune(id))

	http.Redirect(w, r, url, http.StatusSeeOther)

}

func (app *application) updateBalance(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

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
	user.Password = []byte("")
	data.User = user

	app.render(w, http.StatusOK, "changeBalance.tmpl", data)
}

func (app *application) updateBalancePut(w http.ResponseWriter, r *http.Request) {
	var form userUpdateBalanceForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(string(rune(form.Balance))), "balance", "Balance cannot be blank")
	form.CheckField(validator.MaxInt(form.Balance, 100), "balance", "you cannot take more than 100$ freely")

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = app.users.UpdateBalance(id, form.Balance, form.CurrentBalance)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your update was successful.")

	url := "/user/view/" + string(rune(id))

	http.Redirect(w, r, url, http.StatusSeeOther)

}

func (app *application) updatePasswordPut(w http.ResponseWriter, r *http.Request) {
	var form userChangePasswordForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Password), "password", "Password cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password cannot be less than 8")

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	err = app.users.UpdatePassword(id, form.Password)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your update was successful.")

	url := "/user/view/" + string(rune(id))

	http.Redirect(w, r, url, http.StatusSeeOther)

}

func (app *application) updatePassword(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	app.render(w, http.StatusOK, "changePassword.tmpl", data)
}

//func (app *application) userCreatePost(w http.ResponseWriter, r *http.Request) {
//	err := r.ParseForm()
//	if err != nil {
//		app.serverError(w, err)
//		return
//	}
//
//	var form userSignupForm
//
//	err = app.formDecoder.Decode(&form, r.PostForm)
//	if err != nil {
//		app.clientError(w, http.StatusBadRequest)
//	}
//
//	form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be blank")
//	form.CheckField(validator.MaxChars(form.Name, 255), "name", "Name field cannot be more than 50 character long")
//	form.CheckField(validator.NotBlank(form.Nickname), "nickname", "Nickname cannot be blank")
//	form.CheckField(validator.MaxChars(form.Nickname, 255), "nickname", "Nickname field cannot be more than 50 character long")
//	form.CheckField(validator.NotBlank(form.Email), "email", "Email cannot be blank")
//	form.CheckField(validator.NotBlank(form.Email), "password", "Password cannot be blank")
//	form.CheckField(validator.MinChars(form.Password, 8), "password", "Password cannot be less than 8")
//
//	if !form.Valid() {
//		data := app.newTemplateData(r)
//		data.Form = form
//		app.render(w, http.StatusUnprocessableEntity, "userSignup.tmpl", data)
//		return
//	}
//
//	err := app.users.Insert(form.Name, form.Nickname, 1000, form.Email, form.Password)
//	if err != nil {
//		app.serverError(w, err)
//		return
//	}
//
//	app.sessionManager.Put(r.Context(), "flash", "User successfully created!")
//
//	http.Redirect(w, r, fmt.Sprintf("/user/view/%d", id), http.StatusSeeOther)
//}
