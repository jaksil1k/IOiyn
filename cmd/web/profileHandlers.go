package main

import (
	"IOiyn.kz/internal/models"
	"IOiyn.kz/internal/validator"
	"errors"
	"fmt"
	"net/http"
)

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

type userChangePasswordForm struct {
	Password            string `form:"password"`
	RePassword          string `form:"rePassword"`
	validator.Validator `form:"-"`
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
	data.Form = userLoginForm{}
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

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

	form.CheckField(validator.NotBlank(form.Name), "name", "Name cannot be blank")
	form.CheckField(validator.MaxChars(form.Name, 255), "name", "Name field cannot be more than 255 character long")
	form.CheckField(validator.NotBlank(form.Nickname), "nickname", "Nickname cannot be blank")
	form.CheckField(validator.MaxChars(form.Nickname, 255), "nickname", "Nickname field cannot be more than 255 character long")

	if !form.Valid() {
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
		data.Form = form
		data.User = user
		app.render(w, http.StatusUnprocessableEntity, "changeInfo.tmpl", data)
		return
	}

	err = app.users.UpdateUserInfo(id, form.Name, form.Nickname)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your update was successful.")

	//url := "/user/view/" + string(rune(id))

	http.Redirect(w, r, fmt.Sprintf("/user/view/%d", id), http.StatusSeeOther)

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
	data.Form = userUpdateBalanceForm{}
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

	form.CheckField(validator.MaxInt(form.Balance, 10000), "balance", "you cannot take more than 100$ freely")

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}

	id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")

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

	if !form.Valid() {

		data := app.newTemplateData(r)
		data.User = user
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "changeBalance.tmpl", data)
		return
	}

	fmt.Println(form.Balance, " ", user.Balance)
	err = app.users.UpdateBalance(id, form.Balance, user.Balance)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "Your update was successful.")

	http.Redirect(w, r, fmt.Sprintf("/user/view/%d", id), http.StatusSeeOther)

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

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "changeBalance.tmpl", data)
		return
	}

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

	http.Redirect(w, r, fmt.Sprintf("/user/view/%d", id), http.StatusSeeOther)

}

func (app *application) updatePassword(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, err)
		return
	}
	data := app.newTemplateData(r)
	data.Form = userChangePasswordForm{}
	app.render(w, http.StatusOK, "changePassword.tmpl", data)
}
